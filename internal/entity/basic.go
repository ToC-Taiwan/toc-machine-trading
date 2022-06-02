// Package entity package entity
package entity

// Stock -.
type Stock struct {
	ID        int64   `json:"id"`
	Number    string  `json:"number"`
	Name      string  `json:"name"`
	Exchange  string  `json:"exchange"`
	Category  string  `json:"category"`
	DayTrade  bool    `json:"day_trade"`
	LastClose float64 `json:"last_close"`
}
