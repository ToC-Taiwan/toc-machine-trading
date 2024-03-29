package entity

import (
	"time"
)

// SinopacEvent SinopacEvent
type SinopacEvent struct {
	ID        int64     `json:"id"`
	EventCode int64     `json:"event_code"`
	Response  int64     `json:"response"`
	Event     string    `json:"event"`
	Info      string    `json:"info"`
	EventTime time.Time `json:"event_time"`
}

// RealTimeStockTick -.
type RealTimeStockTick struct {
	StockNum string `json:"stock_num"`
	Stock    *Stock `json:"stock"`

	TickTime        time.Time `json:"tick_time"`
	Open            float64   `json:"open"`
	AvgPrice        float64   `json:"avg_price"`
	Close           float64   `json:"close"`
	High            float64   `json:"high"`
	Low             float64   `json:"low"`
	Amount          float64   `json:"amount"`
	AmountSum       float64   `json:"amount_sum"`
	Volume          int64     `json:"volume"`
	VolumeSum       int64     `json:"volume_sum"`
	TickType        int64     `json:"tick_type"`
	ChgType         int64     `json:"chg_type"`
	PriceChg        float64   `json:"price_chg"`
	PctChg          float64   `json:"pct_chg"`
	BidSideTotalVol int64     `json:"bid_side_total_vol"`
	AskSideTotalVol int64     `json:"ask_side_total_vol"`
	BidSideTotalCnt int64     `json:"bid_side_total_cnt"`
	AskSideTotalCnt int64     `json:"ask_side_total_cnt"`
}

// RealTimeFutureTick -.
type RealTimeFutureTick struct {
	Code string `json:"code"`

	TickTime        time.Time `json:"tick_time"`
	Open            float64   `json:"open"`
	UnderlyingPrice float64   `json:"underlying_price"`
	BidSideTotalVol int64     `json:"bid_side_total_vol"`
	AskSideTotalVol int64     `json:"ask_side_total_vol"`
	AvgPrice        float64   `json:"avg_price"`
	Close           float64   `json:"close"`
	High            float64   `json:"high"`
	Low             float64   `json:"low"`
	Amount          float64   `json:"amount"`
	TotalAmount     float64   `json:"total_amount"`
	Volume          int64     `json:"volume"`
	TotalVolume     int64     `json:"total_volume"`
	TickType        int64     `json:"tick_type"`
	ChgType         int64     `json:"chg_type"`
	PriceChg        float64   `json:"price_chg"`
	PctChg          float64   `json:"pct_chg"`
}

// StockSnapShot -.
type StockSnapShot struct {
	StockNum  string `json:"stock_num"`
	StockName string `json:"stock_name"`

	SnapShotBase
}

// FutureSnapShot -.
type FutureSnapShot struct {
	Code       string `json:"code"`
	FutureName string `json:"future_name"`

	SnapShotBase
}

func (f *FutureSnapShot) ToRealTimeFutureTick() *RealTimeFutureTick {
	var tickType, chgType int64
	switch f.TickType {
	case "Sell":
		tickType = 1
	case "Buy":
		tickType = 2
	}

	switch f.ChgType {
	case "LimitUp":
		chgType = 1
	case "Up":
		chgType = 2
	case "Unchanged":
		chgType = 3
	case "Dowm":
		chgType = 4
	case "LimitDown":
		chgType = 5
	}

	return &RealTimeFutureTick{
		Code:        f.Code,
		TickTime:    f.SnapTime,
		Open:        f.Open,
		Close:       f.Close,
		High:        f.High,
		Low:         f.Low,
		Amount:      float64(f.Amount),
		TotalAmount: float64(f.AmountSum),
		Volume:      f.Volume,
		TotalVolume: f.VolumeSum,
		TickType:    tickType,
		ChgType:     chgType,
		PriceChg:    f.PriceChg,
		PctChg:      f.PctChg,
	}
}

// SnapShotBase -.
type SnapShotBase struct {
	SnapTime        time.Time `json:"snap_time"`
	Open            float64   `json:"open"`
	High            float64   `json:"high"`
	Low             float64   `json:"low"`
	Close           float64   `json:"close"`
	TickType        string    `json:"tick_type"`
	PriceChg        float64   `json:"price_chg"`
	PctChg          float64   `json:"pct_chg"`
	ChgType         string    `json:"chg_type"`
	Volume          int64     `json:"volume"`
	VolumeSum       int64     `json:"volume_sum"`
	Amount          int64     `json:"amount"`
	AmountSum       int64     `json:"amount_sum"`
	YesterdayVolume float64   `json:"yesterday_volume"`
	VolumeRatio     float64   `json:"volume_ratio"`
}

// YahooPrice -.
type YahooPrice struct {
	Last      float64   `json:"last"`
	Price     float64   `json:"price"`
	UpdatedAt time.Time `json:"updated_at"`
}

type IndexStatus struct {
	BreakCount int64   `json:"break_count"`
	PriceChg   float64 `json:"price_chg"`
}

func NewIndexStatus() *IndexStatus {
	return &IndexStatus{
		BreakCount: 0,
		PriceChg:   0,
	}
}

type TradeIndex struct {
	TSE    *IndexStatus `json:"tse"`
	OTC    *IndexStatus `json:"otc"`
	Nasdaq *IndexStatus `json:"nasdaq"`
	NF     *IndexStatus `json:"nf"`
}

func (i *IndexStatus) UpdateIndexStatus(priceChange float64) {
	if i.PriceChg == 0 {
		i.PriceChg = priceChange
		return
	}

	switch {
	case priceChange > i.PriceChg:
		if i.BreakCount < 0 {
			i.BreakCount = 0
		}
		i.BreakCount++
	case priceChange < i.PriceChg:
		if i.BreakCount > 0 {
			i.BreakCount = 0
		}
		i.BreakCount--
	}
	i.PriceChg = priceChange
}
