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

	analyzePeriod int

	lastArr       realTimeFutureTickArr
	magnification float64
	avgVolume     float64

	orderMapLock sync.RWMutex
	orderMap     map[entity.OrderAction][]*entity.FutureOrder

	waitingOrder *entity.FutureOrder

	tickArr    realTimeFutureTickArr
	lastTick   *entity.RealTimeFutureTick
	lastBidAsk *entity.FutureRealTimeBidAsk

	tickChan   chan *entity.RealTimeFutureTick
	bidAskChan chan *entity.FutureRealTimeBidAsk

	tradeSwitch config.FutureTradeSwitch
	analyzeCfg  config.FutureAnalyze

	futureTradeInSwitch bool
}

// NewFutureTrader -.
func NewFutureTrader(code string, tradeSwitch config.FutureTradeSwitch, analyzeCfg config.FutureAnalyze) *FutureTrader {
	t := &FutureTrader{
		code:                code,
		orderQuantity:       tradeSwitch.Quantity,
		orderMap:            make(map[entity.OrderAction][]*entity.FutureOrder),
		tickChan:            make(chan *entity.RealTimeFutureTick),
		bidAskChan:          make(chan *entity.FutureRealTimeBidAsk),
		tradeSwitch:         tradeSwitch,
		analyzeCfg:          analyzeCfg,
		futureTradeInSwitch: false,
		analyzePeriod:       int(analyzeCfg.TickAnalyzePeriod),
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
func (o *FutureTrader) GetBidAskChan() chan *entity.FutureRealTimeBidAsk {
	return o.bidAskChan
}

// TradingRoom -.
func (o *FutureTrader) TradingRoom() {
	go func() {
		for {
			o.lastBidAsk = <-o.bidAskChan
		}
	}()

	for {
		tick := <-o.tickChan
		for i, v := range o.tickArr {
			if time.Since(v.TickTime) < time.Duration(o.analyzePeriod)*time.Millisecond {
				o.tickArr = o.tickArr[i:]
				break
			}
		}
		o.tickArr = append(o.tickArr, tick)
		tmp := o.tickArr.splitBySecond()
		var total, last int
		for i, v := range tmp {
			if i != len(tmp)-1 {
				total += int(v.getTotalVolume())
			} else {
				last = int(v.getTotalVolume())
				o.lastArr = v
			}
		}
		if avg := float64(total) / float64(len(tmp)-2); avg != 0 {
			o.magnification = float64(last) / avg
			o.avgVolume = avg
		}

		o.lastTick = tick
		o.placeFutureOrder(o.generateOrder())
	}
}

func (o *FutureTrader) placeFutureOrder(order *entity.FutureOrder) {
	if order == nil {
		return
	}

	if order.Price == 0 {
		log.Errorf("%s Future Order price is 0", order.Code)
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

func (o *FutureTrader) generateOrder() *entity.FutureOrder {
	if o.waitingOrder != nil {
		return nil
	}

	if postOrderAction, preOrder := o.checkNeededPost(); postOrderAction != entity.ActionNone {
		return o.generateTradeOutOrder(postOrderAction, preOrder)
	}

	if o.magnification <= 3 {
		return nil
	}

	// // get out in ration in period
	outInRation := o.lastArr.getOutInRatio()
	order := &entity.FutureOrder{
		Code: o.code,
		BaseOrder: entity.BaseOrder{
			Quantity: o.orderQuantity,
			TickTime: o.lastTick.TickTime,
			GroupID:  uuid.New().String(),
			Price:    o.lastTick.Close,
		},
	}

	switch {
	case outInRation > o.analyzeCfg.AllOutInRatio:
		order.Action = entity.ActionBuy
		return order
	case 100-outInRation > o.analyzeCfg.AllInOutRatio:
		order.Action = entity.ActionSellFirst
		return order
	default:
		return nil
	}
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

	if o.lastTick.TickTime.After(preOrder.TradeTime.Add(time.Duration(o.analyzeCfg.MaxHoldTime) * time.Minute)) {
		return order
	}

	outInRation := o.tickArr.getOutInRatio()
	switch order.Action {
	case entity.ActionSell:
		if order.Price-preOrder.Price < -2 {
			return order
		}

		if 100-outInRation > o.analyzeCfg.AllInOutRatio && order.Price-preOrder.Price > 2 {
			return order
		}

	case entity.ActionBuyLater:
		if order.Price-preOrder.Price > 2 {
			return order
		}

		if outInRation > o.analyzeCfg.AllInOutRatio && order.Price-preOrder.Price < -2 {
			return order
		}
	}
	return nil
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
			log.Warnf("Future Order Filled -> Future: %s, Action: %d, Price: %.2f, Qty: %d", order.Code, order.Action, order.Price, order.Quantity)
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

	log.Error("check place order status raise unknown error")
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
				log.Warnf("Future Order Canceled -> Future: %s, Action: %d, Price: %.2f, Qty: %d", order.Code, order.Action, order.Price, order.Quantity)
				o.waitingOrder = nil
				return
			} else if order.TradeTime.Add(time.Duration(o.tradeSwitch.CancelWaitTime) * time.Second).Before(time.Now()) {
				log.Warnf("Try Cancel Future Order Again -> Future: %s, Action: %d, Price: %.2f, Qty: %d", order.Code, order.Action, order.Price, order.Quantity)
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
