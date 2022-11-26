// Package future package future
package future

import (
	"encoding/json"
	"errors"
	"sync"
	"time"

	"tmt/internal/entity"
	"tmt/internal/usecase"

	"tmt/internal/controller/http/websocket"

	"github.com/gin-gonic/gin"
)

type WSFutureTrade struct {
	*websocket.WSRouter

	s usecase.Stream
	o usecase.Order

	futureOrderMap map[string]*entity.FutureOrder
	orderLock      sync.Mutex

	assistTrader         *assistTrader
	assistTraderTickChan chan *entity.RealTimeFutureTick
}

// StartWSFutureTrade -.
func StartWSFutureTrade(c *gin.Context, s usecase.Stream, o usecase.Order) {
	w := &WSFutureTrade{
		s:              s,
		o:              o,
		futureOrderMap: make(map[string]*entity.FutureOrder),
		WSRouter:       websocket.NewWSRouter(c),
	}
	w.assistTrader = newAssistTrader(w.Ctx(), o)
	w.assistTraderTickChan = w.assistTrader.getTickChan()

	forwardChan := make(chan []byte)
	go func() {
		for {
			msg, ok := <-forwardChan
			if !ok {
				return
			}

			var fMsg futureTradeClientMsg
			if err := json.Unmarshal(msg, &fMsg); err != nil {
				w.SendToClient(errMsg{ErrMsg: err.Error()})
				continue
			}
			w.processTrade(fMsg)
		}
	}()
	go w.sendFuture()
	w.ReadFromClient(forwardChan)
}

func (w *WSFutureTrade) processTrade(clientMsg futureTradeClientMsg) {
	if !w.o.IsFutureTradeTime() {
		w.SendToClient(errMsg{ErrMsg: "Not trade time"})
		return
	}

	if !w.assistTrader.isAssistDone() {
		w.SendToClient(errMsg{ErrMsg: "Assist trader is running"})
		return
	}

	order := &entity.FutureOrder{
		Code: clientMsg.Code,
		BaseOrder: entity.BaseOrder{
			Action:   clientMsg.Action,
			Quantity: clientMsg.Qty,
			Price:    clientMsg.Price,
		},
	}

	if e := w.placeOrder(order); e != nil {
		w.SendToClient(errMsg{ErrMsg: e.Error()})
		return
	}

	w.orderLock.Lock()
	order.TradeTime = time.Now()
	if clientMsg.Option.AutomationType != AutomationNone {
		w.assistTrader.addAssistOrder(order, clientMsg.Option)
	} else {
		w.futureOrderMap[order.OrderID] = order
	}
	w.orderLock.Unlock()
}

func (w *WSFutureTrade) placeOrder(order *entity.FutureOrder) error {
	var err error
	switch order.Action {
	case entity.ActionBuy:
		order.OrderID, order.Status, err = w.o.BuyFuture(order)
		if err != nil {
			return err
		}

	case entity.ActionSell:
		order.OrderID, order.Status, err = w.o.SellFuture(order)
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *WSFutureTrade) sendFuture() {
	snapshot, err := w.s.GetFutureSnapshotByCode(w.s.GetMainFutureCode())
	if err != nil {
		w.SendToClient(errMsg{ErrMsg: err.Error()})
	} else {
		w.SendToClient(snapshot.ToRealTimeFutureTick())
	}

	tickChan := make(chan *entity.RealTimeFutureTick)
	orderStatusChan := make(chan interface{})
	go w.processTickArr(tickChan)
	go w.processOrderStatus(orderStatusChan)

	go w.sendTradeIndex()
	go w.sendPosition()
	go w.cancelOverTimeOrder()

	w.s.NewFutureRealTimeConnection(tickChan, w.ConnectionID)
	w.s.NewOrderStatusConnection(orderStatusChan, w.ConnectionID)

	<-w.Ctx().Done()

	w.s.DeleteFutureRealTimeConnection(w.ConnectionID)
	w.s.DeleteOrderStatusConnection(w.ConnectionID)
}

func (w *WSFutureTrade) processTickArr(tickChan chan *entity.RealTimeFutureTick) {
	var tickArr entity.RealTimeFutureTickArr
	for {
		tick, ok := <-tickChan
		if !ok {
			close(w.assistTraderTickChan)
			return
		}
		tickArr = append(tickArr, tick)
		w.SendToClient(tick)
		w.assistTraderTickChan <- tick

		var firstPeriod, secondPeriod, thirdPeriod, fourthPeriod entity.RealTimeFutureTickArr
	L:
		for i := len(tickArr) - 1; i >= 0; i-- {
			switch {
			case time.Since(tickArr[i].TickTime) <= 10*time.Second:
				fourthPeriod = append(fourthPeriod, tickArr[i])
				thirdPeriod = append(thirdPeriod, tickArr[i])
				secondPeriod = append(secondPeriod, tickArr[i])
				firstPeriod = append(firstPeriod, tickArr[i])

			case time.Since(tickArr[i].TickTime) <= 20*time.Second:
				fourthPeriod = append(fourthPeriod, tickArr[i])
				thirdPeriod = append(thirdPeriod, tickArr[i])
				secondPeriod = append(secondPeriod, tickArr[i])

			case time.Since(tickArr[i].TickTime) <= 30*time.Second:
				fourthPeriod = append(fourthPeriod, tickArr[i])
				thirdPeriod = append(thirdPeriod, tickArr[i])

			case time.Since(tickArr[i].TickTime) <= 40*time.Second:
				fourthPeriod = append(fourthPeriod, tickArr[i])

			default:
				tickArr = tickArr[i+1:]
				break L
			}
		}

		w.SendToClient(periodTradeVolume{
			FirstPeriod:  firstPeriod.GetOutInVolume(),
			SecondPeriod: secondPeriod.GetOutInVolume(),
			ThirdPeriod:  thirdPeriod.GetOutInVolume(),
			FourthPeriod: fourthPeriod.GetOutInVolume(),
		})
	}
}

func (w *WSFutureTrade) processOrderStatus(orderStatusChan chan interface{}) {
	for {
		order, ok := <-orderStatusChan
		if !ok {
			return
		}

		if o, ok := order.(*entity.FutureOrder); ok {
			w.orderLock.Lock()
			if cache, ok := w.futureOrderMap[o.OrderID]; ok {
				if cache.Status != o.Status {
					w.SendToClient(o)
					o.TradeTime = cache.TradeTime
					w.futureOrderMap[o.OrderID] = o
				}
			} else {
				if a := w.assistTrader.updateOrderStatus(o); a != nil {
					w.SendToClient(a)
				}
			}
			w.orderLock.Unlock()
		}
	}
}

func (w *WSFutureTrade) sendTradeIndex() {
	w.SendToClient(w.generateTradeIndex())

	for {
		select {
		case <-w.Ctx().Done():
			return

		case <-time.After(5 * time.Second):
			w.SendToClient(w.generateTradeIndex())
		}
	}
}

func (w *WSFutureTrade) sendPosition() {
	if position, err := w.generatePosition(); err != nil {
		w.SendToClient(errMsg{ErrMsg: err.Error()})
	} else {
		w.SendToClient(&futurePosition{position})
	}

	for {
		select {
		case <-w.Ctx().Done():
			return

		case <-time.After(5 * time.Second):
			if position, err := w.generatePosition(); err != nil {
				w.SendToClient(errMsg{ErrMsg: err.Error()})
			} else {
				w.SendToClient(&futurePosition{position})
			}
		}
	}
}

func (w *WSFutureTrade) cancelOverTimeOrder() {
	cancelOrderMap := make(map[string]*entity.FutureOrder)
	for {
		select {
		case <-w.Ctx().Done():
			return

		case <-time.After(time.Second):
			w.orderLock.Lock()
			for id, order := range w.futureOrderMap {
				if !order.Cancellabel() {
					delete(w.futureOrderMap, id)
					delete(cancelOrderMap, id)
				} else if time.Since(order.TradeTime) > 10*time.Second && cancelOrderMap[id] == nil {
					if e := w.cancelOrderByID(id); e != nil {
						w.SendToClient(errMsg{ErrMsg: e.Error()})
					}
					cancelOrderMap[id] = order
				}
			}
			w.orderLock.Unlock()
		}
	}
}

func (w *WSFutureTrade) cancelOrderByID(orderID string) error {
	_, s, err := w.o.CancelFutureOrderID(orderID)
	if err != nil {
		return err
	}

	if s != entity.StatusCancelled {
		return errors.New("cancel order failed")
	}

	return nil
}

func (w *WSFutureTrade) generateTradeIndex() *entity.TradeIndex {
	t := w.s.GetTradeIndex()
	switch {
	case t.Nasdaq == nil, t.NF == nil, t.TSE == nil, t.OTC == nil:
		time.Sleep(time.Second)
		return w.generateTradeIndex()
	default:
		return t
	}
}

func (w *WSFutureTrade) generatePosition() ([]*entity.FuturePosition, error) {
	position, err := w.o.GetFuturePosition()
	if err != nil {
		return nil, err
	}
	return position, nil
}
