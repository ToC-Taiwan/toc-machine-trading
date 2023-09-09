package entity

import (
	"time"
)

type Inventory struct {
	ID       int       `json:"id" yaml:"id"`
	BankID   int       `json:"bank_id" yaml:"bank_id"`
	AvgPrice float64   `json:"avg_price" yaml:"avg_price"`
	Quantity int       `json:"quantity" yaml:"quantity"`
	Updated  time.Time `json:"updated" yaml:"updated"`
}

type InventoryStock struct {
	StockNum string `json:"stock_num" yaml:"stock_num"`
	Inventory
}

type InventoryFuture struct {
	Code string `json:"code" yaml:"code"`
	Inventory
}
