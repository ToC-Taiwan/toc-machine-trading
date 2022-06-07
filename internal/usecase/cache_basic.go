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

// CacheSetStockDetail CacheSetStockDetail
func CacheSetStockDetail(stock *entity.Stock) {
	key := cache.Key{
		Type: keyTypeStockDetail,
		Name: fmt.Sprintf("%s:%s", keyTypeStockDetail, stock.Number),
	}
	cache.Set(key, stock)
}

// CacheGetStockDetail CacheGetStockDetail
func CacheGetStockDetail(stockNum string) *entity.Stock {
	key := cache.Key{
		Type: keyTypeStockDetail,
		Name: fmt.Sprintf("%s:%s", keyTypeStockDetail, stockNum),
	}
	if value, ok := cache.Get(key); ok {
		return value.(*entity.Stock)
	}
	return nil
}

// CacheSetCalendar CacheSetCalendar
func CacheSetCalendar(calendar map[time.Time]bool) {
	key := cache.Key{
		Type: keyTypeCalendar,
		Name: keyTypeCalendar,
	}
	cache.Set(key, calendar)
}

// CacheGetCalendar CacheGetCalendar
func CacheGetCalendar() map[time.Time]bool {
	key := cache.Key{
		Type: keyTypeCalendar,
		Name: keyTypeCalendar,
	}
	if value, ok := cache.Get(key); ok {
		return value.(map[time.Time]bool)
	}
	return nil
}

// CacheSetTradeDay CacheSetTradeDay
func CacheSetTradeDay(tradeDay time.Time) {
	key := cache.Key{
		Type: keyTypeTradeDay,
		Name: keyTypeTradeDay,
	}
	cache.Set(key, tradeDay)
}

// CacheGetTradeDay CacheGetTradeDay
func CacheGetTradeDay() time.Time {
	key := cache.Key{
		Type: keyTypeTradeDay,
		Name: keyTypeTradeDay,
	}
	if value, ok := cache.Get(key); ok {
		return value.(time.Time)
	}
	return time.Time{}
}
