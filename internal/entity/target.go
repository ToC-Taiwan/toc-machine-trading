package entity

import "time"

// Target -.
type Target struct {
	ID          int64     `json:"id"`
	Rank        int       `json:"rank"`
	Volume      int64     `json:"volume"`
	Subscribe   bool      `json:"subscribe"`
	RealTimeAdd bool      `json:"real_time_add"`
	TradeDay    time.Time `json:"trade_day"`

	StockNum string `json:"stock_num"`
	Stock    *Stock `json:"stock"`
}
