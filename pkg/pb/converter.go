// Package pb package pb
package pb

import "toc-machine-trading/internal/entity"

// ToStockEntity -.
func (c *StockDetailMessage) ToStockEntity() *entity.Stock {
	return &entity.Stock{
		Number:    c.GetCode(),
		Name:      c.GetName(),
		Exchange:  c.GetExchange(),
		Category:  c.GetCategory(),
		DayTrade:  c.GetDayTrade() == "Yes",
		LastClose: c.GetReference(),
	}
}
