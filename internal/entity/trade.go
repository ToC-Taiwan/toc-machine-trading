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
)

const (
	ActionStringNone string = "None"
	ActionStringBuy  string = "Buy"
	ActionStringSell string = "Sell"
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
	default:
		return ""
	}
}

func StringToOrderAction(s string) OrderAction {
	switch s {
	case ActionStringBuy:
		return ActionBuy
	case ActionStringSell:
		return ActionSell
	default:
		return ActionNone
	}
}

func StringToOrderStatus(s string) OrderStatus {
	switch s {
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
	default:
		return StatusUnknow
	}
}

// OrderDetail -.
type OrderDetail struct {
	OrderID   string      `json:"order_id"`
	Price     float64     `json:"price"`
	Status    OrderStatus `json:"status"`
	Action    OrderAction `json:"action"`
	OrderTime time.Time   `json:"order_time"`
}

func (o *OrderDetail) Cancellable() bool {
	switch o.Status {
	case StatusPendingSubmit, StatusPreSubmitted, StatusSubmitted, StatusPartFilled:
		return true
	default:
		return false
	}
}

// StockOrder -.
type StockOrder struct {
	StockNum string `json:"stock_num"`
	Lot      int64  `json:"lot"`
	Share    int64  `json:"share"`
	Stock    *Stock `json:"stock"`

	OrderDetail `json:"base_order"`
}

func (s *StockOrder) StockOrderStatusString() string {
	return fmt.Sprintf("%s %s %s %.0f x (%d+%d)", s.OrderDetail.Status.String(), s.OrderDetail.Action.String(), s.StockNum, s.OrderDetail.Price, s.Lot*1000, s.Share)
}

// func (s *StockOrder) FixTime() *StockOrder {
// 	if time.Since(s.OrderTime) > 12*time.Hour {
// 		s.OrderTime = time.Now()
// 	}
// 	return s
// }

// FutureOrder -.
type FutureOrder struct {
	Code     string  `json:"code"`
	Position int64   `json:"position"`
	Future   *Future `json:"future"`

	OrderDetail `json:"base_order"`
}

func (f *FutureOrder) FutureOrderStatusString() string {
	return fmt.Sprintf("%s %s %s %.0f x %d", f.OrderDetail.Status.String(), f.OrderDetail.Action.String(), f.Code, f.OrderDetail.Price, f.Position)
}

func (f *FutureOrder) String() string {
	return fmt.Sprintf("%s %s %.0f x %d", f.OrderDetail.Action.String(), f.Code, f.OrderDetail.Price, f.Position)
}

// func (f *FutureOrder) FixTime() *FutureOrder {
// 	if time.Since(f.OrderTime) > 12*time.Hour {
// 		f.OrderTime = time.Now()
// 	}
// 	return f
// }

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
	Position  int64   `json:"position"`
	Price     float64 `json:"price"`
	LastPrice float64 `json:"last_price"`
	Pnl       float64 `json:"pnl"`
}

type FuturePositionArr []*FuturePosition
