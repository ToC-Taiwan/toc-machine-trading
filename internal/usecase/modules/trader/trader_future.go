// Package trader package trader
package trader

import (
	"sync"
	"time"

	"tmt/cmd/config"
	"tmt/internal/entity"
	"tmt/internal/usecase/modules/event"
	"tmt/pkg/logger"
	"tmt/pkg/utils"

	"github.com/google/uuid"
)

var log = logger.Get()

// FutureTrader -.
type FutureTrader struct {
	code          string
	orderQuantity int64

	analyzePeriod int

	magnification float64
	outInRation   float64
	avgVolume     float64

	orderMapLock sync.RWMutex
	orderMap     map[entity.OrderAction][]*entity.FutureOrder

	waitingOrder *entity.FutureOrder

	tickArr    RealTimeFutureTickArr
	lastTick   *entity.RealTimeFutureTick
	lastBidAsk *entity.FutureRealTimeBidAsk

	tickChan   chan *entity.RealTimeFutureTick
	bidAskChan chan *entity.FutureRealTimeBidAsk

	tradeSwitch config.FutureTradeSwitch
	analyzeCfg  config.FutureAnalyze
}

// NewFutureTrader -.
func NewFutureTrader(code string, tradeSwitch config.FutureTradeSwitch, analyzeCfg config.FutureAnalyze) *FutureTrader {
	new := &FutureTrader{
		code:          code,
		analyzePeriod: 15,
		orderQuantity: tradeSwitch.Quantity,
		orderMap:      make(map[entity.OrderAction][]*entity.FutureOrder),
		tickChan:      make(chan *entity.RealTimeFutureTick),
		bidAskChan:    make(chan *entity.FutureRealTimeBidAsk),
		tradeSwitch:   tradeSwitch,
		analyzeCfg:    analyzeCfg,
	}
	return new
}

// ReceiveTick -.
func (o *FutureTrader) ReceiveTick(tick *entity.RealTimeFutureTick) *entity.RealTimeFutureTick {
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
	if avg := float64(total) / float64(len(tmp)-2); avg != 0 {
		o.magnification = float64(last) / avg
		o.avgVolume = avg
	}
	return tick
}

// GetLastBidAsk -.
func (o *FutureTrader) GetLastBidAsk() *entity.FutureRealTimeBidAsk {
	if o.lastBidAsk == nil {
		return &entity.FutureRealTimeBidAsk{}
	}
	return o.lastBidAsk
}

// ReceiveBidAsk -.
func (o *FutureTrader) ReceiveBidAsk(bidAsk *entity.FutureRealTimeBidAsk) {
	o.lastBidAsk = bidAsk
}

// GetTickChan -.
func (o *FutureTrader) GetTickChan() chan *entity.RealTimeFutureTick {
	return o.tickChan
}

// GetBidAskChan -.
func (o *FutureTrader) GetBidAskChan() chan *entity.FutureRealTimeBidAsk {
	return o.bidAskChan
}

// WaitOrder -.
func (o *FutureTrader) WaitOrder(order *entity.FutureOrder) {
	o.waitingOrder = order
}

// CancelWaiting -.
func (o *FutureTrader) CancelWaiting() {
	o.waitingOrder = nil
}

// IsWaiting -.
func (o *FutureTrader) IsWaiting() bool {
	return o.waitingOrder != nil
}

// GenerateOrder -.
func (o *FutureTrader) GenerateOrder() *entity.FutureOrder {
	if postOrderAction, preOrder := o.checkNeededPost(); postOrderAction != entity.ActionNone {
		return o.generateTradeOutOrder(postOrderAction, preOrder)
	}

	if o.magnification <= 2 {
		o.outInRation = 0.0
		return nil
	}

	// // get out in ration in period
	outInRation := o.tickArr.getOutInRatio()
	o.outInRation = outInRation
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

// GetAvgVolume -.
func (o *FutureTrader) GetAvgVolume() float64 {
	if o.avgVolume == 0 {
		return 0
	}
	return utils.Round(o.avgVolume, 2)
}

// GetOutInRatio -.
func (o *FutureTrader) GetOutInRatio() float64 {
	return utils.Round(o.outInRation, 2)
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
func (o *FutureTrader) CheckPlaceOrderStatus(order *entity.FutureOrder) {
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
				go o.CheckPlaceOrderStatus(order)
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
