package entity

import "time"

// HistoryClose HistoryClose
type HistoryClose struct {
	ID       int64     `json:"id"`
	Date     time.Time `json:"date"`
	StockNum string    `json:"stock_num"`
	Close    float64   `json:"close"`
	Stock    *Stock    `json:"stock"`
}

// HistoryKbar HistoryKbar
type HistoryKbar struct {
	ID       int64     `json:"id"`
	StockNum string    `json:"stock_num"`
	KbarTime time.Time `json:"kbar_time"`
	Open     float64   `json:"open"`
	High     float64   `json:"high"`
	Low      float64   `json:"low"`
	Close    float64   `json:"close"`
	Volume   int64     `json:"volume"`
	Stock    *Stock    `json:"stock"`
}

// HistoryTick HistoryTick
type HistoryTick struct {
	ID        int64     `json:"id"`
	StockNum  string    `json:"stock_num"`
	TickTime  time.Time `json:"tick_time"`
	Close     float64   `json:"close"`
	TickType  int64     `json:"tick_type"`
	Volume    int64     `json:"volume"`
	BidPrice  float64   `json:"bid_price"`
	BidVolume int64     `json:"bid_volume"`
	AskPrice  float64   `json:"ask_price"`
	AskVolume int64     `json:"ask_volume"`
	Stock     *Stock    `json:"stock"`
}
