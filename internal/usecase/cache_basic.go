package usecase

import (
	"fmt"
	"time"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/pkg/cache"
)

const (
	keyTypeStockDetail string = "stock_detail"

	keyTypeCalendar string = "calendar"
	keyTypeTradeDay string = "trade_day"
)

// SetStockDetail SetStockDetail
func SetStockDetail(stock *entity.Stock) {
	key := cache.Key{
		Type: keyTypeStockDetail,
		Name: fmt.Sprintf("%s:%s", keyTypeStockDetail, stock.Number),
	}
	cache.Set(key, stock)
}

// GetStockDetail GetStockDetail
func GetStockDetail(stockNum string) *entity.Stock {
	key := cache.Key{
		Type: keyTypeStockDetail,
		Name: fmt.Sprintf("%s:%s", keyTypeStockDetail, stockNum),
	}
	if value, ok := cache.Get(key); ok {
		return value.(*entity.Stock)
	}
	return nil
}

// SetCalendar SetCalendar
func SetCalendar(calendar map[time.Time]bool) {
	key := cache.Key{
		Type: keyTypeCalendar,
		Name: keyTypeCalendar,
	}
	cache.Set(key, calendar)
}

// GetCalendar GetCalendar
func GetCalendar() map[time.Time]bool {
	key := cache.Key{
		Type: keyTypeCalendar,
		Name: keyTypeCalendar,
	}
	if value, ok := cache.Get(key); ok {
		return value.(map[time.Time]bool)
	}
	return nil
}

// SetTradeDay SetTradeDay
func SetTradeDay(tradeDay time.Time) {
	key := cache.Key{
		Type: keyTypeTradeDay,
		Name: keyTypeTradeDay,
	}
	cache.Set(key, tradeDay)
}

// GetTradeDay GetTradeDay
func GetTradeDay() time.Time {
	key := cache.Key{
		Type: keyTypeTradeDay,
		Name: keyTypeTradeDay,
	}
	if value, ok := cache.Get(key); ok {
		return value.(time.Time)
	}
	return time.Time{}
}
