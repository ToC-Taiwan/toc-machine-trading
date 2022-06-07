package entity

import "time"

// OrderStatus OrderStatus
type OrderStatus struct {
	Stock     *Stock    `json:"stock"`
	StockNum  string    `json:"stock_num"`
	OrderTime time.Time `json:"order_time"`

	Action   int64   `json:"action"`
	Price    float64 `json:"price"`
	Quantity int64   `json:"quantity"`
	Status   int64   `json:"status"`
	OrderID  string  `json:"order_id"`
}

// Order Order
type Order struct {
	StockNum string  `json:"stock_num"`
	Action   int64   `json:"action"`
	Price    float64 `json:"price"`
	Quantity int64   `json:"quantity"`
}
