package entity

import (
	"time"
)

type InventoryBankDetail struct {
	ID       int       `json:"id" yaml:"id"`
	BankID   int       `json:"bank_id" yaml:"bank_id"`
	AvgPrice float64   `json:"avg_price" yaml:"avg_price"`
	Updated  time.Time `json:"updated" yaml:"updated"`
}

type InventoryStock struct {
	InventoryBankDetail
	StockNum string `json:"stock_num" yaml:"stock_num"`
	Lot      int    `json:"lot" yaml:"lot"`
	Share    int    `json:"share" yaml:"share"`
}

type InventoryFuture struct {
	InventoryBankDetail
	Code     string `json:"code" yaml:"code"`
	Position int    `json:"position" yaml:"position"`
}
