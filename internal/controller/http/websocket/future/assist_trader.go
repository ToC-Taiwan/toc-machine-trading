package future

import (
	"context"
	"errors"
	"sync"
	"time"

	"tmt/internal/entity"
	"tmt/internal/usecase"
)

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

	go a.checkAssistStatus()
	go a.processTick()

	return a
}

func (a *assistTrader) getTickChan() chan *entity.RealTimeFutureTick {
	return a.tickChan
}

func (a *assistTrader) checkAssistStatus() {
	for {
		select {
		case <-a.ctx.Done():
			return

		case <-time.After(5 * time.Second):
			a.orderLock.Lock()
			var qty int64
			for _, order := range a.assistOrderMap {
				if order.Status == entity.StatusFilled {
					switch order.Action {
					case entity.ActionBuy:
						qty += order.Quantity
					case entity.ActionSell:
						qty -= order.Quantity
					}
				}
			}

			if qty == 0 {
				a.assistOrderMap = make(map[string]*entity.FutureOrder)
				a.assisting = false
			}
			a.orderLock.Unlock()
		}
	}
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
			a.cheeckByBalance(tick)
		case AutomationByTimePeriod:
			a.cheeckByTime(tick)
		case AutomationByTimePeriodAndBalance:
			a.cheeckByTimeAndBalance(tick)
		}
	}
}

func (a *assistTrader) cheeckByTime(tick *entity.RealTimeFutureTick) {
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

func (a *assistTrader) cheeckByBalance(tick *entity.RealTimeFutureTick) {
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

func (a *assistTrader) cheeckByTimeAndBalance(tick *entity.RealTimeFutureTick) {
	if tick.Close > a.assistOption.ByBalanceLow+a.basePrice && tick.Close < a.assistOption.ByBalanceHigh+a.basePrice {
		return
	}

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

		case <-time.After(3 * time.Second):
			a.orderLock.RLock()
			order, ok := a.assistOrderMap[orderID]
			if !ok {
				a.orderLock.RUnlock()
				continue
			}

			if !order.Cancellable() {
				delete(cancelOrderMap, orderID)
			} else if time.Since(order.TradeTime) > 10*time.Second && cancelOrderMap[orderID] == nil {
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

func (a *assistTrader) addAssistOrder(order *entity.FutureOrder, option HalfAutomationOption) {
	a.orderLock.Lock()
	defer a.orderLock.Unlock()

	a.assistOrderMap[order.OrderID] = order
	a.assistOption = option

	a.qty = order.Quantity
	a.code = order.Code

	switch order.Action {
	case entity.ActionBuy:
		a.action = entity.ActionSell
	case entity.ActionSell:
		a.action = entity.ActionBuy
	}

	a.basePrice = order.Price
	a.tradeTime = order.TradeTime

	a.assisting = true
	a.waitingOrder = nil
}

func (a *assistTrader) updateOrderStatus(o *entity.FutureOrder) *entity.FutureOrder {
	a.orderLock.Lock()
	defer a.orderLock.Unlock()

	if cache, ok := a.assistOrderMap[o.OrderID]; ok && cache.Status != o.Status {
		o.TradeTime = cache.TradeTime
		a.assistOrderMap[o.OrderID] = o
		return o
	}
	return nil
}

func (a *assistTrader) isAssisting() bool {
	return a.assisting
}
