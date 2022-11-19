package entity

import (
	"time"
)

// OrderAction -.
type OrderAction int64

// OrderStatus -.
type OrderStatus int64

const (
	// ActionNone -.
	ActionNone OrderAction = iota
	// ActionBuy -.
	ActionBuy
	// ActionSell -.
	ActionSell
	// ActionSellFirst -.
	ActionSellFirst
	// ActionBuyLater -.
	ActionBuyLater
)

const (
	// StatusUnknow -.
	StatusUnknow OrderStatus = iota
	// StatusPendingSubmit -.
	StatusPendingSubmit
	// StatusPreSubmitted -.
	StatusPreSubmitted
	// StatusSubmitted -.
	StatusSubmitted
	// StatusFailed -.
	StatusFailed
	// StatusCancelled -.
	StatusCancelled
	// StatusFilled -.
	StatusFilled
	// StatusFilling -.
	StatusFilling
	// StatusAborted -.
	StatusAborted
)

// ActionListMap ActionListMap
var ActionListMap = map[string]OrderAction{
	"Buy":  ActionBuy,
	"Sell": ActionSell,
}

// StatusListMap StatusListMap
var StatusListMap = map[string]OrderStatus{
	"PendingSubmit": StatusPendingSubmit, // 傳送中
	"PreSubmitted":  StatusPreSubmitted,  // 預約單
	"Submitted":     StatusSubmitted,     // 傳送成功
	"Failed":        StatusFailed,        // 失敗
	"Cancelled":     StatusCancelled,     // 已刪除
	"Filled":        StatusFilled,        // 完全成交
	"Filling":       StatusFilling,       // 部分成交
}

// BaseOrder -.
type BaseOrder struct {
	OrderID   string      `json:"order_id"`
	Status    OrderStatus `json:"status"`
	OrderTime time.Time   `json:"order_time"`
	Action    OrderAction `json:"action"`
	Price     float64     `json:"price"`
	Quantity  int64       `json:"quantity"`
	TradeTime time.Time   `json:"trade_time"`
	TickTime  time.Time   `json:"tick_time"`
	GroupID   string      `json:"group_id"`
}

func (o *BaseOrder) Cancellabel() bool {
	switch o.Status {
	case StatusPendingSubmit, StatusPreSubmitted, StatusSubmitted, StatusFilling:
		return true
	default:
		return false
	}
}

// StockOrder -.
type StockOrder struct {
	BaseOrder `json:"base_order"`

	StockNum string `json:"stock_num"`
	Stock    *Stock `json:"stock"`
	Manual   bool   `json:"manual"`
}

// FutureOrder -.
type FutureOrder struct {
	BaseOrder `json:"base_order"`

	Code   string  `json:"code"`
	Future *Future `json:"future"`
	Manual bool    `json:"manual"`
}

// StockTradeBalance -.
type StockTradeBalance struct {
	ID              int64     `json:"id"`
	TradeCount      int64     `json:"trade_count"`
	Forward         int64     `json:"forward"`
	Reverse         int64     `json:"reverse"`
	OriginalBalance int64     `json:"original_balance"`
	Discount        int64     `json:"discount"`
	Total           int64     `json:"total"`
	TradeDay        time.Time `json:"trade_day"`
}

// FutureTradeBalance -.
type FutureTradeBalance struct {
	ID         int64     `json:"id"`
	TradeCount int64     `json:"trade_count"`
	Forward    int64     `json:"forward"`
	Reverse    int64     `json:"reverse"`
	Total      int64     `json:"total"`
	TradeDay   time.Time `json:"trade_day"`
}

// FuturePosition -.
type FuturePosition struct {
	Code      string  `json:"code"`
	Direction string  `json:"direction"`
	Quantity  int64   `json:"quantity"`
	Price     float64 `json:"price"`
	LastPrice float64 `json:"last_price"`
	Pnl       float64 `json:"pnl"`
}
