package future

// type assistTrader struct {
// 	assistOrderMap map[string]*entity.FutureOrder
// 	assistOption   map[string]HalfAutomationOption
// 	mutex          sync.Mutex

// 	msgChan  chan interface{}
// 	tickChan chan *entity.RealTimeFutureTick
// }

// func newAssistTrader(msgChan chan interface{}) *assistTrader {
// 	return &assistTrader{
// 		assistOrderMap: make(map[string]*entity.FutureOrder),
// 		assistOption:   make(map[string]HalfAutomationOption),
// 		msgChan:        msgChan,
// 	}
// }

// func (a *assistTrader) addAssistOrder(order *entity.FutureOrder, option HalfAutomationOption) chan *entity.RealTimeFutureTick {
// 	a.mutex.Lock()
// 	defer a.mutex.Unlock()

// 	a.assistOrderMap[order.OrderID] = order
// 	a.assistOption[order.OrderID] = option
// 	a.tickChan = make(chan *entity.RealTimeFutureTick)
// 	return a.tickChan
// }

// func (w *WSFutureTrade) prepareAssist(orderID string) {
// 	for {
// 		time.Sleep(time.Second)
// 		w.orderLock.Lock()
// 		if order, ok := w.assistOrderMap[orderID]; ok && order.Status == entity.StatusFilled {
// 			tickChan := make(chan *entity.RealTimeFutureTick)
// 			go w.assistTrader(order.OrderID, tickChan)
// 			w.addAssistTrader(order.OrderID, tickChan)
// 			return
// 		}
// 		w.orderLock.Unlock()
// 	}
// }

// func (w *WSFutureTrade) assistTrader(orderID string, tickChan chan *entity.RealTimeFutureTick) {
// 	w.orderLock.Lock()
// 	option := w.assistOption[orderID]
// 	w.orderLock.Unlock()
// 	for {
// 		select {
// 		case <-w.ctx.Done():
// 			return

// 		case tick := <-tickChan:
// 			if option.AutomationType == AutomationByBalance {
// 				log.Warn(tick)
// 			}

// 		case <-time.After(time.Second):
// 			w.checkAssistStatus()
// 		}
// 	}
// }

// func (w *WSFutureTrade) checkAssistStatus() {
// 	var qty int64
// 	w.orderLock.Lock()
// 	for _, order := range w.assistOrderMap {
// 		if order.Status == entity.StatusFilled {
// 			w.msgChan <- assistStatus{true}
// 			switch order.Action {
// 			case entity.ActionBuy:
// 				qty += order.Quantity
// 			case entity.ActionSell:
// 				qty -= order.Quantity
// 			}
// 		}
// 	}

// 	if qty == 0 {
// 		w.assistOrderMap = make(map[string]*entity.FutureOrder)
// 		w.msgChan <- assistStatus{false}
// 	}
// 	w.orderLock.Unlock()
// }
