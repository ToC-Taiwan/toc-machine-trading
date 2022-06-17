package entity

import (
	"time"
)

// ActionListMap ActionListMap
var ActionListMap = map[string]int64{
	"Buy":  1,
	"Sell": 2,
}

// StatusListMap StatusListMap
var StatusListMap = map[string]int64{
	"PendingSubmit": 1, // 傳送中
	"PreSubmitted":  2, // 預約單
	"Submitted":     3, // 傳送成功
	"Failed":        4, // 失敗
	"Cancelled":     5, // 已刪除
	"Filled":        6, // 完全成交
	"Filling":       7, // 部分成交
}

// OrderStatus OrderStatus
type OrderStatus struct {
	OrderID   string    `json:"order_id"`
	StockNum  string    `json:"stock_num"`
	Action    int64     `json:"action"`
	Price     float64   `json:"price"`
	Quantity  int64     `json:"quantity"`
	Status    int64     `json:"status"`
	OrderTime time.Time `json:"order_time"`
	Stock     *Stock    `json:"stock"`
}

// TradeBalance -.
type TradeBalance struct {
	ID              int64     `json:"id"`
	TradeCount      int64     `json:"trade_count"`
	Forward         int64     `json:"forward"`
	Reverse         int64     `json:"reverse"`
	OriginalBalance int64     `json:"original_balance"`
	Discount        int64     `json:"discount"`
	Total           int64     `json:"total"`
	TradeDay        time.Time `json:"trade_day"`
}

// Order Order
type Order struct {
	StockNum string  `json:"stock_num"`
	Action   int64   `json:"action"`
	Price    float64 `json:"price"`
	Quantity int64   `json:"quantity"`
}
