// Package cache package cache
package cache

import (
	"time"

	"tmt/internal/entity"
	"tmt/pkg/cache"
)

var singleton *Cache

// Cache -.
type Cache struct {
	*cache.Cache
}

// Get -.
func Get() *Cache {
	if singleton != nil {
		return singleton
	}

	singleton = &Cache{
		Cache: cache.New(),
	}

	return singleton
}

// SetStockDetail -.
func (c *Cache) SetStockDetail(stock *entity.Stock) {
	c.Set(c.stockDetailKey(stock.Number), stock)
}

// GetStockDetail -.
func (c *Cache) GetStockDetail(stockNum string) *entity.Stock {
	if value, ok := c.Get(c.stockDetailKey(stockNum)); ok {
		return value.(*entity.Stock)
	}
	return nil
}

func (c *Cache) GetAllStockDetail() map[string]*entity.Stock {
	result := make(map[string]*entity.Stock)
	for k, v := range c.GetAll(c.allStockDetailKey()) {
		result[k] = v.(*entity.Stock)
	}
	return result
}

// SetFutureDetail -.
func (c *Cache) SetFutureDetail(future *entity.Future) {
	c.Set(c.futureDetailKey(future.Code), future)
}

// GetFutureDetail -.
func (c *Cache) GetFutureDetail(code string) *entity.Future {
	if value, ok := c.Get(c.futureDetailKey(code)); ok {
		return value.(*entity.Future)
	}
	return nil
}

func (c *Cache) GetAllFutureDetail() map[string]*entity.Future {
	result := make(map[string]*entity.Future)
	for k, v := range c.GetAll(c.allFutureDetailKey()) {
		result[k] = v.(*entity.Future)
	}
	return result
}

// SetHistoryOpen -.
func (c *Cache) SetHistoryOpen(stockNum string, date time.Time, open float64) {
	c.Set(c.historyOpenKey(stockNum, date), open)
}

// GetHistoryOpen -.
func (c *Cache) GetHistoryOpen(stockNum string, date time.Time) float64 {
	if value, ok := c.Get(c.historyOpenKey(stockNum, date)); ok {
		return value.(float64)
	}
	return 0
}

// SetHistoryClose -.
func (c *Cache) SetHistoryClose(stockNum string, date time.Time, close float64) {
	c.Set(c.historyCloseKey(stockNum, date), close)
}

// GetHistoryClose -.
func (c *Cache) GetHistoryClose(stockNum string, date time.Time) float64 {
	if value, ok := c.Get(c.historyCloseKey(stockNum, date)); ok {
		return value.(float64)
	}
	return 0
}

// AppendStockTargets -.
func (c *Cache) AppendStockTargets(targets []*entity.StockTarget) {
	original := c.GetStockTargets()
	original = append(original, targets...)
	c.setStockTargets(original)
}

func (c *Cache) setStockTargets(targets []*entity.StockTarget) {
	c.Set(c.targetsKey(), targets)
}

// GetStockTargets -.
func (c *Cache) GetStockTargets() []*entity.StockTarget {
	if value, ok := c.Get(c.targetsKey()); ok {
		return value.([]*entity.StockTarget)
	}
	return []*entity.StockTarget{}
}

// GetHistoryTickAnalyze -.
func (c *Cache) GetHistoryTickAnalyze(stockNum string) []int64 {
	if value, ok := c.Get(c.historyTickAnalyzeKey(stockNum)); ok {
		return value.([]int64)
	}
	return []int64{}
}

func (c *Cache) setHistoryTickAnalyze(stockNum string, arr []int64) {
	c.Set(c.historyTickAnalyzeKey(stockNum), arr)
}

// AppendHistoryTickAnalyze -.
func (c *Cache) AppendHistoryTickAnalyze(stockNum string, arr []int64) {
	original := c.GetHistoryTickAnalyze(stockNum)
	original = append(original, arr...)
	c.setHistoryTickAnalyze(stockNum, original)
}

// SetDaykbar -.
func (c *Cache) SetDaykbar(stockNum string, date time.Time, daykbar *entity.StockHistoryKbar) {
	c.Set(c.dayKbarKey(stockNum, date), daykbar)
}

// GetDaykbar -.
func (c *Cache) GetDaykbar(stockNum string, date time.Time) *entity.StockHistoryKbar {
	if value, ok := c.Get(c.dayKbarKey(stockNum, date)); ok {
		return value.(*entity.StockHistoryKbar)
	}
	return nil
}

// SetHistoryTickArr -.
func (c *Cache) SetHistoryTickArr(stockNum string, date time.Time, tickArr []*entity.StockHistoryTick) {
	c.Set(c.historyTickArrKey(stockNum, date), tickArr)
}

// GetHistoryTickArr -.
func (c *Cache) GetHistoryTickArr(stockNum string, date time.Time) []*entity.StockHistoryTick {
	if value, ok := c.Get(c.historyTickArrKey(stockNum, date)); ok {
		return value.([]*entity.StockHistoryTick)
	}
	return nil
}
