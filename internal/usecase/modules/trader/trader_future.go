// Package trader package trader
package trader

import (
	"sync"
	"time"

	"tmt/cmd/config"
	"tmt/internal/entity"
	"tmt/internal/usecase/events"
	"tmt/pkg/eventbus"
	"tmt/pkg/logger"

	"github.com/google/uuid"
)

var log = logger.Get()

// FutureTradeAgent -.
type FutureTradeAgent struct {
	code          string
	orderQuantity int64
	analyzePeriod int
	magnification float64
	orderMapLock  sync.RWMutex
	bus           *eventbus.Bus
	waitingOrder  *entity.FutureOrder
	tickArr       RealTimeFutureTickArr
	periodMap     map[int64]RealTimeFutureTickArr
	orderMap      map[entity.OrderAction][]*entity.FutureOrder
	tickChan      chan *entity.RealTimeFutureTick
	bidAskChan    chan *entity.FutureRealTimeBidAsk
	lastTick      *entity.RealTimeFutureTick
	lastBidAsk    *entity.FutureRealTimeBidAsk

	tradeSwitch config.FutureTradeSwitch
	analyzeCfg  config.FutureAnalyze
}

// NewFutureAgent -.
func NewFutureAgent(code string, tradeSwitch config.FutureTradeSwitch, analyzeCfg config.FutureAnalyze, bus *eventbus.Bus) *FutureTradeAgent {
	new := &FutureTradeAgent{
		code:          code,
		analyzePeriod: 15,
		orderQuantity: tradeSwitch.Quantity,
		periodMap:     make(map[int64]RealTimeFutureTickArr),
		orderMap:      make(map[entity.OrderAction][]*entity.FutureOrder),
		tickChan:      make(chan *entity.RealTimeFutureTick),
		bidAskChan:    make(chan *entity.FutureRealTimeBidAsk),
		tradeSwitch:   tradeSwitch,
		analyzeCfg:    analyzeCfg,
		bus:           bus,
	}
	return new
}

// ReceiveTick -.
func (o *FutureTradeAgent) ReceiveTick(tick *entity.RealTimeFutureTick) *entity.RealTimeFutureTick {
	o.lastTick = tick
	for i, v := range o.tickArr {
		if time.Since(v.TickTime) < time.Duration(o.analyzePeriod)*time.Second {
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
		}
	}
	avg := float64(total) / float64(len(tmp)-2)
	o.magnification = float64(last) / avg
	return tick
}

// GetLastBidAsk -.
func (o *FutureTradeAgent) GetLastBidAsk() *entity.FutureRealTimeBidAsk {
	return o.lastBidAsk
}

// ReceiveBidAsk -.
func (o *FutureTradeAgent) ReceiveBidAsk(bidAsk *entity.FutureRealTimeBidAsk) {
	o.lastBidAsk = bidAsk
}

// GetTickChan -.
func (o *FutureTradeAgent) GetTickChan() chan *entity.RealTimeFutureTick {
	return o.tickChan
}

// GetBidAskChan -.
func (o *FutureTradeAgent) GetBidAskChan() chan *entity.FutureRealTimeBidAsk {
	return o.bidAskChan
}

// WaitOrder -.
func (o *FutureTradeAgent) WaitOrder(order *entity.FutureOrder) {
	o.waitingOrder = order
}

// CancelWaiting -.
func (o *FutureTradeAgent) CancelWaiting() {
	o.waitingOrder = nil
}

// IsWaiting -.
func (o *FutureTradeAgent) IsWaiting() bool {
	return o.waitingOrder != nil
}

// GenerateOrder -.
func (o *FutureTradeAgent) GenerateOrder() *entity.FutureOrder {
	if postOrderAction, preOrder := o.checkNeededPost(); postOrderAction != entity.ActionNone {
		return o.generateTradeOutOrder(postOrderAction, preOrder)
	}

	if o.magnification < 5 {
		return nil
	}

	// // get out in ration in period
	outInRation := o.tickArr.getOutInRatio()
	log.Warnf("magnification: %.2f, outInRatio: %.2f", o.magnification, outInRation)
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

func (o *FutureTradeAgent) generateTradeOutOrder(postOrderAction entity.OrderAction, preOrder *entity.FutureOrder) *entity.FutureOrder {
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

	switch order.Action {
	case entity.ActionSell:
		if order.Price-preOrder.Price > 3 || order.Price-preOrder.Price < -2 {
			return order
		}
	case entity.ActionBuyLater:
		if order.Price-preOrder.Price < -3 || order.Price-preOrder.Price > 2 {
			return order
		}
	}
	return nil
}

// CheckPlaceOrderStatus -.
func (o *FutureTradeAgent) CheckPlaceOrderStatus(order *entity.FutureOrder) {
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

func (o *FutureTradeAgent) cancelOrder(order *entity.FutureOrder) {
	order.TradeTime = time.Time{}
	o.bus.PublishTopicEvent(events.TopicCancelFutureOrder, order)

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
				go o.CheckPlaceOrderStatus(order)
				return
			}
		}
	}()
}

func (o *FutureTradeAgent) checkNeededPost() (entity.OrderAction, *entity.FutureOrder) {
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

// RealTimeFutureTickArr -.
type RealTimeFutureTickArr []*entity.RealTimeFutureTick

func (c RealTimeFutureTickArr) splitBySecond() []RealTimeFutureTickArr {
	if len(c) < 2 {
		return nil
	}

	var result []RealTimeFutureTickArr
	var tmp RealTimeFutureTickArr
	for i, tick := range c {
		if i == len(c)-1 {
			result = append(result, tmp)
			break
		}

		if tick.TickTime.Second() == c[i+1].TickTime.Second() {
			tmp = append(tmp, tick)
		} else {
			result = append(result, tmp)
			tmp = RealTimeFutureTickArr{tick}
		}
	}

	return result
}

func (c RealTimeFutureTickArr) getTotalVolume() int64 {
	var volume int64
	for _, v := range c {
		volume += v.Volume
	}
	return volume
}

func (c RealTimeFutureTickArr) getOutInRatio() float64 {
	if len(c) == 0 {
		return 0
	}

	var outVolume, inVolume int64
	for _, v := range c {
		switch v.TickType {
		case 1:
			outVolume += v.Volume
		case 2:
			inVolume += v.Volume
		default:
			continue
		}
	}
	return 100 * float64(outVolume) / float64(outVolume+inVolume)
}
