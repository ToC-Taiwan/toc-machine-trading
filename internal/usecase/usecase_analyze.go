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
	rebornMap        map[time.Time][]entity.Stock
	rebornLock       sync.Mutex
}

// NewAnalyze -.
func NewAnalyze(r *repo.HistoryRepo) *AnalyzeUseCase {
	uc := &AnalyzeUseCase{repo: r}

	uc.lastBelowMAStock = make(map[string]*entity.HistoryAnalyze)
	uc.rebornMap = make(map[time.Time][]entity.Stock)

	bus.SubscribeTopic(topicAnalyzeTargets, uc.AnalyzeAll)
	return uc
}

// AnalyzeAll -.
func (uc *AnalyzeUseCase) AnalyzeAll(ctx context.Context, targetArr []*entity.Target) {
	uc.findBelowQuaterMATargets(ctx, targetArr)

	bus.PublishTopicEvent(topicStreamTargets, ctx, targetArr)
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

func (uc *AnalyzeUseCase) findBelowQuaterMATargets(ctx context.Context, targetArr []*entity.Target) {
	defer uc.rebornLock.Unlock()
	uc.rebornLock.Lock()
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
				uc.rebornMap[ma.Date] = append(uc.rebornMap[ma.Date], *tmp.Stock)
			}
		}
	}
	log.Info("FindBelowQuaterMATargets Done")
}
