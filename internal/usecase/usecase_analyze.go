package usecase

import (
	"context"
	"sync"
	"time"

	"tmt/internal/config"
	"tmt/internal/entity"
	"tmt/internal/usecase/cache"
	"tmt/internal/usecase/modules/calendar"
	"tmt/internal/usecase/repo"
	"tmt/pkg/eventbus"
	"tmt/pkg/log"
)

// AnalyzeUseCase -.
type AnalyzeUseCase struct {
	repo HistoryRepo

	targetArr []*entity.StockTarget

	lastBelowMAStock map[string]*entity.StockHistoryAnalyze
	rebornMap        map[time.Time][]entity.Stock
	rebornLock       sync.Mutex

	tradeDay *calendar.Calendar

	logger *log.Log
	cc     *cache.Cache
	bus    *eventbus.Bus
}

func NewAnalyze() Analyze {
	uc := &AnalyzeUseCase{
		repo:             repo.NewHistory(config.Get().GetPostgresPool()),
		lastBelowMAStock: make(map[string]*entity.StockHistoryAnalyze),
		rebornMap:        make(map[time.Time][]entity.Stock),
		tradeDay:         calendar.Get(),
		logger:           log.Get(),
		cc:               cache.Get(),
		bus:              eventbus.New(),
	}

	uc.bus.SubscribeAsync(topicAnalyzeStockTargets, true, uc.findBelowQuaterMATargets)
	return uc
}

// GetRebornMap -.
func (uc *AnalyzeUseCase) GetRebornMap(ctx context.Context) map[time.Time][]entity.Stock {
	uc.rebornLock.Lock()
	if len(uc.lastBelowMAStock) != 0 {
		for _, s := range uc.lastBelowMAStock {
			if open := uc.cc.GetHistoryOpen(s.Stock.Number, uc.tradeDay.GetStockTradeDay().TradeDay); open != 0 {
				if open > s.QuaterMA {
					uc.rebornMap[s.Date] = append(uc.rebornMap[s.Date], *s.Stock)
				}
				delete(uc.lastBelowMAStock, s.Stock.Number)
			}
		}
	}
	uc.rebornLock.Unlock()
	return uc.rebornMap
}

func (uc *AnalyzeUseCase) findBelowQuaterMATargets(targetArr []*entity.StockTarget) {
	defer uc.rebornLock.Unlock()
	uc.rebornLock.Lock()
	uc.targetArr = append(uc.targetArr, targetArr...)

	for _, t := range targetArr {
		maMap, err := uc.repo.QueryAllQuaterMAByStockNum(context.Background(), t.StockNum)
		if err != nil {
			uc.logger.Fatal(err)
		}

		for _, ma := range maMap {
			tmp := ma
			if close := uc.cc.GetHistoryClose(ma.StockNum, ma.Date); close != 0 && close-ma.QuaterMA > 0 {
				continue
			}
			if nextTradeDay := uc.tradeDay.GetAbsNextTradeDayTime(ma.Date); nextTradeDay.Equal(uc.tradeDay.GetStockTradeDay().TradeDay) {
				uc.lastBelowMAStock[tmp.StockNum] = tmp
			} else if nextOpen := uc.cc.GetHistoryOpen(ma.StockNum, nextTradeDay); nextOpen != 0 && nextOpen-ma.QuaterMA > 0 {
				uc.rebornMap[ma.Date] = append(uc.rebornMap[ma.Date], *tmp.Stock)
			}
		}
	}

	uc.bus.PublishTopicEvent(topicSubscribeStockTickTargets, targetArr)
	uc.logger.Info("Find below quaterMA targets done")
}
