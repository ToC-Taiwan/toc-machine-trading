package usecase

import (
	"time"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/pkg/cache"
)

// GlobalCache -.
type GlobalCache struct {
	*cache.Cache
}

// NewCache -.
func NewCache() *GlobalCache {
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

// SetHistoryClose -.
func (c *GlobalCache) SetHistoryClose(stockNum string, date time.Time, close float64) {
	c.Set(c.historyCloseKey(stockNum, date), close)
}

// GetHistoryClose -.
func (c *GlobalCache) GetHistoryClose(stockNum string, date time.Time) float64 {
	if value, ok := c.Get(c.historyCloseKey(stockNum, date)); ok {
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

// SetTargets -.
func (c *GlobalCache) SetTargets(targets []*entity.Target) {
	c.Set(c.targetsKey(), targets)
}

// GetTargets -.
func (c *GlobalCache) GetTargets() []*entity.Target {
	if value, ok := c.Get(c.targetsKey()); ok {
		return value.([]*entity.Target)
	}
	return []*entity.Target{}
}

// GetHistoryTickAnalyze -.
func (c *GlobalCache) GetHistoryTickAnalyze(stockNum string) []int64 {
	if value, ok := c.Get(c.historyTickAnalyzeKey(stockNum)); ok {
		return value.([]int64)
	}
	return []int64{}
}

func (c *GlobalCache) setHistoryTickAnalyze(stockNum string, arr []int64) {
	c.Set(c.historyTickAnalyzeKey(stockNum), arr)
}

// AppendHistoryTickAnalyze -.
func (c *GlobalCache) AppendHistoryTickAnalyze(stockNum string, arr []int64) {
	original := c.GetHistoryTickAnalyze(stockNum)
	original = append(original, arr...)
	c.setHistoryTickAnalyze(stockNum, original)
}

// SetDaykbar -.
func (c *GlobalCache) SetDaykbar(stockNum string, date time.Time, daykbar *entity.HistoryKbar) {
	c.Set(c.dayKbarKey(stockNum, date), daykbar)
}

// GetDaykbar -.
func (c *GlobalCache) GetDaykbar(stockNum string, date time.Time) *entity.HistoryKbar {
	if value, ok := c.Get(c.dayKbarKey(stockNum, date)); ok {
		return value.(*entity.HistoryKbar)
	}
	return nil
}
