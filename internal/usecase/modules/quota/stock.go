// Package quota package quota
package quota

import (
	"math"

	"tmt/cmd/config"
	"tmt/internal/entity"
)

// Quota -.
type Quota struct {
	quota         int64
	tradeTaxRatio float64
	tradeFeeRatio float64
	feeDiscount   float64
}

// NewQuota -.
func NewQuota(cfg config.Quota) *Quota {
	return &Quota{
		quota:         cfg.TradeQuota,
		tradeTaxRatio: cfg.TradeTaxRatio,
		tradeFeeRatio: cfg.TradeFeeRatio,
		feeDiscount:   cfg.FeeDiscount,
	}
}

// GetCurrentQuota -.
func (q *Quota) GetCurrentQuota() int64 {
	return q.quota
}

// CosumeQuota -.
func (q *Quota) CosumeQuota(t int64) {
	q.quota -= t
}

// BackQuota -.
func (q *Quota) BackQuota(t int64) {
	q.quota += t
}

// IsEnough -.
func (q *Quota) IsEnough(t int64) bool {
	return q.quota >= t
}

// CalculateOriginalOrderCost -.
func (q *Quota) CalculateOriginalOrderCost(order *entity.StockOrder) int64 {
	if order.Action == entity.ActionBuy || order.Action == entity.ActionSellFirst {
		return int64(math.Ceil(order.Price * float64(order.Quantity) * 1000))
	}
	return 0
}

// GetStockBuyCost -.
func (q *Quota) GetStockBuyCost(price float64, qty int64) int64 {
	return int64(math.Ceil(price*float64(qty)*1000) + math.Floor(price*float64(qty)*1000*q.tradeFeeRatio))
}

// GetStockSellCost -.
func (q *Quota) GetStockSellCost(price float64, qty int64) int64 {
	return int64(math.Ceil(price*float64(qty)*1000) - math.Floor(price*float64(qty)*1000*q.tradeFeeRatio) - math.Floor(price*float64(qty)*1000*q.tradeTaxRatio))
}

// GetStockTradeFeeDiscount -.
func (q *Quota) GetStockTradeFeeDiscount(price float64, qty int64) int64 {
	return int64(math.Floor(price*float64(qty)*1000*q.tradeFeeRatio) * (1 - q.feeDiscount))
}
