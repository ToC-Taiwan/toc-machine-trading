package entity

import (
	"fmt"
	"time"
)

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
	// StatusPartFilled -.
	StatusPartFilled
	// StatusAborted -.
	StatusAborted
)

const (
	ActionStringNone      string = "None"
	ActionStringBuy       string = "Buy"
	ActionStringSell      string = "Sell"
	ActionStringSellFirst string = "SellFirst"
	ActionStringBuyLater  string = "BuyLater"
)

const (
	StatusStringUnknow        string = "Unknow"
	StatusStringPendingSubmit string = "PendingSubmit"
	StatusStringPreSubmitted  string = "PreSubmitted"
	StatusStringSubmitted     string = "Submitted"
	StatusStringFailed        string = "Failed"
	StatusStringCancelled     string = "Cancelled"
	StatusStringFilled        string = "Filled"
	StatusStringPartFilled    string = "PartFilled"
	StatusStringAborted       string = "Aborted"
)

// OrderAction -.
type OrderAction int64

func (a OrderAction) String() string {
	switch a {
	case ActionNone:
		return ActionStringNone
	case ActionBuy:
		return ActionStringBuy
	case ActionSell:
		return ActionStringSell
	case ActionSellFirst:
		return ActionStringSellFirst
	case ActionBuyLater:
		return ActionStringBuyLater
	default:
		return ""
	}
}

// OrderStatus -.
type OrderStatus int64

func (s OrderStatus) String() string {
	switch s {
	case StatusUnknow:
		return "Unknow"
	case StatusPendingSubmit:
		return "PendingSubmit"
	case StatusPreSubmitted:
		return "PreSubmitted"
	case StatusSubmitted:
		return "Submitted"
	case StatusFailed:
		return "Failed"
	case StatusCancelled:
		return "Cancelled"
	case StatusFilled:
		return "Filled"
	case StatusPartFilled:
		return "PartFilled"
	case StatusAborted:
		return "Aborted"
	default:
		return ""
	}
}

func StringToOrderAction(s string) OrderAction {
	switch s {
	case ActionStringNone:
		return ActionNone
	case ActionStringBuy:
		return ActionBuy
	case ActionStringSell:
		return ActionSell
	case ActionStringSellFirst:
		return ActionSellFirst
	case ActionStringBuyLater:
		return ActionBuyLater
	default:
		return ActionNone
	}
}

func StringToOrderStatus(s string) OrderStatus {
	switch s {
	case StatusStringUnknow:
		return StatusUnknow
	case StatusStringPendingSubmit:
		return StatusPendingSubmit
	case StatusStringPreSubmitted:
		return StatusPreSubmitted
	case StatusStringSubmitted:
		return StatusSubmitted
	case StatusStringFailed:
		return StatusFailed
	case StatusStringCancelled:
		return StatusCancelled
	case StatusStringFilled:
		return StatusFilled
	case StatusStringPartFilled:
		return StatusPartFilled
	case StatusStringAborted:
		return StatusAborted
	default:
		return StatusUnknow
	}
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

func (o *BaseOrder) Cancellable() bool {
	switch o.Status {
	case StatusPendingSubmit, StatusPreSubmitted, StatusSubmitted, StatusPartFilled:
		return true
	default:
		return false
	}
}

func (o *BaseOrder) FilledQty() int64 {
	if o.Status != StatusFilled && o.Status != StatusPartFilled {
		return 0
	}

	switch o.Action {
	case ActionBuy:
		return o.Quantity
	case ActionSell:
		return -o.Quantity
	case ActionSellFirst:
		return -o.Quantity
	case ActionBuyLater:
		return o.Quantity
	default:
		return 0
	}
}

type StockOrderArr []*StockOrder

func (s StockOrderArr) SplitManualAndGroupID() (map[string]StockOrderArr, StockOrderArr) {
	group := make(map[string]StockOrderArr)
	var manual StockOrderArr
	for _, v := range s {
		if v.Manual {
			manual = append(manual, v)
		} else {
			group[v.OrderID] = append(group[v.OrderID], v)
		}
	}
	return group, manual
}

func (s StockOrderArr) IsAllDone() bool {
	var qty int64
	for _, v := range s {
		if v.Status != StatusFilled {
			continue
		}

		switch v.Action {
		case ActionBuy:
			qty += v.Quantity
		case ActionSell:
			qty -= v.Quantity
		case ActionSellFirst:
			qty -= v.Quantity
		case ActionBuyLater:
			qty += v.Quantity
		}
	}
	return qty == 0
}

// StockOrder -.
type StockOrder struct {
	BaseOrder `json:"base_order"`

	StockNum string `json:"stock_num"`
	Stock    *Stock `json:"stock"`
	Manual   bool   `json:"manual"`
}

func (s *StockOrder) StockOrderStatusString() string {
	return fmt.Sprintf("%s %s %s %.0f x %d", s.BaseOrder.Status.String(), s.BaseOrder.Action.String(), s.StockNum, s.BaseOrder.Price, s.BaseOrder.Quantity)
}

func (s *StockOrder) ToManual() *StockOrder {
	s.Manual = true
	s.GroupID = "-"

	if time.Since(s.OrderTime) > 12*time.Hour {
		s.OrderTime = time.Now()
	}

	s.TradeTime = s.OrderTime
	s.TickTime = s.OrderTime
	return s
}

type FutureOrderArr []*FutureOrder

func (s FutureOrderArr) SplitManualAndGroupID() (map[string]FutureOrderArr, FutureOrderArr) {
	group := make(map[string]FutureOrderArr)
	var manual FutureOrderArr
	for _, v := range s {
		if v.Manual {
			manual = append(manual, v)
		} else {
			group[v.OrderID] = append(group[v.OrderID], v)
		}
	}
	return group, manual
}

func (s FutureOrderArr) IsAllDone() bool {
	var qty int64
	for _, v := range s {
		if v.Status != StatusFilled {
			continue
		}

		switch v.Action {
		case ActionBuy:
			qty += v.Quantity
		case ActionSell:
			qty -= v.Quantity
		case ActionSellFirst:
			qty -= v.Quantity
		case ActionBuyLater:
			qty += v.Quantity
		}
	}
	return qty == 0
}

// FutureOrder -.
type FutureOrder struct {
	BaseOrder `json:"base_order"`

	Code   string  `json:"code"`
	Future *Future `json:"future"`
	Manual bool    `json:"manual"`
}

func (f *FutureOrder) FutureOrderStatusString() string {
	return fmt.Sprintf("%s %s %s %.0f x %d", f.BaseOrder.Status.String(), f.BaseOrder.Action.String(), f.Code, f.BaseOrder.Price, f.BaseOrder.Quantity)
}

func (f *FutureOrder) ToManual() *FutureOrder {
	f.Manual = true
	f.GroupID = "-"

	if time.Since(f.OrderTime) > 12*time.Hour {
		f.OrderTime = time.Now()
	}

	f.TradeTime = f.OrderTime
	f.TickTime = f.OrderTime
	return f
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

type FuturePositionArr []*FuturePosition
