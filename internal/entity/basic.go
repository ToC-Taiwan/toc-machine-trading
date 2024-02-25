// Package entity package entity
package entity

import (
	"time"
)

type ShioajiUsage struct {
	Connections          int     `json:"connections"`
	TrafficUsage         float64 `json:"traffic_usage"`
	TrafficUsagePercents float64 `json:"traffic_usage_percents"`
}

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
	DeliveryDate   time.Time `json:"delivery_date"`
	UnderlyingKind string    `json:"underlying_kind"`
	Unit           int64     `json:"unit"`
	LimitUp        float64   `json:"limit_up"`
	LimitDown      float64   `json:"limit_down"`
	Reference      float64   `json:"reference"`
	UpdateDate     time.Time `json:"update_date"`
}

// Option -.
type Option struct {
	Code           string    `json:"code"`
	Symbol         string    `json:"symbol"`
	Name           string    `json:"name"`
	Category       string    `json:"category"`
	DeliveryMonth  string    `json:"delivery_month"`
	DeliveryDate   time.Time `json:"delivery_date"`
	StrikePrice    float64   `json:"strike_price"`
	OptionRight    string    `json:"option_right"`
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
