package entity

import (
	"time"
)

type Inventory struct {
	ID       int       `json:"id,omitempty" yaml:"id"`
	BankID   int       `json:"bank_id,omitempty" yaml:"bank_id"`
	AvgPrice float64   `json:"avg_price,omitempty" yaml:"avg_price"`
	Quantity int       `json:"quantity,omitempty" yaml:"quantity"`
	Updated  time.Time `json:"updated,omitempty" yaml:"updated"`
}

type InventoryStock struct {
	StockNum string `json:"stock_num,omitempty" yaml:"stock_num"`
	Inventory
}

type InventoryFuture struct {
	Code string `json:"code,omitempty" yaml:"code"`
	Inventory
}
