package usecase

import (
	"math"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/pkg/config"
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

// CalculateOrderCost -.
func (q *Quota) CalculateOrderCost(order *entity.Order) int64 {
	switch order.Action {
	case entity.ActionBuy, entity.ActionSellFirst:
		return order.Quantity * int64(order.Price) * 1000
	default:
		return 0
	}
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