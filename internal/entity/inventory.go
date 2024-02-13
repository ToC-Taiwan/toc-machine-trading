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
	UUID     string    `json:"UUID"`
	AvgPrice float64   `json:"AvgPrice"`
	Date     time.Time `json:"Date"`
}

type InventoryStock struct {
	InventoryBase
	StockNum string           `json:"StockNum"`
	Lot      int              `json:"Lot"`
	Share    int              `json:"Share"`
	Position []*PositionStock `json:"Position"`
}

type InventoryFuture struct {
	InventoryBase
	Code     string
	Position int
}
