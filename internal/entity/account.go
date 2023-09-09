package entity

import "time"

const (
	// BankIDSinopac is the bank id of Sinopac
	BankIDSinopac = iota + 1
	// BankIDFugle is the bank id of Fugle
	BankIDFugle
)

type AccountBalance struct {
	ID              int       `json:"id" yaml:"id"`
	Date            time.Time `json:"date" yaml:"date"`
	Balance         float64   `json:"balance" yaml:"balance"`
	TodayMargin     float64   `json:"today_margin" yaml:"today_margin"`
	AvailableMargin float64   `json:"available_margin" yaml:"available_margin"`
	YesterdayMargin float64   `json:"yesterday_margin" yaml:"yesterday_margin"`
	RiskIndicator   float64   `json:"risk_indicator" yaml:"risk_indicator"`
	BankID          int       `json:"bank_id" yaml:"bank_id"`
}

type Settlement struct {
	Date    time.Time `json:"date" yaml:"date"`
	Sinopac float64   `json:"sinopac" yaml:"sinopac"`
	Fugle   float64   `json:"fugle" yaml:"fugle"`
}
