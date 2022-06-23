package usecase

import (
	"fmt"
	"time"

	"toc-machine-trading/pkg/cache"
	"toc-machine-trading/pkg/global"
)

const (
	cacheCatagoryBasic cache.Category = "basic"

	cacheCatagoryOrder        cache.Category = "order"
	cacheCatagoryHistoryOpen  cache.Category = "history_open"
	cacheCatagoryHistoryClose cache.Category = "history_close"
	cacheCatagoryBiasRate     cache.Category = "bias_rate"
	cacheCatagoryQuaterMA     cache.Category = "quater_ma"
)

const (
	// no variable id
	cacheIDCalendar  string = "calendar"
	cacheIDBasicInfo string = "basic_info"
	cacheIDTargets   string = "targets"

	// with variable id
	cacheIDStockNum string = "stock_num"
	cacheIDOrderID  string = "order_id"
)

func (c *GlobalCache) generateID(opt ...any) string {
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

// Basic
//
func (c *GlobalCache) stockDetailKey(stockNum string) cache.Key {
	return cache.Key{
		Category: cacheCatagoryBasic,
		ID:       c.generateID(cacheIDStockNum, stockNum),
	}
}

func (c *GlobalCache) calendarKey() cache.Key {
	return cache.Key{
		Category: cacheCatagoryBasic,
		ID:       cacheIDCalendar,
	}
}

func (c *GlobalCache) basicInfoKey() cache.Key {
	return cache.Key{
		Category: cacheCatagoryBasic,
		ID:       cacheIDBasicInfo,
	}
}

func (c *GlobalCache) targetsKey() cache.Key {
	return cache.Key{
		Category: cacheCatagoryBasic,
		ID:       cacheIDTargets,
	}
}

// Order
//
func (c *GlobalCache) orderKey(orderID string) cache.Key {
	return cache.Key{
		Category: cacheCatagoryOrder,
		ID:       c.generateID(cacheIDOrderID, orderID),
	}
}

// HistoryOpen
//
func (c *GlobalCache) historyOpenKey(stockNum string, date time.Time) cache.Key {
	return cache.Key{
		Category: cacheCatagoryHistoryOpen,
		ID:       c.generateID(cacheIDStockNum, stockNum, date.Format(global.ShortTimeLayout)),
	}
}

// HistoryClose
//
func (c *GlobalCache) historyCloseKey(stockNum string, date time.Time) cache.Key {
	return cache.Key{
		Category: cacheCatagoryHistoryClose,
		ID:       c.generateID(cacheIDStockNum, stockNum, date.Format(global.ShortTimeLayout)),
	}
}

// BiasRate
//
func (c *GlobalCache) biasRateKey(stockNum string) cache.Key {
	return cache.Key{
		Category: cacheCatagoryBiasRate,
		ID:       c.generateID(cacheIDStockNum, stockNum),
	}
}
