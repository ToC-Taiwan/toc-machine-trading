package future

import (
	"sync"

	"tmt/internal/entity"
)

type AutomationType int

func (a AutomationType) String() string {
	switch a {
	case AutomationByBalance:
		return "By balance"
	case AutomationByTimePeriod:
		return "By time period"
	case AutomationByTimePeriodAndBalance:
		return "By time period and balance"
	}
	return "unknown"
}

const (
	AutomationNone AutomationType = iota
	AutomationByBalance
	AutomationByTimePeriod
	AutomationByTimePeriodAndBalance
)

type clientOrder struct {
	Code   string               `json:"code"`
	Action entity.OrderAction   `json:"action"`
	Price  float64              `json:"price"`
	Qty    int64                `json:"qty"`
	Option halfAutomationOption `json:"option"`
}

type halfAutomationOption struct {
	AutomationType AutomationType `json:"automation_type"`
	ByBalanceHigh  float64        `json:"by_balance_high"`
	ByBalanceLow   float64        `json:"by_balance_low"`
	ByTimePeriod   int64          `json:"by_time_period"`
}

func (f *clientOrder) toFutureOrderArr() []*entity.FutureOrder {
	var orders []*entity.FutureOrder
	for i := 0; i < int(f.Qty); i++ {
		orders = append(orders, &entity.FutureOrder{
			Code: f.Code,
			BaseOrder: entity.BaseOrder{
				Action:   f.Action,
				Quantity: 1,
				Price:    f.Price,
			},
		})
	}
	return orders
}

type waitingList struct {
	list map[string]*entity.FutureOrder
	m    sync.RWMutex
}

func newWaitingList() *waitingList {
	return &waitingList{
		list: make(map[string]*entity.FutureOrder),
	}
}

func (w *waitingList) empty() bool {
	defer w.m.RUnlock()
	w.m.RLock()
	return len(w.list) == 0
}

func (w *waitingList) orderIDExist(orderID string) bool {
	defer w.m.RUnlock()
	w.m.RLock()
	_, ok := w.list[orderID]
	return ok
}

func (w *waitingList) add(order *entity.FutureOrder) {
	defer w.m.Unlock()
	w.m.Lock()
	w.list[order.OrderID] = order
}

func (w *waitingList) remove(orderID string) {
	defer w.m.Unlock()
	w.m.Lock()
	delete(w.list, orderID)
}
