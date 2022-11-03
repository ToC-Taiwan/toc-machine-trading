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

	orderMapLock sync.RWMutex
	orderMap     map[entity.OrderAction][]*entity.FutureOrder

	waitingOrder *entity.FutureOrder

	tickArr    realTimeFutureTickArr
	kbarArr    realTimeKbarArr
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
		o.lastTick = tick
		o.tickArr = append(o.tickArr, tick)
		break
	}

	for {
		tick := <-o.tickChan
		if tick.TickTime.Minute() != o.lastTick.TickTime.Minute() {
			o.kbarArr = append(o.kbarArr, o.tickArr.getKbar())
			o.tickArr = realTimeFutureTickArr{}
		}

		o.lastTick = tick
		o.tickArr = append(o.tickArr, tick)
		o.placeFutureOrder(o.generateOrder())
	}
}

func (o *FutureTrader) generateOrder() *entity.FutureOrder {
	if o.waitingOrder != nil {
		return nil
	}

	if postOrderAction, preOrder := o.checkNeededPost(); postOrderAction != entity.ActionNone {
		return o.generateTradeOutOrder(postOrderAction, preOrder)
	}

	if !o.kbarArr.isStable(10) {
		return nil
	}

	splitBySecondArr := o.tickArr.splitBySecond()
	if len(splitBySecondArr) < 5 {
		return nil
	}

	base := splitBySecondArr[len(splitBySecondArr)-1].getTotalVolume()
	for i := len(splitBySecondArr) - 2; i >= 0; i-- {
		if splitBySecondArr[i].getTotalVolume() > base {
			return nil
		}
	}

	// get out in ration in period
	outInRation := o.tickArr.getOutInRatio()
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
	case outInRation >= o.analyzeCfg.AllOutInRatio:
		order.Action = entity.ActionBuy
		return order
	case 100-outInRation >= o.analyzeCfg.AllInOutRatio:
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

	// if o.lastTick.TickTime.After(preOrder.TradeTime.Add(time.Duration(o.analyzeCfg.MaxHoldTime) * time.Minute)) {
	// 	return order
	// }

	if !o.kbarArr.isStable(3) {
		return nil
	}

	// outInRation := o.tickArr.getOutInRatio()
	switch order.Action {
	case entity.ActionSell:
		if order.Price-preOrder.Price < -2 || order.Price-preOrder.Price > 1 {
			return order
		}

	case entity.ActionBuyLater:
		if order.Price-preOrder.Price > 2 || order.Price-preOrder.Price < -1 {
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
			log.Infof("Future Order Filled -> Future: %s, Action: %d, Price: %.2f, Qty: %d", order.Code, order.Action, order.Price, order.Quantity)
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
				log.Infof("Future Order Canceled -> Future: %s, Action: %d, Price: %.2f, Qty: %d", order.Code, order.Action, order.Price, order.Quantity)
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
