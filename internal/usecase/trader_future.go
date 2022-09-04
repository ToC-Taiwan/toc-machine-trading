package usecase

import (
	"sync"
	"time"

	"tmt/internal/entity"
	"tmt/pkg/config"
	"tmt/pkg/utils"

	"github.com/google/uuid"
)

// FutureTradeAgent -.
type FutureTradeAgent struct {
	code          string
	orderQuantity int64

	tickArr       RealTimeFutureTickArr
	periodTickArr RealTimeFutureTickArr

	orderMapLock sync.RWMutex
	orderMap     map[entity.OrderAction][]*entity.FutureOrder

	waitingOrder *entity.FutureOrder

	tickChan chan *entity.RealTimeFutureTick
	lastTick *entity.RealTimeFutureTick

	analyzeTickTime time.Time

	tradeInWaitTime  time.Duration
	tradeOutWaitTime time.Duration
	cancelWaitTime   time.Duration
}

// NewFutureAgent -.
func NewFutureAgent(code string, tradeSwitch config.FutureTradeSwitch) *FutureTradeAgent {
	new := &FutureTradeAgent{
		code:             code,
		orderQuantity:    tradeSwitch.Quantity,
		orderMap:         make(map[entity.OrderAction][]*entity.FutureOrder),
		tickChan:         make(chan *entity.RealTimeFutureTick),
		tradeInWaitTime:  time.Duration(tradeSwitch.TradeInWaitTime) * time.Second,
		tradeOutWaitTime: time.Duration(tradeSwitch.TradeOutWaitTime) * time.Second,
		cancelWaitTime:   time.Duration(tradeSwitch.CancelWaitTime) * time.Second,
	}
	return new
}

func (o *FutureTradeAgent) generateOrder(cfg config.Analyze) *entity.FutureOrder {
	if o.lastTick.TickTime.Sub(o.analyzeTickTime) > time.Duration(cfg.TickAnalyzePeriod*1.1)*time.Millisecond {
		o.analyzeTickTime = o.lastTick.TickTime
		o.periodTickArr = RealTimeFutureTickArr{o.lastTick}
		return nil
	}

	if o.lastTick.TickTime.Sub(o.analyzeTickTime) < time.Duration(cfg.TickAnalyzePeriod)*time.Millisecond {
		o.periodTickArr = append(o.periodTickArr, o.lastTick)
		return nil
	}
	// copy new arr before reset
	// analyzeArr := o.periodTickArr

	// reset analyze tick time and arr
	o.analyzeTickTime = o.lastTick.TickTime
	o.periodTickArr = RealTimeFutureTickArr{o.lastTick}

	if postOrderAction, preOrder := o.checkNeededPost(); postOrderAction != entity.ActionNone {
		return o.generateTradeOutOrder(cfg, postOrderAction, preOrder)
	}

	if o.lastTick.PctChg < cfg.CloseChangeRatioLow || o.lastTick.PctChg > cfg.CloseChangeRatioHigh {
		return nil
	}

	// get out in ration in all tick
	allOutInRation := o.tickArr.getOutInRatio()

	// need to compare with all and period
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
	case allOutInRation > cfg.AllOutInRatio && o.lastTick.Low < o.lastTick.Close:
		order.Action = entity.ActionBuy
		return order
	case 100-allOutInRation > cfg.AllInOutRatio && o.lastTick.High > o.lastTick.Close:
		order.Action = entity.ActionSellFirst
		return order
	default:
		return nil
	}
}

func (o *FutureTradeAgent) generateTradeOutOrder(cfg config.Analyze, postOrderAction entity.OrderAction, preOrder *entity.FutureOrder) *entity.FutureOrder {
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

	if o.lastTick.TickTime.After(preOrder.TradeTime.Add(time.Duration(cfg.MaxHoldTime) * time.Minute)) {
		return order
	}

	rsi := o.tickArr.getRSIByTickTime(preOrder.TickTime, cfg.RSIMinCount)
	if rsi == 0 {
		return nil
	}

	if rsi <= 49 || rsi >= 51 {
		return order
	}

	return nil
}

func (o *FutureTradeAgent) checkPlaceOrderStatus(order *entity.FutureOrder) {
	var timeout time.Duration
	switch order.Action {
	case entity.ActionBuy, entity.ActionSellFirst:
		timeout = o.tradeInWaitTime
	case entity.ActionSell, entity.ActionBuyLater:
		timeout = o.tradeOutWaitTime
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
	bus.PublishTopicEvent(topicCancelFutureOrder, order)

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
			} else if order.TradeTime.Add(o.cancelWaitTime).Before(time.Now()) {
				log.Warnf("Try Cancel Future Order Again -> Future: %s, Action: %d, Price: %.2f, Qty: %d", order.Code, order.Action, order.Price, order.Quantity)
				go o.checkPlaceOrderStatus(order)
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

func (c RealTimeFutureTickArr) getRSIByTickTime(preTime time.Time, count int) float64 {
	if len(c) == 0 || preTime.IsZero() {
		return 0
	}

	var tmp []float64
	for _, v := range c {
		if v.TickTime.Equal(preTime) || v.TickTime.After(preTime) {
			tmp = append(tmp, v.Close)
		}
	}

	return utils.GenerateFutureRSI(tmp, count)
}
