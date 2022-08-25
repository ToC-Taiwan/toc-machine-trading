package usecase

import (
	"sync"

	"tmt/internal/entity"
	"tmt/pkg/config"
)

// TargetFilter -.
type TargetFilter struct {
	blackCategory map[string]struct{}
	blackStock    map[string]struct{}
	priceLimit    []config.PriceLimit
	realTimeRank  int64
	volumeLimit   int64

	mutex sync.RWMutex
}

// NewTargetFilter -.
func NewTargetFilter(cond config.TargetCond) *TargetFilter {
	blackCategoryMap := make(map[string]struct{})
	for _, category := range cond.BlackCategory {
		blackCategoryMap[category] = struct{}{}
	}

	blackStockMap := make(map[string]struct{})
	for _, stockNum := range cond.BlackStock {
		blackStockMap[stockNum] = struct{}{}
	}

	return &TargetFilter{
		blackCategory: blackCategoryMap,
		blackStock:    blackStockMap,
		priceLimit:    cond.PriceLimit,
		realTimeRank:  cond.RealTimeRank,
		volumeLimit:   cond.LimitVolume,
	}
}

func (t *TargetFilter) isTarget(stock *entity.Stock, close float64) bool {
	defer t.mutex.RUnlock()
	t.mutex.RLock()

	if _, ok := t.blackCategory[stock.Category]; ok {
		return false
	}

	if _, ok := t.blackStock[stock.Number]; ok {
		return false
	}

	for _, c := range t.priceLimit {
		if close >= c.Low && close < c.High {
			return true
		}
	}
	return false
}

func (t *TargetFilter) checkVolume(volume int64) bool {
	return volume > t.volumeLimit
}
