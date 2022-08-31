package usecase

import (
	"fmt"
	"time"

	"tmt/pkg/cache"
	"tmt/pkg/global"
)

const (
	cacheCatagoryBasic cache.Category = "basic"

	cacheCatagoryOrder cache.Category = "order"

	cacheCatagoryRealTimeFutureTick cache.Category = "real_time_future_tick"
	cacheCatagoryFutureGap          cache.Category = "future_gap"

	cacheCatagoryHistoryOpen        cache.Category = "history_open"
	cacheCatagoryHistoryClose       cache.Category = "history_close"
	cacheCatagoryHistoryTickAnalyze cache.Category = "history_tick_analyze"
	cacheCatagoryHistoryTickArr     cache.Category = "history_tick_arr"
	cacheCatagoryBiasRate           cache.Category = "bias_rate"
	cacheCatagoryDayKbar            cache.Category = "day_kbar"
)

const (
	// no variable id
	cacheIDCalendar  string = "calendar"
	cacheIDBasicInfo string = "basic_info"
	cacheIDTargets   string = "targets"

	// with variable id
	cacheIDStockNum   string = "stock_num"
	cacheIDFutureCode string = "future_code"
	cacheIDOrderID    string = "order_id"
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

func (c *Cache) realTimeFutureTickKey(code string) cache.Key {
	return cache.Key{
		Category: cacheCatagoryRealTimeFutureTick,
		ID:       c.generateID(cacheIDFutureCode, code),
	}
}

func (c *Cache) futureGapKey(date time.Time) cache.Key {
	return cache.Key{
		Category: cacheCatagoryFutureGap,
		ID:       date.Format(global.ShortTimeLayout),
	}
}
