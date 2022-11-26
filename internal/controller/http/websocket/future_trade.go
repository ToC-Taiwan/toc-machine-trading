package websocket

import (
	"encoding/json"
	"errors"
	"sync"
	"time"

	"tmt/internal/entity"
	"tmt/internal/usecase"

	"github.com/gin-gonic/gin"
)

type WSFutureTrade struct {
	*WSRouter

	s usecase.Stream
	o usecase.Order

	futureOrderMap map[string]*entity.FutureOrder
	assistOrderMap map[string]*entity.FutureOrder
	assistOption   map[string]HalfAutomationOption
	orderLock      sync.Mutex
}

type futureTradeClientMsg struct {
	FutureOrder *futureOrder `json:"future_order"`
}

// StartWSFutureTrade -.
func StartWSFutureTrade(c *gin.Context, s usecase.Stream, o usecase.Order) {
	w := &WSFutureTrade{
		s:              s,
		o:              o,
		futureOrderMap: make(map[string]*entity.FutureOrder),
		assistOrderMap: make(map[string]*entity.FutureOrder),
		assistOption:   make(map[string]HalfAutomationOption),
		WSRouter:       NewWSRouter(c),
	}

	forwardChan := make(chan []byte)
	go func() {
		for {
			msg, ok := <-forwardChan
			if !ok {
				return
			}

			var fMsg futureTradeClientMsg
			if err := json.Unmarshal(msg, &fMsg); err != nil {
				w.msgChan <- errMsg{ErrMsg: err.Error()}
				continue
			}
			w.processTrade(fMsg)
		}
	}()
	w.read(forwardChan)
	w.sendFuture()
}

type futureOrder struct {
	Code   string               `json:"code"`
	Action entity.OrderAction   `json:"action"`
	Price  float64              `json:"price"`
	Qty    int64                `json:"qty"`
	Option HalfAutomationOption `json:"option"`
}

type AutomationType int

const (
	AutomationByBalance AutomationType = iota + 1
	AutomationByTimePeriod
)

type HalfAutomationOption struct {
	HalfAutomation bool           `json:"half_automation"`
	AutomationType AutomationType `json:"automation_type"`
	ByBalanceHigh  float64        `json:"by_balance_high"`
	ByBalanceLow   float64        `json:"by_balance_low"`
	ByTimePeriod   int64          `json:"by_time_period"`
}

type periodTradeVolume struct {
	FirstPeriod  entity.OutInVolume `json:"first_period"`
	SecondPeriod entity.OutInVolume `json:"second_period"`
	ThirdPeriod  entity.OutInVolume `json:"third_period"`
	FourthPeriod entity.OutInVolume `json:"fourth_period"`
}

type futurePosition struct {
	Position []*entity.FuturePosition `json:"position"`
}

func (w *WSFutureTrade) processTrade(clientMsg futureTradeClientMsg) {
	if !w.o.IsFutureTradeTime() {
		w.msgChan <- errMsg{ErrMsg: "Not trade time"}
		return
	}

	// w.orderLock.Lock()
	// defer w.orderLock.Unlock()
	// if len(w.assistOrderMap) != 0 {
	// 	w.msgChan <- errMsg{ErrMsg: "assist order is processing"}
	// 	return
	// }

	order := &entity.FutureOrder{
		Code: clientMsg.FutureOrder.Code,
		BaseOrder: entity.BaseOrder{
			Action:   clientMsg.FutureOrder.Action,
			Quantity: clientMsg.FutureOrder.Qty,
			Price:    clientMsg.FutureOrder.Price,
		},
	}

	var err error
	switch order.Action {
	case entity.ActionBuy:
		order.OrderID, order.Status, err = w.o.BuyFuture(order)
		if err != nil {
			w.msgChan <- errMsg{ErrMsg: err.Error()}
			return
		}

	case entity.ActionSell:
		order.OrderID, order.Status, err = w.o.SellFuture(order)
		if err != nil {
			w.msgChan <- errMsg{ErrMsg: err.Error()}
			return
		}
	}

	w.orderLock.Lock()
	order.TradeTime = time.Now()
	if clientMsg.FutureOrder.Option.HalfAutomation {
		w.assistOrderMap[order.OrderID] = order
		w.assistOption[order.OrderID] = clientMsg.FutureOrder.Option
	} else {
		w.futureOrderMap[order.OrderID] = order
	}
	w.orderLock.Unlock()
}

func (w *WSFutureTrade) sendFuture() {
	snapshot, err := w.s.GetFutureSnapshotByCode(w.s.GetMainFutureCode())
	if err != nil {
		w.msgChan <- errMsg{ErrMsg: err.Error()}
	} else {
		w.msgChan <- snapshot.ToRealTimeFutureTick()
	}

	tickChan := make(chan *entity.RealTimeFutureTick)
	orderStatusChan := make(chan interface{})
	go w.processTickArr(tickChan)
	go w.processOrderStatus(orderStatusChan)

	go w.sendTradeIndex()
	go w.sendPosition()
	go w.cancelOverTimeOrder()
	// go w.assistTrader()

	w.s.NewFutureRealTimeConnection(tickChan, w.connectionID)
	w.s.NewOrderStatusConnection(orderStatusChan, w.connectionID)

	<-w.ctx.Done()

	w.s.DeleteFutureRealTimeConnection(w.connectionID)
	w.s.DeleteOrderStatusConnection(w.connectionID)
}

func (w *WSFutureTrade) processTickArr(tickChan chan *entity.RealTimeFutureTick) {
	var tickArr entity.RealTimeFutureTickArr
	for {
		tick, ok := <-tickChan
		if !ok {
			return
		}
		tickArr = append(tickArr, tick)
		w.msgChan <- tick

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

		w.msgChan <- periodTradeVolume{
			FirstPeriod:  firstPeriod.GetOutInVolume(),
			SecondPeriod: secondPeriod.GetOutInVolume(),
			ThirdPeriod:  thirdPeriod.GetOutInVolume(),
			FourthPeriod: fourthPeriod.GetOutInVolume(),
		}
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
			if cache := w.futureOrderMap[o.OrderID]; cache != nil && cache.Status != o.Status {
				w.msgChan <- o
				o.TradeTime = cache.TradeTime
				w.futureOrderMap[o.OrderID] = o
			}

			// if cache := w.assistOrderMap[o.OrderID]; cache != nil && cache.Status != o.Status {
			// 	w.msgChan <- o
			// 	o.TradeTime = cache.TradeTime
			// 	w.futureOrderMap[o.OrderID] = o
			// }
			w.orderLock.Unlock()
		}
	}
}

func (w *WSFutureTrade) sendTradeIndex() {
	w.msgChan <- w.generateTradeIndex()

	for {
		select {
		case <-w.ctx.Done():
			return

		case <-time.After(5 * time.Second):
			w.msgChan <- w.generateTradeIndex()
		}
	}
}

func (w *WSFutureTrade) sendPosition() {
	if position, err := w.generatePosition(); err != nil {
		w.msgChan <- errMsg{ErrMsg: err.Error()}
	} else {
		w.msgChan <- &futurePosition{position}
	}

	for {
		select {
		case <-w.ctx.Done():
			return

		case <-time.After(5 * time.Second):
			if position, err := w.generatePosition(); err != nil {
				w.msgChan <- errMsg{ErrMsg: err.Error()}
			} else {
				w.msgChan <- &futurePosition{position}
			}
		}
	}
}

func (w *WSFutureTrade) cancelOverTimeOrder() {
	cancelOrderMap := make(map[string]*entity.FutureOrder)
	for {
		select {
		case <-w.ctx.Done():
			return

		case <-time.After(time.Second):
			w.orderLock.Lock()
			for id, order := range w.futureOrderMap {
				if !order.Cancellabel() {
					delete(w.futureOrderMap, id)
					delete(cancelOrderMap, id)
				} else if time.Since(order.TradeTime) > 10*time.Second && cancelOrderMap[id] == nil {
					if e := w.cancelOrderByID(id); e != nil {
						w.msgChan <- errMsg{ErrMsg: e.Error()}
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
	return w.s.GetTradeIndex()
}

func (w *WSFutureTrade) generatePosition() ([]*entity.FuturePosition, error) {
	position, err := w.o.GetFuturePosition()
	if err != nil {
		return nil, err
	}
	return position, nil
}

// func (w *WSFutureTrade) assistTrader() {
// 	for {
// 		select {
// 		case <-w.ctx.Done():
// 			return

// 		case <-time.After(time.Second):
// 			w.orderLock.Lock()
// 			for id, order := range w.assistOrderMap {
// 			}
// 			w.orderLock.Unlock()
// 		}
// 	}
// }
