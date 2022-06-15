package entity

import "time"

// Target -.
type Target struct {
	ID          int64     `json:"id"`
	Rank        int       `json:"rank"`
	StockNum    string    `json:"stock_num"`
	Volume      int64     `json:"volume"`
	Subscribe   bool      `json:"subscribe"`
	RealTimeAdd bool      `json:"real_time_add"`
	TradeDay    time.Time `json:"trade_day"`
	Stock       *Stock    `json:"stock"`
}
