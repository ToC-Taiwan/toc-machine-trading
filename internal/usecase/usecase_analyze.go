package usecase

import (
	"context"
	"sync"
	"time"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/internal/usecase/repo"
)

// AnalyzeUseCase -.
type AnalyzeUseCase struct {
	repo HistoryRepo

	lastBelowMAStock map[string]*entity.HistoryAnalyze
	belowQuaterMap   map[time.Time][]entity.Stock
	belowQuaterLock  sync.Mutex
}

// NewAnalyze -.
func NewAnalyze(r *repo.HistoryRepo) *AnalyzeUseCase {
	uc := &AnalyzeUseCase{repo: r}

	uc.lastBelowMAStock = make(map[string]*entity.HistoryAnalyze)
	uc.belowQuaterMap = make(map[time.Time][]entity.Stock)

	bus.SubscribeTopic(topicAnalyzeTargets, uc.AnalyzeAll)
	return uc
}

// AnalyzeAll -.
func (uc *AnalyzeUseCase) AnalyzeAll(ctx context.Context, targetArr []*entity.Target) {
	uc.findBelowQuaterMATargets(ctx, targetArr)

	bus.PublishTopicEvent(topicStreamTargets, ctx, targetArr)
}

// GetBelowQuaterMap GetBelowQuaterMap
func (uc *AnalyzeUseCase) GetBelowQuaterMap(ctx context.Context) map[time.Time][]entity.Stock {
	uc.belowQuaterLock.Lock()
	basicInfo := cc.GetBasicInfo()
	if len(uc.lastBelowMAStock) != 0 {
		for _, s := range uc.lastBelowMAStock {
			if open := cc.GetHistoryOpen(s.Stock.Number, basicInfo.TradeDay); open != 0 {
				if open > s.QuaterMA {
					uc.belowQuaterMap[s.Date] = append(uc.belowQuaterMap[s.Date], *s.Stock)
				}
				delete(uc.lastBelowMAStock, s.Stock.Number)
			}
		}
	}
	uc.belowQuaterLock.Unlock()
	return uc.belowQuaterMap
}

func (uc *AnalyzeUseCase) findBelowQuaterMATargets(ctx context.Context, targetArr []*entity.Target) {
	log.Info("findBelowQuaterMATargets")
	defer uc.belowQuaterLock.Unlock()
	uc.belowQuaterLock.Lock()
	for _, t := range targetArr {
		maMap, err := uc.repo.QueryAllQuaterMAByStockNum(ctx, t.StockNum)
		if err != nil {
			log.Panic(err)
		}

		basicInfo := cc.GetBasicInfo()
		for _, ma := range maMap {
			tmp := ma
			if close := cc.GetHistoryClose(ma.StockNum, ma.Date); close != 0 && close-ma.QuaterMA > 0 {
				continue
			}
			if nextTradeDay := getAbsNextTradeDayTime(ma.Date); nextTradeDay.Equal(basicInfo.TradeDay) {
				uc.lastBelowMAStock[tmp.StockNum] = tmp
			} else if nextOpen := cc.GetHistoryOpen(ma.StockNum, nextTradeDay); nextOpen != 0 && nextOpen-ma.QuaterMA > 0 {
				uc.belowQuaterMap[ma.Date] = append(uc.belowQuaterMap[ma.Date], *tmp.Stock)
			}
		}
	}
}
