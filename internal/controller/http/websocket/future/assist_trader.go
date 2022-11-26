package future

import (
	"context"
	"sync"
	"time"

	"tmt/internal/entity"
	"tmt/internal/usecase"
)

type assistTrader struct {
	o              usecase.Order
	ctx            context.Context
	assistOption   HalfAutomationOption
	tickChan       chan *entity.RealTimeFutureTick
	assistOrderMap map[string]*entity.FutureOrder
	orderLock      sync.Mutex
}

func newAssistTrader(ctx context.Context, o usecase.Order) *assistTrader {
	a := &assistTrader{
		ctx:            ctx,
		assistOrderMap: make(map[string]*entity.FutureOrder),
		tickChan:       make(chan *entity.RealTimeFutureTick),
		o:              o,
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
			}
			a.orderLock.Unlock()
		}
	}
}

func (a *assistTrader) processTick() {
	for {
		_, ok := <-a.tickChan
		if !ok {
			return
		}
	}
}

// func (a *assistTrader) placeOrder(order *entity.FutureOrder) error {
// 	var err error
// 	switch order.Action {
// 	case entity.ActionBuy:
// 		order.OrderID, order.Status, err = a.o.BuyFuture(order)
// 		if err != nil {
// 			return err
// 		}

// 	case entity.ActionSell:
// 		order.OrderID, order.Status, err = a.o.SellFuture(order)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

func (a *assistTrader) addAssistOrder(order *entity.FutureOrder, option HalfAutomationOption) {
	a.orderLock.Lock()
	defer a.orderLock.Unlock()

	a.assistOrderMap[order.OrderID] = order
	a.assistOption = option
	a.tickChan = make(chan *entity.RealTimeFutureTick)
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

func (a *assistTrader) isAssistDone() bool {
	defer a.orderLock.Unlock()
	a.orderLock.Lock()
	return len(a.assistOrderMap) == 0
}
