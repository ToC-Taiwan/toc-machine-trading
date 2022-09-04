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

// Future -.
type Future struct {
	Code           string    `json:"code"`
	Symbol         string    `json:"symbol"`
	Name           string    `json:"name"`
	Category       string    `json:"category"`
	DeliveryMonth  string    `json:"delivery_month"`
	DeliveryDate   string    `json:"delivery_date"`
	UnderlyingKind string    `json:"underlying_kind"`
	Unit           int64     `json:"unit"`
	LimitUp        float64   `json:"limit_up"`
	LimitDown      float64   `json:"limit_down"`
	Reference      float64   `json:"reference"`
	UpdateDate     time.Time `json:"update_date"`
}

const (
	// DayTradeYes -.
	DayTradeYes string = "Yes"
	// DayTradeNo -.
	DayTradeNo string = "No"
	// DayTradeOnlyBuy -.
	DayTradeOnlyBuy string = "OnlyBuy"
)

// BasicInfo -.
type BasicInfo struct {
	TradeDay           time.Time
	LastTradeDay       time.Time
	BefroeLastTradeDay time.Time

	OpenTime       time.Time
	EndTime        time.Time
	TradeInEndTime time.Time

	HistoryCloseRange []time.Time
	HistoryKbarRange  []time.Time
	HistoryTickRange  []time.Time

	AllStocks  map[string]*Stock
	AllFutures map[string]*Future
}
