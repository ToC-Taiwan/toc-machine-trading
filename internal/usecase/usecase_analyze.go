package usecase

import (
	"context"
	"sync"
	"time"

	"tmt/internal/entity"
	"tmt/internal/usecase/module/tradeday"
	"tmt/internal/usecase/topic"
)

// AnalyzeUseCase -.
type AnalyzeUseCase struct {
	repo      HistoryRepo
	targetArr []*entity.StockTarget

	lastBelowMAStock map[string]*entity.StockHistoryAnalyze
	rebornMap        map[time.Time][]entity.Stock
	rebornLock       sync.Mutex

	tradeDay *tradeday.TradeDay
}

// NewAnalyze -.
func NewAnalyze(r HistoryRepo) Analyze {
	uc := &AnalyzeUseCase{
		repo:             r,
		lastBelowMAStock: make(map[string]*entity.StockHistoryAnalyze),
		rebornMap:        make(map[time.Time][]entity.Stock),
		tradeDay:         tradeday.NewTradeDay(),
	}

	bus.SubscribeTopic(topic.TopicAnalyzeStockTargets, uc.findBelowQuaterMATargets)
	return uc
}

// GetRebornMap -.
func (uc *AnalyzeUseCase) GetRebornMap(ctx context.Context) map[time.Time][]entity.Stock {
	uc.rebornLock.Lock()
	basicInfo := cc.GetBasicInfo()
	if len(uc.lastBelowMAStock) != 0 {
		for _, s := range uc.lastBelowMAStock {
			if open := cc.GetHistoryOpen(s.Stock.Number, basicInfo.TradeDay); open != 0 {
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
			logger.Fatal(err)
		}

		basicInfo := cc.GetBasicInfo()
		for _, ma := range maMap {
			tmp := ma
			if close := cc.GetHistoryClose(ma.StockNum, ma.Date); close != 0 && close-ma.QuaterMA > 0 {
				continue
			}
			if nextTradeDay := uc.tradeDay.GetAbsNextTradeDayTime(ma.Date); nextTradeDay.Equal(basicInfo.TradeDay) {
				uc.lastBelowMAStock[tmp.StockNum] = tmp
			} else if nextOpen := cc.GetHistoryOpen(ma.StockNum, nextTradeDay); nextOpen != 0 && nextOpen-ma.QuaterMA > 0 {
				uc.rebornMap[ma.Date] = append(uc.rebornMap[ma.Date], *tmp.Stock)
			}
		}
	}
	bus.PublishTopicEvent(topic.TopicSubscribeStockTickTargets, targetArr)
	logger.Info("Find below quaterMA targets done")
}
