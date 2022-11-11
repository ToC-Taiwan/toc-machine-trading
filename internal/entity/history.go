package entity

import "time"

// StockHistoryTick -.
type StockHistoryTick struct {
	StockNum string `json:"stock_num"`
	Stock    *Stock `json:"stock"`
	HistoryTickBase
}

// FutureHistoryTick -.
type FutureHistoryTick struct {
	Code   string  `json:"code"`
	Future *Future `json:"future"`
	HistoryTickBase
}

// HistoryTickBase -.
type HistoryTickBase struct {
	ID        int64     `json:"id"`
	TickTime  time.Time `json:"tick_time"`
	Close     float64   `json:"close"`
	TickType  int64     `json:"tick_type"`
	Volume    int64     `json:"volume"`
	BidPrice  float64   `json:"bid_price"`
	BidVolume int64     `json:"bid_volume"`
	AskPrice  float64   `json:"ask_price"`
	AskVolume int64     `json:"ask_volume"`
}

type StockHistoryClose struct {
	StockNum string `json:"stock_num"`
	Stock    *Stock `json:"stock"`
	HistoryCloseBase
}

// FutureHistoryClose -.
type FutureHistoryClose struct {
	Code   string  `json:"code"`
	Future *Future `json:"future"`
	HistoryCloseBase
}

// HistoryCloseBase -.
type HistoryCloseBase struct {
	ID    int64     `json:"id"`
	Date  time.Time `json:"date"`
	Close float64   `json:"close"`
}

type StockHistoryKbar struct {
	StockNum string `json:"stock_num"`
	Stock    *Stock `json:"stock"`
	HistoryKbarBase
}

// HistoryKbarBase -.
type HistoryKbarBase struct {
	ID       int64     `json:"id"`
	KbarTime time.Time `json:"kbar_time"`
	Open     float64   `json:"open"`
	High     float64   `json:"high"`
	Low      float64   `json:"low"`
	Close    float64   `json:"close"`
	Volume   int64     `json:"volume"`
}

// StockHistoryAnalyze -.
type StockHistoryAnalyze struct {
	StockNum string `json:"stock_num"`
	Stock    *Stock `json:"stock"`
	HistoryAnalyzeBase
}

type HistoryAnalyzeBase struct {
	ID       int64     `json:"id"`
	Date     time.Time `json:"date"`
	QuaterMA float64   `json:"quater_ma"`
}
