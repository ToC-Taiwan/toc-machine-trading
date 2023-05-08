// Package cache package cache
package cache

import (
	"fmt"
	"time"

	"tmt/internal/entity"
	"tmt/pkg/cache"
)

const (
	cacheCatagoryBasic int64 = iota + 1
	cacheCatagoryStockDetail
	cacheCatagoryFutureDetail
	cacheCatagoryHistoryOpen
	cacheCatagoryHistoryClose
	cacheCatagoryHistoryTickAnalyze
	cacheCatagoryHistoryTickArr
	cacheCatagoryDayKbar
)

const (
	cacheStaticIndexTargets string = "targets"
)

var singleton *Cache

type Cache struct {
	*cache.Cache
}

func Get() *Cache {
	if singleton != nil {
		return singleton
	}

	singleton = &Cache{
		Cache: cache.New(),
	}

	return singleton
}

func (c *Cache) key(category int64, index ...string) string {
	if len(index) == 0 {
		panic("index is empty")
	}

	key := fmt.Sprintf("%d", category)
	for _, v := range index {
		key = fmt.Sprintf("%s:%s", key, v)
	}
	return key
}

func (c *Cache) SetStockDetail(stock *entity.Stock) {
	c.Set(c.key(cacheCatagoryStockDetail, stock.Number), stock)
}

func (c *Cache) GetStockDetail(stockNum string) *entity.Stock {
	if value, ok := c.Get(c.key(cacheCatagoryStockDetail, stockNum)); ok {
		return value.(*entity.Stock)
	}
	return nil
}

func (c *Cache) GetAllStockDetail() map[string]*entity.Stock {
	result := make(map[string]*entity.Stock)
	for k, v := range c.GetAll(cacheCatagoryStockDetail) {
		result[k] = v.(*entity.Stock)
	}
	return result
}

func (c *Cache) SetFutureDetail(future *entity.Future) {
	c.Set(c.key(cacheCatagoryFutureDetail, future.Code), future)
}

func (c *Cache) GetFutureDetail(code string) *entity.Future {
	if value, ok := c.Get(c.key(cacheCatagoryFutureDetail, code)); ok {
		return value.(*entity.Future)
	}
	return nil
}

func (c *Cache) GetAllFutureDetail() map[string]*entity.Future {
	result := make(map[string]*entity.Future)
	for k, v := range c.GetAll(cacheCatagoryFutureDetail) {
		result[k] = v.(*entity.Future)
	}
	return result
}

func (c *Cache) SetHistoryOpen(stockNum string, date time.Time, open float64) {
	c.Set(c.key(cacheCatagoryHistoryOpen, stockNum, date.Format("20060102")), open)
}

func (c *Cache) GetHistoryOpen(stockNum string, date time.Time) float64 {
	if value, ok := c.Get(c.key(cacheCatagoryHistoryOpen, stockNum, date.Format("20060102"))); ok {
		return value.(float64)
	}
	return 0
}

func (c *Cache) SetHistoryClose(stockNum string, date time.Time, close float64) {
	c.Set(c.key(cacheCatagoryHistoryClose, stockNum, date.Format("20060102")), close)
}

func (c *Cache) GetHistoryClose(stockNum string, date time.Time) float64 {
	if value, ok := c.Get(c.key(cacheCatagoryHistoryClose, stockNum, date.Format("20060102"))); ok {
		return value.(float64)
	}
	return 0
}

func (c *Cache) AppendStockTargets(targets []*entity.StockTarget) {
	original := c.GetStockTargets()
	original = append(original, targets...)
	c.setStockTargets(original)
}

func (c *Cache) setStockTargets(targets []*entity.StockTarget) {
	c.Set(c.key(cacheCatagoryBasic, cacheStaticIndexTargets), targets)
}

func (c *Cache) GetStockTargets() []*entity.StockTarget {
	if value, ok := c.Get(c.key(cacheCatagoryBasic, cacheStaticIndexTargets)); ok {
		return value.([]*entity.StockTarget)
	}
	return []*entity.StockTarget{}
}

func (c *Cache) GetHistoryTickAnalyze(stockNum string) []int64 {
	if value, ok := c.Get(c.key(cacheCatagoryHistoryTickAnalyze, stockNum)); ok {
		return value.([]int64)
	}
	return []int64{}
}

func (c *Cache) setHistoryTickAnalyze(stockNum string, arr []int64) {
	c.Set(c.key(cacheCatagoryHistoryTickAnalyze, stockNum), arr)
}

func (c *Cache) AppendHistoryTickAnalyze(stockNum string, arr []int64) {
	original := c.GetHistoryTickAnalyze(stockNum)
	original = append(original, arr...)
	c.setHistoryTickAnalyze(stockNum, original)
}

func (c *Cache) SetDaykbar(stockNum string, date time.Time, daykbar *entity.StockHistoryKbar) {
	c.Set(c.key(cacheCatagoryDayKbar, stockNum, date.Format("20060102")), daykbar)
}

func (c *Cache) GetDaykbar(stockNum string, date time.Time) *entity.StockHistoryKbar {
	if value, ok := c.Get(c.key(cacheCatagoryDayKbar, stockNum, date.Format("20060102"))); ok {
		return value.(*entity.StockHistoryKbar)
	}
	return nil
}

func (c *Cache) SetHistoryTickArr(stockNum string, date time.Time, tickArr []*entity.StockHistoryTick) {
	c.Set(c.key(cacheCatagoryHistoryTickArr, stockNum, date.Format("20060102")), tickArr)
}

func (c *Cache) GetHistoryTickArr(stockNum string, date time.Time) []*entity.StockHistoryTick {
	if value, ok := c.Get(c.key(cacheCatagoryHistoryTickArr, stockNum, date.Format("20060102"))); ok {
		return value.([]*entity.StockHistoryTick)
	}
	return nil
}
