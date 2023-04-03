package entity

import "time"

const (
	// BankIDSinopac is the bank id of Sinopac
	BankIDSinopac = iota + 1
	// BankIDFugle is the bank id of Fugle
	BankIDFugle
)

type AccountBalance struct {
	ID              int       `json:"id,omitempty" yaml:"id"`
	Date            time.Time `json:"date,omitempty" yaml:"date"`
	Balance         float64   `json:"balance,omitempty" yaml:"balance"`
	TodayMargin     float64   `json:"today_margin,omitempty" yaml:"today_margin"`
	AvailableMargin float64   `json:"available_margin,omitempty" yaml:"available_margin"`
	YesterdayMargin float64   `json:"yesterday_margin,omitempty" yaml:"yesterday_margin"`
	RiskIndicator   float64   `json:"risk_indicator,omitempty" yaml:"risk_indicator"`
	BankID          int       `json:"bank_id,omitempty" yaml:"bank_id"`
}

type Settlement struct {
	Date    time.Time `json:"date,omitempty" yaml:"date"`
	Sinopac float64   `json:"sinopac,omitempty" yaml:"sinopac"`
	Fugle   float64   `json:"fugle,omitempty" yaml:"fugle"`
}
