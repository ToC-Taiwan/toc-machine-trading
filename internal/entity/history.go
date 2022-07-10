package entity

import "time"

// HistoryClose -.
type HistoryClose struct {
	ID    int64     `json:"id"`
	Date  time.Time `json:"date"`
	Close float64   `json:"close"`

	StockNum string `json:"stock_num"`
	Stock    *Stock `json:"stock"`
}

// HistoryKbar -.
type HistoryKbar struct {
	ID       int64     `json:"id"`
	KbarTime time.Time `json:"kbar_time"`
	Open     float64   `json:"open"`
	High     float64   `json:"high"`
	Low      float64   `json:"low"`
	Close    float64   `json:"close"`
	Volume   int64     `json:"volume"`

	StockNum string `json:"stock_num"`
	Stock    *Stock `json:"stock"`
}

// HistoryTick -.
type HistoryTick struct {
	ID        int64     `json:"id"`
	TickTime  time.Time `json:"tick_time"`
	Close     float64   `json:"close"`
	TickType  int64     `json:"tick_type"`
	Volume    int64     `json:"volume"`
	BidPrice  float64   `json:"bid_price"`
	BidVolume int64     `json:"bid_volume"`
	AskPrice  float64   `json:"ask_price"`
	AskVolume int64     `json:"ask_volume"`

	StockNum string `json:"stock_num"`
	Stock    *Stock `json:"stock"`
}

// HistoryAnalyze -.
type HistoryAnalyze struct {
	ID       int64     `json:"id"`
	Date     time.Time `json:"date"`
	QuaterMA float64   `json:"quater_ma"`

	StockNum string `json:"stock_num"`
	Stock    *Stock `json:"stock"`
}
