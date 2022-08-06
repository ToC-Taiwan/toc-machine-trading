package entity

import "time"

// Target -.
type Target struct {
	ID       int64     `json:"id"`
	Rank     int       `json:"rank"`
	Volume   int64     `json:"volume"`
	RealTime bool      `json:"real_time"`
	TradeDay time.Time `json:"trade_day"`

	StockNum string `json:"stock_num"`
	Stock    *Stock `json:"stock"`
}
