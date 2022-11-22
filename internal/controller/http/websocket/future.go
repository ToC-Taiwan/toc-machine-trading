package websocket

import (
	"errors"
	"time"

	"tmt/internal/entity"
)

type AutomationType int

const (
	// WSPickStock -
	AutomationByTime AutomationType = iota + 1
	// WSFuture -
	AutomationByBalance
)

type futureOrder struct {
	Code   string             `json:"code"`
	Action entity.OrderAction `json:"action"`
	Price  float64            `json:"price"`
	Qty    int64              `json:"qty"`

	HalfAutomation bool           `json:"half_automation"`
	AutomationType AutomationType `json:"automation_type"`
}

type tradeRate struct {
	OutRate int64 `json:"out_rate"`
	OutChg  int64 `json:"out_chg"`
	InRate  int64 `json:"in_rate"`
	InChg   int64 `json:"in_chg"`
}

type futurePosition struct {
	Position []*entity.FuturePosition `json:"position"`
}

func (w *WSRouter) processTrade(clientMsg clientMsg) {
	if clientMsg.FutureOrder == nil {
		return
	}

	if !w.o.IsFutureTradeTime() {
		w.msgChan <- errMsg{ErrMsg: "Not trade time"}
		return
	}

	order := &entity.FutureOrder{
		Code: clientMsg.FutureOrder.Code,
		BaseOrder: entity.BaseOrder{
			Action:   clientMsg.FutureOrder.Action,
			Quantity: clientMsg.FutureOrder.Qty,
			Price:    clientMsg.FutureOrder.Price,
		},
		Manual: true,
	}

	var err error
	switch order.Action {
	case entity.ActionBuy:
		order.OrderID, order.Status, err = w.o.BuyFuture(order)
		if err != nil {
			w.msgChan <- errMsg{ErrMsg: err.Error()}
			log.Error(err)
			return
		}

	case entity.ActionSell:
		order.OrderID, order.Status, err = w.o.SellFuture(order)
		if err != nil {
			w.msgChan <- errMsg{ErrMsg: err.Error()}
			log.Error(err)
			return
		}
	}

	w.orderLock.Lock()
	order.TradeTime = time.Now()
	w.futureOrderMap[order.OrderID] = order
	log.Infof("New Order: %s", order.FutureOrderStatusString())
	w.orderLock.Unlock()
}

func (w *WSRouter) sendFuture() {
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

	w.s.NewFutureRealTimeConnection(tickChan, w.connectionID)
	w.s.NewOrderStatusConnection(orderStatusChan, w.connectionID)

	<-w.ctx.Done()

	w.s.DeleteFutureRealTimeConnection(w.connectionID)
	w.s.DeleteOrderStatusConnection(w.connectionID)
}

func (w *WSRouter) processTickArr(tickChan chan *entity.RealTimeFutureTick) {
	var tickArr []*entity.RealTimeFutureTick
	var lastTradeRate *tradeRate
	for {
		tick, ok := <-tickChan
		if !ok {
			return
		}
		tickArr = append(tickArr, tick)

		var outVolume, inVolume int64
		for i := len(tickArr) - 1; i >= 0; i-- {
			if time.Since(tickArr[i].TickTime) > 15*time.Second {
				tickArr = tickArr[i+1:]
				break
			}

			switch tickArr[i].TickType {
			case 1:
				outVolume += tickArr[i].Volume
			case 2:
				inVolume += tickArr[i].Volume
			}
		}
		rate := &tradeRate{
			OutRate: outVolume,
			OutChg:  outVolume - lastTradeRate.OutRate,
			InRate:  inVolume,
			InChg:   inVolume - lastTradeRate.InRate,
		}
		lastTradeRate = rate
		w.msgChan <- rate
		w.msgChan <- tick
	}
}

func (w *WSRouter) processOrderStatus(orderStatusChan chan interface{}) {
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
			w.orderLock.Unlock()
		}
	}
}

func (w *WSRouter) sendTradeIndex() {
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

func (w *WSRouter) sendPosition() {
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

func (w *WSRouter) cancelOverTimeOrder() {
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
					log.Warnf("Delete %s", order.FutureOrderStatusString())
				} else if time.Since(order.TradeTime) > 10*time.Second && cancelOrderMap[id] == nil {
					if e := w.cancelOrderByID(id); e != nil {
						w.msgChan <- errMsg{ErrMsg: e.Error()}
					}
					cancelOrderMap[id] = order
					log.Warnf("%s timeout, cancel it", order.FutureOrderStatusString())
				}
			}
			w.orderLock.Unlock()
		}
	}
}

func (w *WSRouter) cancelOrderByID(orderID string) error {
	_, s, err := w.o.CancelFutureOrderID(orderID)
	if err != nil {
		return err
	}

	if s != entity.StatusCancelled {
		return errors.New("cancel order failed")
	}

	return nil
}

func (w *WSRouter) generateTradeIndex() *entity.TradeIndex {
	return w.s.GetTradeIndex()
}

func (w *WSRouter) generatePosition() ([]*entity.FuturePosition, error) {
	position, err := w.o.GetFuturePosition()
	if err != nil {
		return nil, err
	}
	return position, nil
}
