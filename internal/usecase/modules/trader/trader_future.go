// Package trader package trader
package trader

import (
	"sync"
	"time"

	"tmt/cmd/config"
	"tmt/internal/entity"
	"tmt/internal/usecase/modules/event"

	"github.com/google/uuid"
)

// FutureTrader -.
type FutureTrader struct {
	code          string
	orderQuantity int64

	orderMapLock       sync.RWMutex
	orderMap           map[entity.OrderAction][]*entity.FutureOrder
	lastPlaceOrderTime time.Time
	tradeOutPrice      float64

	waitingOrder *entity.FutureOrder

	tickArr  realTimeFutureTickArr
	kbarArr  kbarArr
	lastTick *entity.RealTimeFutureTick
	// lastBidAsk *entity.FutureRealTimeBidAsk

	tickChan chan *entity.RealTimeFutureTick
	// bidAskChan chan *entity.FutureRealTimeBidAsk

	tradeSwitch config.FutureTradeSwitch
	analyzeCfg  config.FutureAnalyze

	futureTradeInSwitch bool

	tradeOutRecord map[string]int

	tradeRate  float64
	checkPoint time.Time
}

// NewFutureTrader -.
func NewFutureTrader(code string, tradeSwitch config.FutureTradeSwitch, analyzeCfg config.FutureAnalyze) *FutureTrader {
	t := &FutureTrader{
		code:          code,
		orderQuantity: tradeSwitch.Quantity,
		orderMap:      make(map[entity.OrderAction][]*entity.FutureOrder),
		tickChan:      make(chan *entity.RealTimeFutureTick),
		// bidAskChan:          make(chan *entity.FutureRealTimeBidAsk),
		tradeSwitch:         tradeSwitch,
		analyzeCfg:          analyzeCfg,
		futureTradeInSwitch: false,
		tradeOutRecord:      make(map[string]int),
	}

	bus.SubscribeTopic(event.TopicUpdateFutureTradeSwitch, t.updateAllowTrade)
	return t
}

func (o *FutureTrader) updateAllowTrade(allow bool) {
	o.futureTradeInSwitch = allow
}

// GetTickChan -.
func (o *FutureTrader) GetTickChan() chan *entity.RealTimeFutureTick {
	return o.tickChan
}

// GetBidAskChan -.
// func (o *FutureTrader) GetBidAskChan() chan *entity.FutureRealTimeBidAsk {
// 	return o.bidAskChan
// }

// TradingRoom -.
func (o *FutureTrader) TradingRoom() {
	// go func() {
	// 	for {
	// 		o.lastBidAsk = <-o.bidAskChan
	// 	}
	// }()

	for {
		tick := <-o.tickChan
		o.tickArr = append(o.tickArr, tick)
		if len(o.tickArr) == 2 {
			o.tickArr = o.tickArr[1:]
			o.lastTick = tick
			o.checkPoint = time.Date(
				tick.TickTime.Year(),
				tick.TickTime.Month(),
				tick.TickTime.Day(),
				tick.TickTime.Hour(),
				tick.TickTime.Minute(),
				tick.TickTime.Second(),
				0,
				tick.TickTime.Location(),
			)
			break
		}
	}

	for {
		tick := <-o.tickChan
		if tick.TickTime.Minute() != o.lastTick.TickTime.Minute() {
			o.placeFutureOrder(o.generateOrder())
		}

		o.lastTick = tick
		o.tickArr = append(o.tickArr, tick)
		o.tickArr = o.tickArr[o.tickArr.appendKbar(&o.kbarArr):]
		if o.checkPoint.Before(tick.TickTime) {
			o.checkPoint = o.checkPoint.Add(time.Second)
			o.placeFutureOrder(o.generateOrder())
		}
	}
}

func (o *FutureTrader) generateOrder() *entity.FutureOrder {
	if o.waitingOrder != nil {
		return nil
	}

	if postOrderAction, preOrder := o.checkNeededPost(); postOrderAction != entity.ActionNone {
		return o.generateTradeOutOrder(postOrderAction, preOrder)
	}

	act := entity.ActionNone
	ratio, rate := o.getTradeRate()
	if o.kbarArr.isStable(5, 3) && o.lastTick.TickTime.After(o.lastPlaceOrderTime.Add(time.Minute)) {
		if rate/o.tradeRate > 5 {
			switch {
			case ratio > 55:
				act = entity.ActionBuy
				o.tradeOutPrice = o.lastTick.Close
				o.lastPlaceOrderTime = o.lastTick.TickTime
			case ratio < 45:
				act = entity.ActionSellFirst
				o.tradeOutPrice = o.lastTick.Close
				o.lastPlaceOrderTime = o.lastTick.TickTime
			}
		}
	}
	o.tradeRate = rate

	if act == entity.ActionNone {
		return nil
	}

	// get out in ration in period
	return &entity.FutureOrder{
		Code: o.code,
		BaseOrder: entity.BaseOrder{
			Action:   act,
			Quantity: o.orderQuantity,
			TickTime: o.lastTick.TickTime,
			GroupID:  uuid.New().String(),
			Price:    o.lastTick.Close,
		},
	}
}

func (o *FutureTrader) getTradeRate() (float64, float64) {
	var tmp realTimeFutureTickArr
	for i := len(o.tickArr) - 1; i >= 0; i-- {
		if o.tickArr[i].TickTime.Before(o.checkPoint.Add(-10 * time.Second)) {
			break
		}
		tmp = append(tmp, o.tickArr[i])
	}
	return tmp.getOutInRatio(), float64(tmp.getTotalVolume()) / 10
}

func (o *FutureTrader) generateTradeOutOrder(postOrderAction entity.OrderAction, preOrder *entity.FutureOrder) *entity.FutureOrder {
	order := &entity.FutureOrder{
		Code: o.code,
		BaseOrder: entity.BaseOrder{
			Action:    postOrderAction,
			Price:     o.lastTick.Close,
			Quantity:  preOrder.Quantity,
			TradeTime: o.lastTick.TickTime,
			TickTime:  o.lastTick.TickTime,
			GroupID:   preOrder.GroupID,
		},
	}

	if o.lastTick.TickTime.After(preOrder.TickTime.Add(time.Duration(o.analyzeCfg.MaxHoldTime) * time.Minute)) {
		return order
	}

	switch postOrderAction {
	case entity.ActionSell:
		if order.Price > o.tradeOutPrice {
			return nil
		}
		o.tradeOutPrice = order.Price

		if order.Price > preOrder.Price+2 || order.Price < preOrder.Price-1 {
			return order
		}

	case entity.ActionBuyLater:
		if order.Price < o.tradeOutPrice {
			return nil
		}
		o.tradeOutPrice = order.Price

		if order.Price < preOrder.Price-2 || order.Price > preOrder.Price+1 {
			return order
		}
	}

	return nil
}

func (o *FutureTrader) placeFutureOrder(order *entity.FutureOrder) {
	if order == nil {
		return
	}

	if order.Price == 0 {
		logger.Errorf("%s Future Order price is 0", order.Code)
		return
	}

	// if out of trade in time, return
	if !o.futureTradeInSwitch && (order.Action == entity.ActionBuy || order.Action == entity.ActionSellFirst) {
		return
	}

	o.waitingOrder = order
	go o.checkPlaceOrderStatus(order)
	bus.PublishTopicEvent(event.TopicPlaceFutureOrder, order)
}

func (o *FutureTrader) checkPlaceOrderStatus(order *entity.FutureOrder) {
	var timeout time.Duration
	switch order.Action {
	case entity.ActionBuy, entity.ActionSellFirst:
		timeout = time.Duration(o.tradeSwitch.TradeInWaitTime) * time.Second
	case entity.ActionSell, entity.ActionBuyLater:
		timeout = time.Duration(o.tradeSwitch.TradeOutWaitTime) * time.Second
	}

	for {
		time.Sleep(time.Second)
		if order.TradeTime.IsZero() {
			continue
		}

		if order.Status == entity.StatusFilled {
			o.orderMapLock.Lock()
			o.orderMap[order.Action] = append(o.orderMap[order.Action], order)
			o.orderMapLock.Unlock()

			o.waitingOrder = nil
			logger.Infof("Future Order Filled -> Future: %s, Action: %d, Price: %.2f, Qty: %d", order.Code, order.Action, order.Price, order.Quantity)
			return
		} else if order.TradeTime.Add(timeout).Before(time.Now()) {
			break
		}
	}

	if order.Status == entity.StatusAborted || order.Status == entity.StatusFailed {
		o.waitingOrder = nil
		return
	}

	if order.OrderID != "" && order.Status != entity.StatusCancelled && order.Status != entity.StatusFilled {
		o.cancelOrder(order)
		return
	}

	logger.Error("check place order status raise unknown error")
}

func (o *FutureTrader) cancelOrder(order *entity.FutureOrder) {
	order.TradeTime = time.Time{}
	bus.PublishTopicEvent(event.TopicCancelFutureOrder, order)

	go func() {
		for {
			time.Sleep(time.Second)
			if order.TradeTime.IsZero() {
				continue
			}

			if order.Status == entity.StatusCancelled {
				logger.Infof("Future Order Canceled -> Future: %s, Action: %d, Price: %.2f, Qty: %d", order.Code, order.Action, order.Price, order.Quantity)
				o.waitingOrder = nil
				return
			} else if order.TradeTime.Add(time.Duration(o.tradeSwitch.CancelWaitTime) * time.Second).Before(time.Now()) {
				logger.Warnf("Try Cancel Future Order Again -> Future: %s, Action: %d, Price: %.2f, Qty: %d", order.Code, order.Action, order.Price, order.Quantity)
				go o.checkPlaceOrderStatus(order)
				return
			}
		}
	}()
}

func (o *FutureTrader) checkNeededPost() (entity.OrderAction, *entity.FutureOrder) {
	defer o.orderMapLock.RUnlock()
	o.orderMapLock.RLock()

	if len(o.orderMap[entity.ActionBuy]) > len(o.orderMap[entity.ActionSell]) {
		return entity.ActionSell, o.orderMap[entity.ActionBuy][len(o.orderMap[entity.ActionSell])]
	}

	if len(o.orderMap[entity.ActionSellFirst]) > len(o.orderMap[entity.ActionBuyLater]) {
		return entity.ActionBuyLater, o.orderMap[entity.ActionSellFirst][len(o.orderMap[entity.ActionBuyLater])]
	}

	return entity.ActionNone, nil
}
