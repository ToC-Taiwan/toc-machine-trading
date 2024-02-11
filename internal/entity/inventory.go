package entity

import (
	"time"
)

type PositionStock struct {
	ID        int       `json:"-"`
	StockNum  string    `json:"StockNum"`
	Date      time.Time `json:"Date"`
	Quantity  int       `json:"Quantity"`
	Price     float64   `json:"Price"`
	LastPrice float64   `json:"LastPrice"`
	Dseq      string    `json:"Dseq"`
	Direction string    `json:"Direction"`
	Pnl       float64   `json:"Pnl"`
	Fee       float64   `json:"Fee"`
	InvID     string    `json:"InvID"`
}

type InventoryBase struct {
	UUID     string
	AvgPrice float64
	Date     time.Time
}

type InventoryStock struct {
	InventoryBase
	StockNum string
	Lot      int
	Share    int
	Position []*PositionStock
}

type InventoryFuture struct {
	InventoryBase
	Code     string
	Position int
}
