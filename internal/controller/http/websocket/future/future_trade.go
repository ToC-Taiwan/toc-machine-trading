// Package future package future
package future

import (
	"encoding/json"
	"sync"
	"time"

	"tmt/internal/entity"
	"tmt/internal/usecase"
	"tmt/pkg/eventbus"

	"tmt/internal/controller/http/websocket"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WSFutureTrade struct {
	*websocket.WSRouter // ws router
	*eventbus.Bus       // event bus

	s usecase.Stream  // stream
	o usecase.Order   // order
	h usecase.History // history

	// save tick chan for assist
	assistTickChanMap     map[string]chan *entity.RealTimeFutureTick
	assistTickChanMapLock sync.RWMutex

	// if waiting manual is not nil, will not accept new order
	waitingList *waitingList

	// save manual order or order from assist
	orderMap           map[string]*entity.FutureOrder
	orderMapLock       sync.Mutex
	cancelOrderMap     map[string]*entity.FutureOrder
	cancelOrderMapLock sync.Mutex

	// limit one process at a time
	processLock sync.Mutex

	// save assist target for assist, if assist order is done, will start to send tick to assist target
	// then delete from map
	assistTargetWaitingMap     map[string]*assistTarget
	assistTargetWaitingMapLock sync.Mutex
}

// StartWSFutureTrade - Start ws future trade with one time bus
func StartWSFutureTrade(c *gin.Context, s usecase.Stream, o usecase.Order, h usecase.History) {
	w := &WSFutureTrade{
		s:                      s,
		o:                      o,
		h:                      h,
		assistTickChanMap:      make(map[string]chan *entity.RealTimeFutureTick),
		assistTargetWaitingMap: make(map[string]*assistTarget),
		orderMap:               make(map[string]*entity.FutureOrder),
		cancelOrderMap:         make(map[string]*entity.FutureOrder),
		WSRouter:               websocket.NewWSRouter(c),
		Bus:                    eventbus.Get(uuid.NewString()),
		waitingList:            newWaitingList(),
	}

	forwardChan := make(chan []byte)
	go func() {
		for {
			msg, ok := <-forwardChan
			if !ok {
				return
			}

			var fMsg clientOrder
			if err := json.Unmarshal(msg, &fMsg); err != nil {
				w.SendToClient(newErrMessageProto(errUnmarshal))
				continue
			}
			w.processClientOrder(fMsg)
		}
	}()

	go w.sendFuture()
	go w.checkAssistTargetStatus()

	w.SubscribeTopic(topicAssistDone, w.closeDoneChan)
	w.SubscribeTopic(topicPlaceOrder, w.addOrderFromAssist)

	w.ReadFromClient(forwardChan)
}

func (w *WSFutureTrade) processClientOrder(client clientOrder) {
	defer w.processLock.Unlock()
	w.processLock.Lock()

	switch {
	case !w.o.IsFutureTradeTime():
		w.SendToClient(newErrMessageProto(errNotTradeTime))
		return
	case !w.waitingList.empty():
		w.SendToClient(newErrMessageProto(errNotFilled))
		return
	case !w.isAssistingFinish():
		w.SendToClient(newErrMessageProto(errAssitingIsNotFinished))
		return
	case client.Option.AutomationType != AutomationNone && client.Qty > 4:
		w.SendToClient(newErrMessageProto(errAssistNotSupport))
		return
	}

	for _, order := range client.toFutureOrderArr() {
		o := w.placeOrder(order)
		if o == nil {
			w.SendToClient(newErrMessageProto(errPlaceOrder))
			return
		}
		w.waitingList.add(o)
		if client.Option.AutomationType != AutomationNone {
			// save assist target, wait for order status update
			w.assistTargetWaitingMapLock.Lock()
			w.assistTargetWaitingMap[o.OrderID] = &assistTarget{
				WSFutureTrade:        w,
				FutureOrder:          o,
				halfAutomationOption: client.Option,
			}
			w.assistTargetWaitingMapLock.Unlock()
		} else {
			// save manual order, it has timeout
			w.orderMapLock.Lock()
			w.orderMap[o.OrderID] = o
			w.orderMapLock.Unlock()
		}
	}
}

func (w *WSFutureTrade) isAssistingFinish() bool {
	defer w.assistTickChanMapLock.RUnlock()
	w.assistTickChanMapLock.RLock()
	return len(w.assistTickChanMap) == 0
}

func (w *WSFutureTrade) closeDoneChan(orderID string) {
	w.assistTickChanMapLock.Lock()
	close(w.assistTickChanMap[orderID])
	delete(w.assistTickChanMap, orderID)
	w.assistTickChanMapLock.Unlock()
}

func (w *WSFutureTrade) addOrderFromAssist(o *entity.FutureOrder) {
	w.orderMapLock.Lock()
	w.orderMap[o.OrderID] = o
	w.orderMapLock.Unlock()
}

func (w *WSFutureTrade) sendFuture() {
	w.sendFutureDetail()
	w.sendLatestKbar()
	w.sendFutureSnapshot()

	tickChan := make(chan *entity.RealTimeFutureTick)
	orderStatusChan := make(chan interface{})
	go w.processTickArr(tickChan)
	go w.processOrderStatus(orderStatusChan)

	go w.sendTradeIndex()
	go w.sendPosition()

	connectionID := uuid.New().String()
	w.s.NewFutureRealTimeConnection(tickChan, connectionID)
	w.s.NewOrderStatusConnection(orderStatusChan, connectionID)

	<-w.Ctx().Done()

	w.s.DeleteFutureRealTimeConnection(connectionID)
	w.s.DeleteOrderStatusConnection(connectionID)
}

func (w *WSFutureTrade) checkAssistTargetStatus() {
	for {
		select {
		case <-w.Ctx().Done():
			return

		case <-time.After(500 * time.Millisecond):
			w.assistTargetWaitingMapLock.Lock()
			for orderID, a := range w.assistTargetWaitingMap {
				if a.Status == entity.StatusFilled {
					assist := newAssistTrader(w.Ctx(), a)
					w.assistTickChanMapLock.Lock()
					w.assistTickChanMap[orderID] = assist.getTickChan()
					w.assistTickChanMapLock.Unlock()
				}

				if !a.Cancellable() {
					delete(w.assistTargetWaitingMap, orderID)
				}
			}
			w.assistTargetWaitingMapLock.Unlock()
		}
	}
}

func (w *WSFutureTrade) placeOrder(order *entity.FutureOrder) *entity.FutureOrder {
	var err error
	switch order.Action {
	case entity.ActionBuy:
		order.OrderID, order.Status, err = w.o.BuyFuture(order)
		if err != nil {
			return nil
		}

	case entity.ActionSell:
		order.OrderID, order.Status, err = w.o.SellFuture(order)
		if err != nil {
			return nil
		}
	}

	// order trade time is not set by order service, set it here
	order.TradeTime = time.Now()
	return order
}

func (w *WSFutureTrade) processTickArr(tickChan chan *entity.RealTimeFutureTick) {
	for {
		tick, ok := <-tickChan
		if !ok {
			return
		}
		w.sendTickToAssit(tick)
		w.SendToClient(newFutureTickProto(tick))
	}
}

func (w *WSFutureTrade) processOrderStatus(orderStatusChan chan interface{}) {
	finishedOrderMap := make(map[string]*entity.FutureOrder)
	for {
		order, ok := <-orderStatusChan
		if !ok {
			return
		}

		if o, ok := order.(*entity.FutureOrder); ok {
			if finishedOrderMap[o.OrderID] != nil {
				continue
			}

			w.updateCacheOrder(o)
			if !o.Cancellable() {
				finishedOrderMap[o.OrderID] = o
				if w.waitingList.orderIDExist(o.OrderID) {
					w.waitingList.remove(o.OrderID)
				}
			} else {
				w.cancelOverTimeOrder(o)
			}
		}
	}
}

// cancelOverTimeOrder cancel order if it is not cancelled or filled, and also update order from assist
func (w *WSFutureTrade) updateCacheOrder(o *entity.FutureOrder) {
	w.orderMapLock.Lock()
	cache, ok := w.orderMap[o.OrderID]
	if ok {
		if cache.Status != o.Status {
			w.SendToClient(newFutureOrderProto(o))
		}
		o.TradeTime = cache.TradeTime
		w.orderMap[o.OrderID] = o
		w.PublishTopicEvent(topicOrderStatus, o) // publish updated order to assist
	}
	w.orderMapLock.Unlock()
	if !ok {
		w.updateAssistTargetWaitingOrder(o)
	}
}

// updateAssistTargetWaitingOrder if filled, assist trader will start
func (w *WSFutureTrade) updateAssistTargetWaitingOrder(o *entity.FutureOrder) {
	w.assistTargetWaitingMapLock.Lock()
	if a, ok := w.assistTargetWaitingMap[o.OrderID]; ok {
		if a.Status != o.Status {
			w.SendToClient(newFutureOrderProto(o))
		}
		o.TradeTime = a.TradeTime
		a.FutureOrder = o
		w.assistTargetWaitingMap[o.OrderID] = a
	}
	w.assistTargetWaitingMapLock.Unlock()
}

func (w *WSFutureTrade) cancelOverTimeOrder(o *entity.FutureOrder) {
	defer w.cancelOrderMapLock.Unlock()
	w.cancelOrderMapLock.Lock()

	if o.Cancellable() && time.Since(o.TradeTime) > 5*time.Second && w.cancelOrderMap[o.OrderID] == nil {
		id, s, err := w.o.CancelFutureOrderID(o.OrderID)
		if err != nil || s != entity.StatusCancelled || id == "" {
			w.SendToClient(newErrMessageProto(errCancelOrderFailed))
			return
		}
		w.cancelOrderMap[o.OrderID] = o
	}
}

func (w *WSFutureTrade) sendTickToAssit(tick *entity.RealTimeFutureTick) {
	w.assistTickChanMapLock.RLock()
	for _, v := range w.assistTickChanMap {
		v <- tick
	}
	w.assistTickChanMapLock.RUnlock()
}

func (w *WSFutureTrade) sendFutureDetail() {
	w.SendToClient(newFutureDetailProto(w.s.GetMainFuture()))
}

func (w *WSFutureTrade) sendFutureSnapshot() {
	snapshot, err := w.s.GetFutureSnapshotByCode(w.s.GetMainFuture().Code)
	if err != nil {
		w.SendToClient(newErrMessageProto(errGetSnapshot))
	} else {
		w.SendToClient(newFutureTickProto(snapshot.ToRealTimeFutureTick()))
	}
}

func (w *WSFutureTrade) sendTradeIndex() {
	w.SendToClient(newTradeIndexProto(w.generateTradeIndex()))

	for {
		select {
		case <-w.Ctx().Done():
			return

		case <-time.After(5 * time.Second):
			w.SendToClient(newTradeIndexProto(w.generateTradeIndex()))
		}
	}
}

func (w *WSFutureTrade) sendPosition() {
	if position, err := w.generatePosition(); err != nil {
		w.SendToClient(newErrMessageProto(errGetPosition))
	} else {
		w.SendToClient(newFuturePositionProto(position))
	}

	for {
		select {
		case <-w.Ctx().Done():
			return

		case <-time.After(10 * time.Second):
			if !w.o.IsFutureTradeTime() {
				return
			}

			if position, err := w.generatePosition(); err != nil {
				w.SendToClient(newErrMessageProto(errGetPosition))
			} else {
				w.SendToClient(newFuturePositionProto(position))
			}
		}
	}
}

func (w *WSFutureTrade) generateTradeIndex() *entity.TradeIndex {
	return w.s.GetTradeIndex()
}

func (w *WSFutureTrade) generatePosition() (entity.FuturePositionArr, error) {
	position, err := w.o.GetFuturePosition()
	if err != nil {
		return nil, err
	}
	return position, nil
}

func (w *WSFutureTrade) sendLatestKbar() {
	go func() {
		if data := w.fetchKbar(); data != nil {
			w.SendToClient(newKbarArrProto(data))
		}

		for {
			select {
			case <-w.Ctx().Done():
				return

			case <-time.After(time.Minute):
				if data := w.fetchKbar(); data != nil {
					w.SendToClient(newKbarArrProto(data))
				}
			}
		}
	}()
}

func (w *WSFutureTrade) fetchKbar() []*entity.FutureHistoryKbar {
	var kbarArr []*entity.FutureHistoryKbar
	var err error
	firstTry := time.Now()

	for {
		if kbarArr, err = w.h.FetchFutureHistoryKbar(w.s.GetMainFuture().Code, firstTry); err != nil {
			w.SendToClient(newErrMessageProto(errGetKbarFail))
			return nil
		} else if len(kbarArr) > 0 {
			break
		}
		firstTry = firstTry.Add(-24 * time.Hour)
	}

	var singleArr []*entity.FutureHistoryKbar
	for i, kbar := range kbarArr {
		if i == 0 {
			singleArr = append(singleArr, kbar)
			continue
		}

		if kbar.KbarTime.Sub(kbarArr[i-1].KbarTime) > time.Minute {
			singleArr = []*entity.FutureHistoryKbar{}
		}
		singleArr = append(singleArr, kbar)
	}
	return singleArr
}
