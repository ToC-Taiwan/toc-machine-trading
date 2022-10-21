// Package cache package cache
package cache

import (
	"fmt"
	"time"

	"tmt/global"
	"tmt/pkg/cache"
)

const (
	cacheCatagoryBasic cache.Category = "basic"

	cacheCatagoryOrder cache.Category = "order"

	cacheCatagoryFutureOrder cache.Category = "future_order"

	cacheCatagoryHistoryOpen cache.Category = "history_open"

	cacheCatagoryHistoryClose cache.Category = "history_close"

	cacheCatagoryHistoryTickAnalyze cache.Category = "history_tick_analyze"

	cacheCatagoryHistoryTickArr cache.Category = "history_tick_arr"

	cacheCatagoryBiasRate cache.Category = "bias_rate"

	cacheCatagoryDayKbar cache.Category = "day_kbar"
)

const (
	cacheIDBasicInfo string = "basic_info"

	cacheIDTargets string = "targets"

	cacheIDStockNum string = "stock_num"

	cacheIDOrderID string = "order_id"
)

func (c *Cache) generateID(opt ...any) string {
	total := len(opt)
	format := ""
	for i := 0; i < total; i++ {
		format += "%s"
		if i != total-1 {
			format += ":"
		}
	}
	return fmt.Sprintf(format, opt...)
}

func (c *Cache) basicInfoKey() cache.Key {
	return cache.Key{
		Category: cacheCatagoryBasic,
		ID:       cacheIDBasicInfo,
	}
}

func (c *Cache) targetsKey() cache.Key {
	return cache.Key{
		Category: cacheCatagoryBasic,
		ID:       cacheIDTargets,
	}
}

func (c *Cache) stockDetailKey(stockNum string) cache.Key {
	return cache.Key{
		Category: cacheCatagoryBasic,
		ID:       c.generateID(cacheIDStockNum, stockNum),
	}
}

func (c *Cache) orderKey(orderID string) cache.Key {
	return cache.Key{
		Category: cacheCatagoryOrder,
		ID:       c.generateID(cacheIDOrderID, orderID),
	}
}

func (c *Cache) futureOrderKey(orderID string) cache.Key {
	return cache.Key{
		Category: cacheCatagoryFutureOrder,
		ID:       c.generateID(cacheIDOrderID, orderID),
	}
}

func (c *Cache) historyOpenKey(stockNum string, date time.Time) cache.Key {
	return cache.Key{
		Category: cacheCatagoryHistoryOpen,
		ID:       c.generateID(cacheIDStockNum, stockNum, date.Format(global.ShortTimeLayout)),
	}
}

func (c *Cache) historyCloseKey(stockNum string, date time.Time) cache.Key {
	return cache.Key{
		Category: cacheCatagoryHistoryClose,
		ID:       c.generateID(cacheIDStockNum, stockNum, date.Format(global.ShortTimeLayout)),
	}
}

func (c *Cache) biasRateKey(stockNum string) cache.Key {
	return cache.Key{
		Category: cacheCatagoryBiasRate,
		ID:       c.generateID(cacheIDStockNum, stockNum),
	}
}

func (c *Cache) highBiasRateKey() cache.Key {
	return cache.Key{
		Category: cacheCatagoryBiasRate,
		ID:       "high",
	}
}

func (c *Cache) lowBiasRateKey() cache.Key {
	return cache.Key{
		Category: cacheCatagoryBiasRate,
		ID:       "low",
	}
}

func (c *Cache) historyTickAnalyzeKey(stockNum string) cache.Key {
	return cache.Key{
		Category: cacheCatagoryHistoryTickAnalyze,
		ID:       c.generateID(cacheIDStockNum, stockNum),
	}
}

func (c *Cache) dayKbarKey(stockNum string, date time.Time) cache.Key {
	return cache.Key{
		Category: cacheCatagoryDayKbar,
		ID:       c.generateID(cacheIDStockNum, stockNum, date.Format(global.ShortTimeLayout)),
	}
}

func (c *Cache) historyTickArrKey(stockNum string, date time.Time) cache.Key {
	return cache.Key{
		Category: cacheCatagoryHistoryTickArr,
		ID:       c.generateID(cacheIDStockNum, stockNum, date.Format(global.ShortTimeLayout)),
	}
}
