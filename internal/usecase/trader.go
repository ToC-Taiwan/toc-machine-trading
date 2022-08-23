package usecase

import (
	"sort"
	"sync"
	"time"
	"tmt/internal/entity"
	"tmt/pkg/config"
	"tmt/pkg/utils"

	"github.com/google/uuid"
)

// TradeAgent -.
type TradeAgent struct {
	stockNum      string
	orderQuantity int64
	tickArr       RealTimeTickArr
	periodTickArr RealTimeTickArr

	orderMapLock sync.RWMutex
	orderMap     map[entity.OrderAction][]*entity.Order
	waitingOrder *entity.Order

	tickChan   chan *entity.RealTimeTick
	bidAskChan chan *entity.RealTimeBidAsk

	historyTickAnalyze []int64
	analyzeTickTime    time.Time
	lastTick           *entity.RealTimeTick
	lastBidAsk         *entity.RealTimeBidAsk

	tradeInWaitTime  time.Duration
	tradeOutWaitTime time.Duration
	cancelWaitTime   time.Duration

	openPass bool
}

// NewAgent -.
func NewAgent(stockNum string, tradeSwitch config.TradeSwitch) *TradeAgent {
	var quantity int64 = 1
	if biasRate := cc.GetBiasRate(stockNum); biasRate > 4 || biasRate < -4 {
		quantity = 2
	} else if biasRate == 0 {
		time.Sleep(time.Second)
		return NewAgent(stockNum, tradeSwitch)
	}

	arr := cc.GetHistoryTickAnalyze(stockNum)
	sort.Slice(arr, func(i, j int) bool {
		return arr[i] > arr[j]
	})

	new := &TradeAgent{
		stockNum:           stockNum,
		orderQuantity:      quantity,
		orderMap:           make(map[entity.OrderAction][]*entity.Order),
		tickChan:           make(chan *entity.RealTimeTick),
		bidAskChan:         make(chan *entity.RealTimeBidAsk),
		historyTickAnalyze: arr,
		tradeInWaitTime:    time.Duration(tradeSwitch.TradeInWaitTime) * time.Second,
		tradeOutWaitTime:   time.Duration(tradeSwitch.TradeOutWaitTime) * time.Second,
		cancelWaitTime:     time.Duration(tradeSwitch.CancelWaitTime) * time.Second,
	}

	go new.checkFirstTickArrive()

	return new
}

func (o *TradeAgent) generateOrder(cfg config.Analyze) *entity.Order {
	if o.lastTick.TickTime.Sub(o.analyzeTickTime) < time.Duration(cfg.TickAnalyzePeriod)*time.Millisecond {
		o.periodTickArr = append(o.periodTickArr, o.lastTick)
		return nil
	}
	// copy new arr before reset
	analyzeArr := o.periodTickArr
	// reset analyze tick time and arr
	o.analyzeTickTime = o.lastTick.TickTime
	o.periodTickArr = RealTimeTickArr{o.lastTick}

	if postOrderAction, preOrder := o.checkNeededPost(); postOrderAction != entity.ActionNone {
		return o.generateTradeOutOrder(cfg, postOrderAction, preOrder)
	}

	if o.lastTick.PctChg < cfg.CloseChangeRatioLow || o.lastTick.PctChg > cfg.CloseChangeRatioHigh {
		return nil
	}

	periodVolume := analyzeArr.getTotalVolume()
	if pr := o.getPRByVolume(periodVolume); pr < cfg.VolumePRLimit {
		return nil
	}

	if o.lastBidAsk == nil {
		return nil
	}

	// get out in ration in this period
	periodOutInRation := analyzeArr.getOutInRatio()
	allOutInRation := o.tickArr.getOutInRatio()

	// need to compare with all and period
	order := &entity.Order{
		StockNum: o.stockNum,
		Quantity: o.orderQuantity,
		TickTime: o.lastTick.TickTime,
		GroupID:  uuid.New().String(),
	}

	switch {
	case periodOutInRation-allOutInRation > cfg.AllOutInRatio*0.1 && allOutInRation > cfg.AllOutInRatio:
		order.Action = entity.ActionBuy
		order.Price = o.lastBidAsk.AskPrice1
	case allOutInRation-periodOutInRation > cfg.AllInOutRatio*0.1 && 100-allOutInRation > cfg.AllInOutRatio:
		order.Action = entity.ActionSellFirst
		order.Price = o.lastBidAsk.BidPrice1
	default:
		return nil
	}

	return order
}

func (o *TradeAgent) generateTradeOutOrder(cfg config.Analyze, postOrderAction entity.OrderAction, preOrder *entity.Order) *entity.Order {
	if o.lastTick.TickTime.After(preOrder.TradeTime.Add(time.Duration(cfg.MaxHoldTime) * time.Minute)) {
		return &entity.Order{
			StockNum:  o.stockNum,
			Action:    postOrderAction,
			Price:     o.lastTick.Close,
			Quantity:  preOrder.Quantity,
			TradeTime: o.lastTick.TickTime,
			TickTime:  o.lastTick.TickTime,
			GroupID:   preOrder.GroupID,
		}
	}

	rsi := o.tickArr.getRSIByTickTime(preOrder.TickTime, cfg.RSIMinCount)
	if rsi == 0 {
		return nil
	}

	if rsi <= 49 || rsi >= 51 {
		return &entity.Order{
			StockNum:  o.stockNum,
			Action:    postOrderAction,
			Price:     o.lastTick.Close,
			Quantity:  preOrder.Quantity,
			TradeTime: o.lastTick.TickTime,
			TickTime:  o.lastTick.TickTime,
			GroupID:   preOrder.GroupID,
		}
	}

	return nil
}

func (o *TradeAgent) checkPlaceOrderStatus(order *entity.Order) {
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
			log.Warnf("Order Filled -> Stock: %s, Action: %d, Price: %.2f, Qty: %d", order.StockNum, order.Action, order.Price, order.Quantity)
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

func (o *TradeAgent) cancelOrder(order *entity.Order) {
	order.TradeTime = time.Time{}
	bus.PublishTopicEvent(topicCancelOrder, order)

	go func() {
		for {
			time.Sleep(time.Second)
			if order.TradeTime.IsZero() {
				continue
			}

			if order.Status == entity.StatusCancelled {
				log.Warnf("Order Canceled -> Stock: %s, Action: %d, Price: %.2f, Qty: %d", order.StockNum, order.Action, order.Price, order.Quantity)
				// if order.Action == entity.ActionBuy || order.Action == entity.ActionSellFirst {
				// 	bus.PublishTopicEvent(topicUnSubscribeTickTargets, order.StockNum)
				// 	return
				// }
				o.waitingOrder = nil
				return
			} else if order.TradeTime.Add(o.cancelWaitTime).Before(time.Now()) {
				log.Warnf("Try Cancel Order Again -> Stock: %s, Action: %d, Price: %.2f, Qty: %d", order.StockNum, order.Action, order.Price, order.Quantity)
				go o.checkPlaceOrderStatus(order)
				return
			}
		}
	}()
}

func (o *TradeAgent) checkNeededPost() (entity.OrderAction, *entity.Order) {
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

func (o *TradeAgent) checkFirstTickArrive() {
	// calculate open change rate here
	//
	basic := cc.GetBasicInfo()
	tradeDay := basic.TradeDay
	lastClose := cc.GetHistoryClose(o.stockNum, basic.LastTradeDay)
	for {
		if len(o.tickArr) > 1 {
			firstTick := o.tickArr[0]
			cc.SetHistoryOpen(o.stockNum, tradeDay, firstTick.Open)

			if firstTick.Open != lastClose {
				bus.PublishTopicEvent(topicUnSubscribeTickTargets, o.stockNum)
				break
			}

			o.analyzeTickTime = o.tickArr[1].TickTime
			o.openPass = true
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (o *TradeAgent) getPRByVolume(volume int64) float64 {
	if len(o.historyTickAnalyze) < 2 {
		return 0
	}
	total := len(o.historyTickAnalyze)

	var position int
	for i, v := range o.historyTickAnalyze {
		if volume >= v {
			position = i
			break
		}
		if i == total-1 && position == 0 {
			position = total
		}
	}
	return 100 * float64(total-position) / float64(total)
}

// func (o *TradeAgent) alreadyTrade() bool {
// 	defer o.orderMapLock.RUnlock()
// 	o.orderMapLock.RLock()
// 	return len(o.orderMap) != 0 && len(o.orderMap)%2 == 0
// }

// RealTimeTickArr -.
type RealTimeTickArr []*entity.RealTimeTick

func (c RealTimeTickArr) getTotalVolume() int64 {
	var volume int64
	for _, v := range c {
		volume += v.Volume
	}
	return volume
}

func (c RealTimeTickArr) getOutInRatio() float64 {
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
		}
	}

	return 100 * float64(outVolume) / float64(outVolume+inVolume)
}

func (c RealTimeTickArr) getRSIByTickTime(preTime time.Time, count int) float64 {
	if len(c) == 0 || preTime.IsZero() {
		return 0
	}

	var tmp []float64
	for _, v := range c {
		if v.TickTime.Equal(preTime) || v.TickTime.After(preTime) {
			tmp = append(tmp, v.Close)
		}
	}

	return utils.GenerateRSI(tmp, count)
}
