package usecase

import (
	"sync"
	"time"

	"tmt/cmd/config"
	"tmt/internal/entity"
	"tmt/pkg/utils"
)

// FutureTradeAgent -.
type FutureTradeAgent struct {
	code          string
	orderQuantity int64

	tickArr       RealTimeFutureTickArr
	periodTickArr RealTimeFutureTickArr
	periodMap     map[int64]RealTimeFutureTickArr

	orderMapLock sync.RWMutex
	orderMap     map[entity.OrderAction][]*entity.FutureOrder

	waitingOrder *entity.FutureOrder

	tickChan chan *entity.RealTimeFutureTick
	lastTick *entity.RealTimeFutureTick

	analyzeTickTime time.Time
	tradeSwitch     config.FutureTradeSwitch
	analyzeCfg      config.FutureAnalyze
}

// NewFutureAgent -.
func NewFutureAgent(code string, tradeSwitch config.FutureTradeSwitch, analyzeCfg config.FutureAnalyze) *FutureTradeAgent {
	new := &FutureTradeAgent{
		code:          code,
		orderQuantity: tradeSwitch.Quantity,
		periodMap:     make(map[int64]RealTimeFutureTickArr),
		orderMap:      make(map[entity.OrderAction][]*entity.FutureOrder),
		tickChan:      make(chan *entity.RealTimeFutureTick),
		tradeSwitch:   tradeSwitch,
		analyzeCfg:    analyzeCfg,
	}
	return new
}

func (o *FutureTradeAgent) generateOrder() *entity.FutureOrder {
	o.periodTickArr = append(o.periodTickArr, o.lastTick)
	tmp := o.periodTickArr.splitBySecond()
	if len(tmp) < 2 {
		return nil
	}
	lastSecond := tmp[len(tmp)-2]
	if _, ok := o.periodMap[lastSecond.getFirstTickTimestamp()]; !ok {
		o.periodMap[lastSecond.getFirstTickTimestamp()] = lastSecond
		log.Warn(lastSecond.getFirstTickTimestamp(), lastSecond.getTotalVolume())
	}
	return nil

	// TODO: use volume from config
	// if volume := o.periodTickArr.getTotalVolume(); volume < 25 {
	// 	return nil
	// }

	// // get out in ration in period
	// outInRation := o.periodTickArr.getOutInRatio()

	// // reset analyze tick time and arr
	// o.analyzeTickTime = o.lastTick.TickTime
	// o.periodTickArr = RealTimeFutureTickArr{o.lastTick}

	// if postOrderAction, preOrder := o.checkNeededPost(); postOrderAction != entity.ActionNone {
	// 	return o.generateTradeOutOrder(postOrderAction, preOrder)
	// }

	// // need to compare with all and period
	// order := &entity.FutureOrder{
	// 	Code: o.code,
	// 	BaseOrder: entity.BaseOrder{
	// 		Quantity: o.orderQuantity,
	// 		TickTime: o.lastTick.TickTime,
	// 		GroupID:  uuid.New().String(),
	// 		Price:    o.lastTick.Close,
	// 	},
	// }

	// switch {
	// case outInRation > o.analyzeCfg.AllOutInRatio && o.lastTick.Low < o.lastTick.Close:
	// 	order.Action = entity.ActionBuy
	// 	return order
	// case 100-outInRation > o.analyzeCfg.AllInOutRatio && o.lastTick.High > o.lastTick.Close:
	// 	order.Action = entity.ActionSellFirst
	// 	return order
	// default:
	// 	return nil
	// }
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

	rsi := o.tickArr.getRSIByTickTime(preOrder.TickTime, o.analyzeCfg.RSIMinCount)
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
			} else if order.TradeTime.Add(time.Duration(o.tradeSwitch.CancelWaitTime) * time.Second).Before(time.Now()) {
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

func (c RealTimeFutureTickArr) splitBySecond() []RealTimeFutureTickArr {
	if len(c) < 2 {
		return nil
	}

	var result []RealTimeFutureTickArr
	var tmp RealTimeFutureTickArr
	for i, tick := range c {
		if i == len(tmp)-1 {
			break
		}

		if i == 0 {
			tmp = append(tmp, tick)
			continue
		}

		if tick.TickTime.Second() == tmp[i+1].TickTime.Second() {
			tmp = append(tmp, tick)
		} else {
			result = append(result, tmp)
			tmp = RealTimeFutureTickArr{tick}
		}
	}
	return result
}

func (c RealTimeFutureTickArr) getFirstTickTimestamp() int64 {
	return c[0].TickTime.Unix()
}

func (c RealTimeFutureTickArr) getTotalTime() float64 {
	if len(c) < 2 {
		return 0
	}
	firstTickTime := c[0].TickTime
	lastTickTime := c[len(c)-1].TickTime
	return lastTickTime.Sub(firstTickTime).Seconds()
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
