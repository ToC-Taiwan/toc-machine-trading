package future

import (
	"context"
	"sync"
	"time"

	"tmt/internal/entity"
)

// assistTarget is the target order to assist
type assistTarget struct {
	*WSFutureTrade       // use for place order
	*entity.FutureOrder  // base order
	halfAutomationOption // option for assist trader
}

// toFinishOrder will return a order to finish the assist order
func (o *assistTarget) toFinishOrder(price float64) *entity.FutureOrder {
	order := &entity.FutureOrder{
		Code: o.Code,
		BaseOrder: entity.BaseOrder{
			Quantity: o.Quantity,
			Price:    price,
		},
	}
	switch o.Action {
	case entity.ActionBuy:
		order.Action = entity.ActionSell
	case entity.ActionSell:
		order.Action = entity.ActionBuy
	}
	return order
}

// assistTrader is a trader to assist the main order
type assistTrader struct {
	*assistTarget                                      // assist target
	ctx                context.Context                 // context from gin context
	finishOrderMap     map[string]*entity.FutureOrder  // map of finish order
	finishOrderMapLock sync.RWMutex                    // lock for finishOrderMap
	waitingOrder       *entity.FutureOrder             // waiting order
	tickChan           chan *entity.RealTimeFutureTick // tick channel
	done               bool                            // done flag
	tradeOutPrice      float64                         // trade out price
	holdTimes          int                             // hold times
}

// newAssistTrader will return a assist trader
func newAssistTrader(ctx context.Context, target *assistTarget) *assistTrader {
	a := &assistTrader{
		ctx:            ctx,
		assistTarget:   target,
		finishOrderMap: make(map[string]*entity.FutureOrder),
		tickChan:       make(chan *entity.RealTimeFutureTick),
		tradeOutPrice:  target.Price,
		holdTimes:      3,
	}

	go a.processTick()

	a.SubscribeTopic(topicOrderStatus, a.updateOrderStatus)
	a.SendToClient(newAssistStatusProto(true))

	return a
}

// updateOrderStatus will update order status
func (a *assistTrader) updateOrderStatus(o *entity.FutureOrder) {
	defer a.finishOrderMapLock.Unlock()
	a.finishOrderMapLock.Lock()

	if _, ok := a.finishOrderMap[o.OrderID]; ok {
		a.finishOrderMap[o.OrderID] = o
		if a.waitingOrder != nil && !o.Cancellable() && a.waitingOrder.OrderID == o.OrderID {
			a.waitingOrder = nil
		}
	}
}

// processTick will process tick
func (a *assistTrader) processTick() {
	for {
		tick, ok := <-a.tickChan
		if !ok {
			return
		}

		if a.waitingOrder != nil {
			continue
		}

		if a.isAssistDone() {
			continue
		}

		switch a.AutomationType {
		case AutomationByBalance:
			a.checkByBalance(tick)
		case AutomationByTimePeriod:
			a.checkByTime(tick)
		case AutomationByTimePeriodAndBalance:
			a.checkByTimeAndBalance(tick)
		}
	}
}

func (a *assistTrader) isAssistDone() bool {
	if a.done {
		return true
	}

	var endQty int64
	a.finishOrderMapLock.RLock()
	for _, o := range a.finishOrderMap {
		if o.Status == entity.StatusFilled {
			endQty += o.Quantity
		}
	}
	a.finishOrderMapLock.RUnlock()

	if endQty == a.Quantity {
		a.UnSubscribeTopic(topicOrderStatus, a.updateOrderStatus)
		a.PublishTopicEvent(topicAssistDone, a.OrderID)
		a.SendToClient(newAssistStatusProto(false))
		a.done = true
		return true
	}
	return false
}

// halfAutomationOption is the option for assist trader
func (a *assistTrader) checkByTime(tick *entity.RealTimeFutureTick) {
	if time.Since(a.TradeTime) > time.Duration(a.ByTimePeriod)*time.Minute {
		a.placeAssistOrder(tick.Close)
	}
}

// halfAutomationOption is the option for assist trader
func (a *assistTrader) checkByBalance(tick *entity.RealTimeFutureTick) {
	defer func() {
		a.tradeOutPrice = tick.Close
	}()

	switch a.Action {
	case entity.ActionBuy:
		switch {
		case tick.Close >= a.Price+a.ByBalanceHigh:
			if tick.Close >= a.tradeOutPrice && a.holdTimes > 0 {
				a.holdTimes--
				return
			}
			a.placeAssistOrder(tick.Close)

		case tick.Close <= a.Price+a.ByBalanceLow:
			a.placeAssistOrder(tick.Close)
		}

	case entity.ActionSell:
		switch {
		case tick.Close <= a.Price-a.ByBalanceHigh:
			if tick.Close <= a.tradeOutPrice && a.holdTimes > 0 {
				a.holdTimes--
				return
			}
			a.placeAssistOrder(tick.Close)

		case tick.Close >= a.Price-a.ByBalanceLow:
			a.placeAssistOrder(tick.Close)
		}
	}
}

// halfAutomationOption is the option for assist trader
func (a *assistTrader) checkByTimeAndBalance(tick *entity.RealTimeFutureTick) {
	defer func() {
		a.tradeOutPrice = tick.Close
	}()

	switch a.Action {
	case entity.ActionBuy:
		switch {
		case tick.Close >= a.Price+a.ByBalanceHigh:
			if tick.Close >= a.tradeOutPrice && a.holdTimes > 0 {
				a.holdTimes--
				return
			}
			a.placeAssistOrder(tick.Close)

		case tick.Close <= a.Price+a.ByBalanceLow:
			a.placeAssistOrder(tick.Close)
		}

	case entity.ActionSell:
		switch {
		case tick.Close <= a.Price-a.ByBalanceHigh:
			if tick.Close <= a.tradeOutPrice && a.holdTimes > 0 {
				a.holdTimes--
				return
			}
			a.placeAssistOrder(tick.Close)

		case tick.Close >= a.Price-a.ByBalanceLow:
			a.placeAssistOrder(tick.Close)
		}
	}

	if time.Since(a.TradeTime) > time.Duration(a.ByTimePeriod)*time.Minute {
		a.placeAssistOrder(tick.Close)
	}
}

// placeAssistOrder will place assist order
func (a *assistTrader) placeAssistOrder(close float64) {
	if o := a.placeOrder(a.toFinishOrder(close)); o != nil {
		a.finishOrderMapLock.Lock()
		a.finishOrderMap[o.OrderID] = o
		a.finishOrderMapLock.Unlock()
		a.waitingOrder = o
		a.PublishTopicEvent(topicPlaceOrder, o)
	}
}

// getTickChan will return tick channel
func (a *assistTrader) getTickChan() chan *entity.RealTimeFutureTick {
	return a.tickChan
}
