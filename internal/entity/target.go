package entity

import "time"

// Target -.
type Target struct {
	Stock    *Stock    `json:"stock"`
	StockNum string    `json:"stock_num"`
	TradeDay time.Time `json:"trade_day"`

	Rank        int   `json:"rank"`
	Volume      int64 `json:"volume"`
	Subscribe   bool  `json:"subscribe"`
	RealTimeAdd bool  `json:"real_time_add"`
}
