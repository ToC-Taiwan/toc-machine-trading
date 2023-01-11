// Package target package target
package target

import (
	"sync"

	"tmt/cmd/config"
	"tmt/internal/entity"
)

// Filter -.
type Filter struct {
	blackCategory map[string]struct{}
	blackStock    map[string]struct{}
	mutex         sync.RWMutex

	priceLimit  []config.PriceLimit
	volumeLimit int64

	RealTimeRank int64
}

// NewFilter -.
func NewFilter(cond config.TargetCond) *Filter {
	blackCategoryMap := make(map[string]struct{})
	for _, category := range cond.BlackCategory {
		blackCategoryMap[category] = struct{}{}
	}

	blackStockMap := make(map[string]struct{})
	for _, stockNum := range cond.BlackStock {
		blackStockMap[stockNum] = struct{}{}
	}

	return &Filter{
		blackCategory: blackCategoryMap,
		blackStock:    blackStockMap,
		priceLimit:    cond.PriceLimit,
		RealTimeRank:  cond.RealTimeRank,
		volumeLimit:   cond.LimitVolume,
	}
}

// IsTarget -.
func (t *Filter) IsTarget(stock *entity.Stock, close float64) bool {
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

// CheckVolume -.
func (t *Filter) CheckVolume(volume int64) bool {
	return volume > t.volumeLimit
}
