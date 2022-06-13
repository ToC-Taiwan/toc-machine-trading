package entity

import (
	"time"
)

// SinopacEvent SinopacEvent
type SinopacEvent struct {
	Event     string    `json:"event"`
	EventCode int64     `json:"event_code"`
	Info      string    `json:"info"`
	Response  int64     `json:"response"`
	EventTime time.Time `json:"event_time"`
}

// RealTimeBidAsk -.
type RealTimeBidAsk struct {
	Stock    *Stock    `json:"stock"`
	StockNum string    `json:"stock_num"`
	TickTime time.Time `json:"tick_time"`

	BidPrice1   float64 `json:"bid_price_1"`
	BidVolume1  int64   `json:"bid_volume_1"`
	DiffBidVol1 int64   `json:"diff_bid_vol_1"`
	BidPrice2   float64 `json:"bid_price_2"`
	BidVolume2  int64   `json:"bid_volume_2"`
	DiffBidVol2 int64   `json:"diff_bid_vol_2"`
	BidPrice3   float64 `json:"bid_price_3"`
	BidVolume3  int64   `json:"bid_volume_3"`
	DiffBidVol3 int64   `json:"diff_bid_vol_3"`
	BidPrice4   float64 `json:"bid_price_4"`
	BidVolume4  int64   `json:"bid_volume_4"`
	DiffBidVol4 int64   `json:"diff_bid_vol_4"`
	BidPrice5   float64 `json:"bid_price_5"`
	BidVolume5  int64   `json:"bid_volume_5"`
	DiffBidVol5 int64   `json:"diff_bid_vol_5"`
	AskPrice1   float64 `json:"ask_price_1"`
	AskVolume1  int64   `json:"ask_volume_1"`
	DiffAskVol1 int64   `json:"diff_ask_vol_1"`
	AskPrice2   float64 `json:"ask_price_2"`
	AskVolume2  int64   `json:"ask_volume_2"`
	DiffAskVol2 int64   `json:"diff_ask_vol_2"`
	AskPrice3   float64 `json:"ask_price_3"`
	AskVolume3  int64   `json:"ask_volume_3"`
	DiffAskVol3 int64   `json:"diff_ask_vol_3"`
	AskPrice4   float64 `json:"ask_price_4"`
	AskVolume4  int64   `json:"ask_volume_4"`
	DiffAskVol4 int64   `json:"diff_ask_vol_4"`
	AskPrice5   float64 `json:"ask_price_5"`
	AskVolume5  int64   `json:"ask_volume_5"`
	DiffAskVol5 int64   `json:"diff_ask_vol_5"`
	Suspend     int64   `json:"suspend"`
	Simtrade    int64   `json:"simtrade"`
}

// RealTimeTick RealTimeTick
type RealTimeTick struct {
	Stock    *Stock    `json:"stock"`
	StockNum string    `json:"stock_num"`
	TickTime time.Time `json:"tick_time"`

	Open            float64 `json:"open"`
	AvgPrice        float64 `json:"avg_price"`
	Close           float64 `json:"close"`
	High            float64 `json:"high"`
	Low             float64 `json:"low"`
	Amount          float64 `json:"amount"`
	AmountSum       float64 `json:"amount_sum"`
	Volume          int64   `json:"volume"`
	VolumeSum       int64   `json:"volume_sum"`
	TickType        int64   `json:"tick_type"`
	ChgType         int64   `json:"chg_type"`
	PriceChg        float64 `json:"price_chg"`
	PctChg          float64 `json:"pct_chg"`
	BidSideTotalVol int64   `json:"bid_side_total_vol"`
	AskSideTotalVol int64   `json:"ask_side_total_vol"`
	BidSideTotalCnt int64   `json:"bid_side_total_cnt"`
	AskSideTotalCnt int64   `json:"ask_side_total_cnt"`
	Suspend         int64   `json:"suspend"`
	Simtrade        int64   `json:"simtrade"`
}
