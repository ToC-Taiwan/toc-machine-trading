package usecase

import (
	"fmt"
	"time"

	"toc-machine-trading/pkg/cache"
	"toc-machine-trading/pkg/global"
)

const (
	cacheCatagoryBasic       cache.Category = "basic"
	cacheCatagoryOrder       cache.Category = "order"
	cacheCatagoryHistoryOpen cache.Category = "history_open"
)

const (
	// no variable id
	cacheIDCalendar  string = "calendar"
	cacheIDBasicInfo string = "basic_info"

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
