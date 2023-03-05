// Package cache package cache
package cache

import (
	"time"

	"tmt/global"
	"tmt/pkg/cache"
)

const (
	cacheCatagoryBasic              string = "basic"
	cacheCatagoryHistoryOpen        string = "history_open"
	cacheCatagoryHistoryClose       string = "history_close"
	cacheCatagoryHistoryTickAnalyze string = "history_tick_analyze"
	cacheCatagoryHistoryTickArr     string = "history_tick_arr"
	cacheCatagoryBiasRate           string = "bias_rate"
	cacheCatagoryDayKbar            string = "day_kbar"
)

const (
	cacheIndexBasicInfo string = "basic_info"
	cacheIndexTargets   string = "targets"

	cacheIndexStockNum   string = "stock_num"
	cacheIndexFutureCode string = "future_code"
	cacheIndexOrderID    string = "order_id"
)

func (c *Cache) basicInfoKey() *cache.Key {
	return cache.NewKey(cacheCatagoryBasic, cacheIndexBasicInfo)
}

func (c *Cache) targetsKey() *cache.Key {
	return cache.NewKey(cacheCatagoryBasic, cacheIndexTargets)
}

func (c *Cache) stockDetailKey(stockNum string) *cache.Key {
	return cache.NewKey(cacheCatagoryBasic, cacheIndexStockNum).ExtendIndex(stockNum)
}

func (c *Cache) futureDetailKey(code string) *cache.Key {
	return cache.NewKey(cacheCatagoryBasic, cacheIndexFutureCode).ExtendIndex(code)
}

func (c *Cache) historyOpenKey(stockNum string, date time.Time) *cache.Key {
	return cache.NewKey(cacheCatagoryHistoryOpen, cacheIndexStockNum).ExtendIndex(stockNum, date.Format(global.ShortTimeLayout))
}

func (c *Cache) historyCloseKey(stockNum string, date time.Time) *cache.Key {
	return cache.NewKey(cacheCatagoryHistoryClose, cacheIndexStockNum).ExtendIndex(stockNum, date.Format(global.ShortTimeLayout))
}

func (c *Cache) historyTickAnalyzeKey(stockNum string) *cache.Key {
	return cache.NewKey(cacheCatagoryHistoryTickAnalyze, cacheIndexStockNum).ExtendIndex(stockNum)
}

func (c *Cache) dayKbarKey(stockNum string, date time.Time) *cache.Key {
	return cache.NewKey(cacheCatagoryDayKbar, cacheIndexStockNum).ExtendIndex(stockNum, date.Format(global.ShortTimeLayout))
}

func (c *Cache) historyTickArrKey(stockNum string, date time.Time) *cache.Key {
	return cache.NewKey(cacheCatagoryHistoryTickArr, cacheIndexStockNum).ExtendIndex(stockNum, date.Format(global.ShortTimeLayout))
}
