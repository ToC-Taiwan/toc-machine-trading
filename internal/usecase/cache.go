package usecase

import (
	"time"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/pkg/cache"
)

var cc = New()

// GlobalCache -.
type GlobalCache struct {
	*cache.Cache
}

// New -.
func New() *GlobalCache {
	return &GlobalCache{
		Cache: cache.New(),
	}
}

// SetStockDetail -.
func (c *GlobalCache) SetStockDetail(stock *entity.Stock) {
	c.Set(c.stockDetailKey(stock.Number), stock)
}

// GetStockDetail -.
func (c *GlobalCache) GetStockDetail(stockNum string) *entity.Stock {
	if value, ok := c.Get(c.stockDetailKey(stockNum)); ok {
		return value.(*entity.Stock)
	}
	return nil
}

// SetCalendar -.
func (c *GlobalCache) SetCalendar(calendar map[time.Time]bool) {
	c.Set(c.calendarKey(), calendar)
}

// GetCalendar -.
func (c *GlobalCache) GetCalendar() map[time.Time]bool {
	if value, ok := c.Get(c.calendarKey()); ok {
		return value.(map[time.Time]bool)
	}
	return nil
}

// SetBasicInfo -.
func (c *GlobalCache) SetBasicInfo(info *entity.BasicInfo) {
	c.Set(c.basicInfoKey(), info)
}

// GetBasicInfo -.
func (c *GlobalCache) GetBasicInfo() *entity.BasicInfo {
	if value, ok := c.Get(c.basicInfoKey()); ok {
		return value.(*entity.BasicInfo)
	}
	return nil
}

// SetOrderByOrderID -.
func (c *GlobalCache) SetOrderByOrderID(order *entity.Order) {
	c.Set(c.orderKey(order.OrderID), order)
}

// GetOrderByOrderID -.
func (c *GlobalCache) GetOrderByOrderID(orderID string) *entity.Order {
	if value, ok := c.Get(c.orderKey(orderID)); ok {
		return value.(*entity.Order)
	}
	return nil
}

// SetHistoryOpen -.
func (c *GlobalCache) SetHistoryOpen(stockNum string, date time.Time, open float64) {
	c.Set(c.historyOpenKey(stockNum, date), open)
}

// GetHistoryOpen -.
func (c *GlobalCache) GetHistoryOpen(stockNum string, date time.Time) float64 {
	if value, ok := c.Get(c.historyOpenKey(stockNum, date)); ok {
		return value.(float64)
	}
	return 0
}

// SetBiasRate -.
func (c *GlobalCache) SetBiasRate(stockNum string, biasRate float64) {
	c.Set(c.biasRateKey(stockNum), biasRate)
}

// GetBiasRate -.
func (c *GlobalCache) GetBiasRate(stockNum string) float64 {
	if value, ok := c.Get(c.biasRateKey(stockNum)); ok {
		return value.(float64)
	}
	return 0
}

// SetQuaterMA -.
func (c *GlobalCache) SetQuaterMA(stockNum string, quaterMA float64) {
	c.Set(c.quaterMAKey(stockNum), quaterMA)
}

// GetQuaterMA -.
func (c *GlobalCache) GetQuaterMA(stockNum string) float64 {
	if value, ok := c.Get(c.quaterMAKey(stockNum)); ok {
		return value.(float64)
	}
	return 0
}
