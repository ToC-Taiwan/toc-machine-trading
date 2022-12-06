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
	Code   string  `json:"code"`
	Future *Future `json:"future"`

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

type RealTimeFutureTickArr []*RealTimeFutureTick

func (r RealTimeFutureTickArr) GetOutInVolume() OutInVolume {
	var outInVolume OutInVolume
	for _, tick := range r {
		switch tick.TickType {
		case 1:
			outInVolume.OutVolume += tick.Volume
		case 2:
			outInVolume.InVolume += tick.Volume
		}
	}
	return outInVolume
}

type OutInVolume struct {
	OutVolume int64 `json:"out_volume"`
	InVolume  int64 `json:"in_volume"`
}

// RealTimeStockBidAsk -.
type RealTimeStockBidAsk struct {
	StockNum string `json:"stock_num"`
	Stock    *Stock `json:"stock"`
	BidAskBase
}

// FutureRealTimeBidAsk -.
type FutureRealTimeBidAsk struct {
	Code                 string  `json:"code"`
	BidTotalVol          int64   `json:"bid_total_vol"`
	AskTotalVol          int64   `json:"ask_total_vol"`
	UnderlyingPrice      float64 `json:"underlying_price"`
	FirstDerivedBidPrice float64 `json:"first_derived_bid_price"`
	FirstDerivedAskPrice float64 `json:"first_derived_ask_price"`
	FirstDerivedBidVol   int64   `json:"first_derived_bid_vol"`
	FirstDerivedAskVol   int64   `json:"first_derived_ask_vol"`
	BidAskBase
}

// BidAskBase -.
type BidAskBase struct {
	BidAskTime  time.Time `json:"bid_ask_time"`
	BidPrice1   float64   `json:"bid_price_1"`
	BidVolume1  int64     `json:"bid_volume_1"`
	DiffBidVol1 int64     `json:"diff_bid_vol_1"`
	BidPrice2   float64   `json:"bid_price_2"`
	BidVolume2  int64     `json:"bid_volume_2"`
	DiffBidVol2 int64     `json:"diff_bid_vol_2"`
	BidPrice3   float64   `json:"bid_price_3"`
	BidVolume3  int64     `json:"bid_volume_3"`
	DiffBidVol3 int64     `json:"diff_bid_vol_3"`
	BidPrice4   float64   `json:"bid_price_4"`
	BidVolume4  int64     `json:"bid_volume_4"`
	DiffBidVol4 int64     `json:"diff_bid_vol_4"`
	BidPrice5   float64   `json:"bid_price_5"`
	BidVolume5  int64     `json:"bid_volume_5"`
	DiffBidVol5 int64     `json:"diff_bid_vol_5"`
	AskPrice1   float64   `json:"ask_price_1"`
	AskVolume1  int64     `json:"ask_volume_1"`
	DiffAskVol1 int64     `json:"diff_ask_vol_1"`
	AskPrice2   float64   `json:"ask_price_2"`
	AskVolume2  int64     `json:"ask_volume_2"`
	DiffAskVol2 int64     `json:"diff_ask_vol_2"`
	AskPrice3   float64   `json:"ask_price_3"`
	AskVolume3  int64     `json:"ask_volume_3"`
	DiffAskVol3 int64     `json:"diff_ask_vol_3"`
	AskPrice4   float64   `json:"ask_price_4"`
	AskVolume4  int64     `json:"ask_volume_4"`
	DiffAskVol4 int64     `json:"diff_ask_vol_4"`
	AskPrice5   float64   `json:"ask_price_5"`
	AskVolume5  int64     `json:"ask_volume_5"`
	DiffAskVol5 int64     `json:"diff_ask_vol_5"`
}

// StockSnapShot -.
type StockSnapShot struct {
	StockNum  string `json:"stock_num"`
	StockName string `json:"stock_name"`
	SnapShotBase
}

// FutureSnapShot -.
type FutureSnapShot struct {
	Code   string  `json:"code"`
	Future *Future `json:"future"`
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
		Future:      f.Future,
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

type TradeIndex struct {
	TSE    *StockSnapShot `json:"tse"`
	OTC    *StockSnapShot `json:"otc"`
	Nasdaq *YahooPrice    `json:"nasdaq"`
	NF     *YahooPrice    `json:"nf"`
}
