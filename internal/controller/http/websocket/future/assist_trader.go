package future

import (
	"context"
	"errors"
	"sync"
	"time"

	"tmt/internal/entity"
	"tmt/internal/usecase"
)

// TODO: need to consider partfilled condition

type assistTrader struct {
	o        usecase.Order
	ctx      context.Context
	tickChan chan *entity.RealTimeFutureTick

	assistOrderMap map[string]*entity.FutureOrder
	orderLock      sync.RWMutex

	assistOption HalfAutomationOption
	action       entity.OrderAction
	tradeTime    time.Time
	code         string
	qty          int64
	basePrice    float64

	assisting    bool
	waitingOrder *entity.FutureOrder
}

func newAssistTrader(ctx context.Context, o usecase.Order) *assistTrader {
	a := &assistTrader{
		o:              o,
		ctx:            ctx,
		tickChan:       make(chan *entity.RealTimeFutureTick),
		assistOrderMap: make(map[string]*entity.FutureOrder),
	}
	go a.processTick()
	return a
}

func (a *assistTrader) processTick() {
	for {
		tick, ok := <-a.tickChan
		if !ok {
			return
		}

		if !a.assisting || a.waitingOrder != nil {
			continue
		}

		switch a.assistOption.AutomationType {
		case AutomationByBalance:
			a.checkByBalance(tick)
		case AutomationByTimePeriod:
			a.checkByTime(tick)
		case AutomationByTimePeriodAndBalance:
			a.checkByTimeAndBalance(tick)
		}
	}
}

func (a *assistTrader) isAssisting() bool {
	return a.assisting
}

func (a *assistTrader) getTickChan() chan *entity.RealTimeFutureTick {
	return a.tickChan
}

func (a *assistTrader) addAssistOrder(order *entity.FutureOrder, option HalfAutomationOption) {
	defer a.orderLock.Unlock()
	a.orderLock.Lock()

	a.assistOrderMap = make(map[string]*entity.FutureOrder)
	a.assistOrderMap[order.OrderID] = order
	a.assistOption = option
	a.code = order.Code
	a.basePrice = order.Price
	a.qty = order.Quantity
	a.tradeTime = order.TradeTime

	switch order.Action {
	case entity.ActionBuy:
		a.action = entity.ActionSell
	case entity.ActionSell:
		a.action = entity.ActionBuy
	}

	a.assisting = true
	a.waitingOrder = order
	go a.cancelOverTimeBaseOrder(order.OrderID)
	go a.checkAssistStatus(order.OrderID)
}

func (a *assistTrader) cancelOverTimeBaseOrder(orderID string) {
	cancelOrderMap := make(map[string]*entity.FutureOrder)
	for {
		select {
		case <-a.ctx.Done():
			return

		case <-time.After(time.Second):
			a.orderLock.Lock()
			order := a.assistOrderMap[orderID]

			if !order.Cancellable() {
				a.orderLock.Unlock()
				return
			}

			if time.Since(order.TradeTime) > 10*time.Second && cancelOrderMap[order.OrderID] == nil {
				if e := a.cancelOrderByID(order.OrderID); e != nil {
					a.orderLock.Unlock()
					continue
				}
				cancelOrderMap[order.OrderID] = order
			}

			a.orderLock.Unlock()
		}
	}
}

func (a *assistTrader) checkAssistStatus(orderID string) {
	for {
		select {
		case <-a.ctx.Done():
			return

		case <-time.After(3 * time.Second):
			a.orderLock.Lock()
			order := a.assistOrderMap[orderID]

			if order.Status == entity.StatusCancelled {
				a.assisting = false
				a.orderLock.Unlock()
				return
			}

			if order.Status == entity.StatusFilled {
				if a.waitingOrder != nil && a.waitingOrder.OrderID == order.OrderID {
					a.waitingOrder = nil
				}

				if done := a.assistIsDone(a.assistOrderMap); done {
					a.assisting = false
					a.orderLock.Unlock()
					return
				}
			}

			a.orderLock.Unlock()
		}
	}
}

func (a *assistTrader) assistIsDone(orderMap map[string]*entity.FutureOrder) bool {
	var qty int64
	for _, order := range orderMap {
		qty += order.FilledQty()
	}

	return qty == 0
}

func (a *assistTrader) placeOrder(order *entity.FutureOrder) error {
	var err error
	switch order.Action {
	case entity.ActionBuy:
		order.OrderID, order.Status, err = a.o.BuyFuture(order)
		if err != nil {
			return err
		}

	case entity.ActionSell:
		order.OrderID, order.Status, err = a.o.SellFuture(order)
		if err != nil {
			return err
		}
	}

	a.orderLock.Lock()
	order.TradeTime = time.Now()
	a.waitingOrder = order
	a.assistOrderMap[order.OrderID] = order
	a.orderLock.Unlock()

	go a.checkOrderStatusByID(order.OrderID)
	return nil
}

func (a *assistTrader) checkOrderStatusByID(orderID string) {
	cancelOrderMap := make(map[string]*entity.FutureOrder)
	for {
		select {
		case <-a.ctx.Done():
			return

		case <-time.After(time.Second):
			a.orderLock.RLock()
			order := a.assistOrderMap[orderID]
			if !order.Cancellable() {
				a.waitingOrder = nil
				a.orderLock.RUnlock()
				return
			}

			if time.Since(order.TradeTime) > 10*time.Second && cancelOrderMap[orderID] == nil {
				if e := a.cancelOrderByID(orderID); e != nil {
					a.orderLock.RUnlock()
					continue
				}
				cancelOrderMap[orderID] = order
			}
			a.orderLock.RUnlock()
		}
	}
}

func (a *assistTrader) cancelOrderByID(orderID string) error {
	_, s, err := a.o.CancelFutureOrderID(orderID)
	if err != nil {
		return err
	}
	if s != entity.StatusCancelled {
		return errors.New("cancel order failed")
	}
	return nil
}

func (a *assistTrader) checkByTime(tick *entity.RealTimeFutureTick) {
	if time.Since(a.tradeTime) < time.Duration(a.assistOption.ByTimePeriod)*time.Minute {
		return
	}

	order := &entity.FutureOrder{
		Code: a.code,
		BaseOrder: entity.BaseOrder{
			Action:   a.action,
			Price:    tick.Close,
			Quantity: a.qty,
		},
	}

	if e := a.placeOrder(order); e != nil {
		return
	}
}

func (a *assistTrader) checkByBalance(tick *entity.RealTimeFutureTick) {
	if tick.Close > a.assistOption.ByBalanceLow+a.basePrice && tick.Close < a.assistOption.ByBalanceHigh+a.basePrice {
		return
	}

	order := &entity.FutureOrder{
		Code: a.code,
		BaseOrder: entity.BaseOrder{
			Action:   a.action,
			Price:    tick.Close,
			Quantity: a.qty,
		},
	}

	if e := a.placeOrder(order); e != nil {
		return
	}
}

func (a *assistTrader) checkByTimeAndBalance(tick *entity.RealTimeFutureTick) {
	var byBalance, byTime bool
	if tick.Close > a.assistOption.ByBalanceLow+a.basePrice && tick.Close < a.assistOption.ByBalanceHigh+a.basePrice {
		byBalance = false
	}

	if time.Since(a.tradeTime) < time.Duration(a.assistOption.ByTimePeriod)*time.Minute {
		byTime = false
	}

	if !byBalance && !byTime {
		return
	}

	order := &entity.FutureOrder{
		Code: a.code,
		BaseOrder: entity.BaseOrder{
			Action:   a.action,
			Price:    tick.Close,
			Quantity: a.qty,
		},
	}

	if e := a.placeOrder(order); e != nil {
		return
	}
}

func (a *assistTrader) updateOrderStatus(o *entity.FutureOrder) *entity.FutureOrder {
	defer a.orderLock.Unlock()
	a.orderLock.Lock()

	if cache, ok := a.assistOrderMap[o.OrderID]; ok && cache.Status != o.Status {
		o.TradeTime = cache.TradeTime
		a.assistOrderMap[o.OrderID] = o
		return o
	}
	return nil
}
