// Package quota package quota
package quota

import (
	"math"

	"tmt/cmd/config"
	"tmt/internal/entity"
)

const (
	stockTradeTaxRatio float64 = 0.0015
	stockTradeFeeRatio float64 = 0.001425

	futureTradeTaxRatio float64 = 0.00002
)

// Quota -.
type Quota struct {
	stockQuota       int64
	stockFeeDiscount float64
	futureTradeFee   int64
}

// NewQuota -.
func NewQuota(cfg config.Quota) *Quota {
	return &Quota{
		stockQuota:       cfg.StockTradeQuota,
		stockFeeDiscount: cfg.StockFeeDiscount,
		futureTradeFee:   cfg.FutureTradeFee,
	}
}

// GetCurrentQuota -.
func (q *Quota) GetCurrentQuota() int64 {
	return q.stockQuota
}

// CosumeQuota -.
func (q *Quota) CosumeQuota(t int64) {
	q.stockQuota -= t
}

// BackQuota -.
func (q *Quota) BackQuota(t int64) {
	q.stockQuota += t
}

// IsEnough -.
func (q *Quota) IsEnough(t int64) bool {
	return q.stockQuota >= t
}

// CalculateOriginalOrderCost -.
func (q *Quota) CalculateOriginalOrderCost(order *entity.StockOrder) int64 {
	if order.Action == entity.ActionBuy {
		return int64(math.Ceil(order.Price * float64(order.Quantity) * 1000))
	}
	return 0
}

// GetStockBuyCost -.
func (q *Quota) GetStockBuyCost(price float64, qty int64) int64 {
	base := price * float64(qty) * 1000
	return int64(math.Ceil(base) + math.Floor(base*stockTradeFeeRatio))
}

// GetStockSellCost -.
func (q *Quota) GetStockSellCost(price float64, qty int64) int64 {
	base := price * float64(qty) * 1000
	return int64(math.Ceil(base) - math.Floor(base*stockTradeFeeRatio) - math.Floor(base*stockTradeTaxRatio))
}

// GetStockTradeFeeDiscount -.
func (q *Quota) GetStockTradeFeeDiscount(price float64, qty int64) int64 {
	base := price * float64(qty) * 1000
	return int64(math.Floor(base*stockTradeFeeRatio) * (1 - q.stockFeeDiscount))
}

// GetFutureBuyCost -.
func (q *Quota) GetFutureBuyCost(price float64, qty int64) int64 {
	base := price * float64(qty) * 50
	return int64(math.Ceil(base)+math.Floor(base*futureTradeTaxRatio)) + q.futureTradeFee*qty
}

// GetFutureSellCost -.
func (q *Quota) GetFutureSellCost(price float64, qty int64) int64 {
	base := price * float64(qty) * 50
	return int64(math.Ceil(base)-math.Floor(base*futureTradeTaxRatio)) - q.futureTradeFee*qty
}
