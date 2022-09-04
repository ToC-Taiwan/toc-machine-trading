package usecase

import (
	"time"

	"tmt/internal/entity"
	"tmt/pkg/cache"
)

// Cache -.
type Cache struct {
	*cache.Cache
}

// NewCache -.
func NewCache() *Cache {
	return &Cache{
		Cache: cache.New(),
	}
}

// SetStockDetail -.
func (c *Cache) SetStockDetail(stock *entity.Stock) {
	c.Set(c.stockDetailKey(stock.Number), stock)
}

// GetStockDetail -.
func (c *Cache) GetStockDetail(stockNum string) *entity.Stock {
	if value, ok := c.Get(c.stockDetailKey(stockNum)); ok {
		return value.(*entity.Stock)
	}
	return nil
}

// SetBasicInfo -.
func (c *Cache) SetBasicInfo(info *entity.BasicInfo) {
	c.Set(c.basicInfoKey(), info)
}

// GetBasicInfo -.
func (c *Cache) GetBasicInfo() *entity.BasicInfo {
	if value, ok := c.Get(c.basicInfoKey()); ok {
		return value.(*entity.BasicInfo)
	}
	return nil
}

// SetOrderByOrderID -.
func (c *Cache) SetOrderByOrderID(order *entity.Order) {
	c.Set(c.orderKey(order.OrderID), order)
}

// GetOrderByOrderID -.
func (c *Cache) GetOrderByOrderID(orderID string) *entity.Order {
	if value, ok := c.Get(c.orderKey(orderID)); ok {
		return value.(*entity.Order)
	}
	return nil
}

// SetFutureOrderByOrderID -.
func (c *Cache) SetFutureOrderByOrderID(order *entity.FutureOrder) {
	c.Set(c.futureOrderKey(order.OrderID), order)
}

// GetFutureOrderByOrderID -.
func (c *Cache) GetFutureOrderByOrderID(orderID string) *entity.FutureOrder {
	if value, ok := c.Get(c.futureOrderKey(orderID)); ok {
		return value.(*entity.FutureOrder)
	}
	return nil
}

// SetHistoryOpen -.
func (c *Cache) SetHistoryOpen(stockNum string, date time.Time, open float64) {
	c.Set(c.historyOpenKey(stockNum, date), open)
}

// GetHistoryOpen -.
func (c *Cache) GetHistoryOpen(stockNum string, date time.Time) float64 {
	if value, ok := c.Get(c.historyOpenKey(stockNum, date)); ok {
		return value.(float64)
	}
	return 0
}

// SetHistoryClose -.
func (c *Cache) SetHistoryClose(stockNum string, date time.Time, close float64) {
	c.Set(c.historyCloseKey(stockNum, date), close)
}

// GetHistoryClose -.
func (c *Cache) GetHistoryClose(stockNum string, date time.Time) float64 {
	if value, ok := c.Get(c.historyCloseKey(stockNum, date)); ok {
		return value.(float64)
	}
	return 0
}

// SetBiasRate -.
func (c *Cache) SetBiasRate(stockNum string, biasRate float64) {
	c.Set(c.biasRateKey(stockNum), biasRate)
}

// GetBiasRate -.
func (c *Cache) GetBiasRate(stockNum string) float64 {
	if value, ok := c.Get(c.biasRateKey(stockNum)); ok {
		return value.(float64)
	}
	return 0
}

// SetHighBiasRate -.
func (c *Cache) SetHighBiasRate(biasRate float64) {
	c.Set(c.highBiasRateKey(), biasRate)
}

// GetHighBiasRate -.
func (c *Cache) GetHighBiasRate() float64 {
	if value, ok := c.Get(c.highBiasRateKey()); ok {
		return value.(float64)
	}
	return 0
}

// SetLowBiasRate -.
func (c *Cache) SetLowBiasRate(biasRate float64) {
	c.Set(c.lowBiasRateKey(), biasRate)
}

// GetLowBiasRate -.
func (c *Cache) GetLowBiasRate() float64 {
	if value, ok := c.Get(c.lowBiasRateKey()); ok {
		return value.(float64)
	}
	return 0
}

// AppendTargets -.
func (c *Cache) AppendTargets(targets []*entity.Target) {
	original := c.GetTargets()
	original = append(original, targets...)
	c.setTargets(original)
}

func (c *Cache) setTargets(targets []*entity.Target) {
	c.Set(c.targetsKey(), targets)
}

// GetTargets -.
func (c *Cache) GetTargets() []*entity.Target {
	if value, ok := c.Get(c.targetsKey()); ok {
		return value.([]*entity.Target)
	}
	return []*entity.Target{}
}

// GetHistoryTickAnalyze -.
func (c *Cache) GetHistoryTickAnalyze(stockNum string) []int64 {
	if value, ok := c.Get(c.historyTickAnalyzeKey(stockNum)); ok {
		return value.([]int64)
	}
	return []int64{}
}

func (c *Cache) setHistoryTickAnalyze(stockNum string, arr []int64) {
	c.Set(c.historyTickAnalyzeKey(stockNum), arr)
}

// AppendHistoryTickAnalyze -.
func (c *Cache) AppendHistoryTickAnalyze(stockNum string, arr []int64) {
	original := c.GetHistoryTickAnalyze(stockNum)
	original = append(original, arr...)
	c.setHistoryTickAnalyze(stockNum, original)
}

// SetDaykbar -.
func (c *Cache) SetDaykbar(stockNum string, date time.Time, daykbar *entity.HistoryKbar) {
	c.Set(c.dayKbarKey(stockNum, date), daykbar)
}

// GetDaykbar -.
func (c *Cache) GetDaykbar(stockNum string, date time.Time) *entity.HistoryKbar {
	if value, ok := c.Get(c.dayKbarKey(stockNum, date)); ok {
		return value.(*entity.HistoryKbar)
	}
	return nil
}

// SetHistoryTickArr -.
func (c *Cache) SetHistoryTickArr(stockNum string, date time.Time, tickArr []*entity.HistoryTick) {
	c.Set(c.historyTickArrKey(stockNum, date), tickArr)
}

// GetHistoryTickArr -.
func (c *Cache) GetHistoryTickArr(stockNum string, date time.Time) []*entity.HistoryTick {
	if value, ok := c.Get(c.historyTickArrKey(stockNum, date)); ok {
		return value.([]*entity.HistoryTick)
	}
	return nil
}

// SetFutureGap -.
func (c *Cache) SetFutureGap(gap float64, date time.Time) {
	c.Set(c.futureGapKey(date), gap)
}

// GetFutureGap -.
func (c *Cache) GetFutureGap(date time.Time) float64 {
	if value, ok := c.Get(c.futureGapKey(date)); ok {
		return value.(float64)
	}
	return 0
}

// SetFutureHistoryTick -.
func (c *Cache) SetFutureHistoryTick(code string, tick *entity.HistoryTick) {
	c.Set(c.futureHistoryTickKey(code), tick)
}

// GetFutureHistoryTick -.
func (c *Cache) GetFutureHistoryTick(code string) *entity.HistoryTick {
	if value, ok := c.Get(c.futureHistoryTickKey(code)); ok {
		return value.(*entity.HistoryTick)
	}
	return nil
}
