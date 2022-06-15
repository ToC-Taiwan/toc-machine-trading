package entity

import (
	"time"
)

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
