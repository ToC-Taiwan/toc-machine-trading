// Package entity package entity
package entity

import (
	"time"
)

// CalendarDate -.
type CalendarDate struct {
	Date       time.Time `json:"date"`
	IsTradeDay bool      `json:"is_trade_day"`
}

// Stock -.
type Stock struct {
	Number     string    `json:"number"`
	Name       string    `json:"name"`
	Exchange   string    `json:"exchange"`
	Category   string    `json:"category"`
	DayTrade   bool      `json:"day_trade"`
	LastClose  float64   `json:"last_close"`
	UpdateDate time.Time `json:"update_date"`
}

// BasicInfo -.
type BasicInfo struct {
	TradeDay     time.Time
	LastTradeDay time.Time

	OpenTime        time.Time
	EndTime         time.Time
	TradeInEndTime  time.Time
	TradeOutEndTime time.Time

	HistoryCloseRange []time.Time
	HistoryKbarRange  []time.Time
	HistoryTickRange  []time.Time
}
