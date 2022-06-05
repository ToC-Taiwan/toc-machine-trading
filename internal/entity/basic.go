// Package entity package entity
package entity

import (
	"time"

	"toc-machine-trading/pkg/pb"
)

// Stock -.
type Stock struct {
	Number    string  `json:"number"`
	Name      string  `json:"name"`
	Exchange  string  `json:"exchange"`
	Category  string  `json:"category"`
	DayTrade  bool    `json:"day_trade"`
	LastClose float64 `json:"last_close"`
}

// CalendarDate CalendarDate
type CalendarDate struct {
	Date       time.Time `json:"date"`
	IsTradeDay bool      `json:"is_trade_day"`
}

// FromProto -.
func (c *Stock) FromProto(data *pb.StockDetailMessage) *Stock {
	c.Number = data.GetCode()
	c.Name = data.GetName()
	c.Exchange = data.GetExchange()
	c.Category = data.GetCategory()
	c.DayTrade = data.GetDayTrade() == "Yes"
	c.LastClose = data.GetReference()
	return c
}
