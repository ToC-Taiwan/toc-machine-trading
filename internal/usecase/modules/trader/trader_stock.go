// Package trader package trader
package trader

import (
	"sort"
	"sync"
	"time"

	"tmt/cmd/config"
	"tmt/internal/entity"
	"tmt/internal/usecase/events"

	"github.com/google/uuid"
)

// TradeAgent -.
type TradeAgent struct {
	stockNum      string
	orderQuantity int64
	tickArr       RealTimeTickArr
	periodTickArr RealTimeTickArr

	orderMapLock sync.RWMutex
	orderMap     map[entity.OrderAction][]*entity.StockOrder
	waitingOrder *entity.StockOrder

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
	if biasRate := cc.GetBiasRate(stockNum); biasRate > cc.GetHighBiasRate() || biasRate < cc.GetLowBiasRate() {
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
		orderMap:           make(map[entity.OrderAction][]*entity.StockOrder),
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

// GenerateOrder -.
func (o *TradeAgent) GenerateOrder(cfg config.StockAnalyze) *entity.StockOrder {
	if o.lastTick.TickTime.Sub(o.analyzeTickTime) > time.Duration(cfg.TickAnalyzePeriod*1.1)*time.Millisecond {
		o.analyzeTickTime = o.lastTick.TickTime
		o.periodTickArr = RealTimeTickArr{o.lastTick}
		return nil
	}

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

	if pr := o.getPRByVolume(analyzeArr.getTotalVolume()); pr < cfg.VolumePRLimit {
		return nil
	}

	// get out in ration in all tick
	allOutInRation := o.tickArr.getOutInRatio()

	// need to compare with all and period
	order := &entity.StockOrder{
		StockNum: o.stockNum,
		BaseOrder: entity.BaseOrder{
			Quantity: o.orderQuantity,
			TickTime: o.lastTick.TickTime,
			GroupID:  uuid.New().String(),
		},
	}

	switch {
	case allOutInRation > cfg.AllOutInRatio && o.lastTick.Low < o.lastTick.Close:
		order.Action = entity.ActionBuy
		order.Price = o.lastBidAsk.AskPrice1
		return order
	case 100-allOutInRation > cfg.AllInOutRatio && o.lastTick.High > o.lastTick.Close:
		order.Action = entity.ActionSellFirst
		order.Price = o.lastBidAsk.BidPrice1
		return order
	default:
		return nil
	}
}

func (o *TradeAgent) generateTradeOutOrder(cfg config.StockAnalyze, postOrderAction entity.OrderAction, preOrder *entity.StockOrder) *entity.StockOrder {
	order := &entity.StockOrder{
		StockNum: o.stockNum,
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
		switch postOrderAction {
		case entity.ActionSell:
			order.Price = o.lastBidAsk.AskPrice1
		case entity.ActionBuyLater:
			order.Price = o.lastBidAsk.BidPrice1
		}
		return order
	}

	return nil
}

// CheckPlaceOrderStatus -.
func (o *TradeAgent) CheckPlaceOrderStatus(order *entity.StockOrder) {
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

func (o *TradeAgent) cancelOrder(order *entity.StockOrder) {
	order.TradeTime = time.Time{}
	bus.PublishTopicEvent(events.TopicCancelOrder, order)

	go func() {
		for {
			time.Sleep(time.Second)
			if order.TradeTime.IsZero() {
				continue
			}

			if order.Status == entity.StatusCancelled {
				log.Warnf("Order Canceled -> Stock: %s, Action: %d, Price: %.2f, Qty: %d", order.StockNum, order.Action, order.Price, order.Quantity)
				if order.Action == entity.ActionBuy || order.Action == entity.ActionSellFirst {
					bus.PublishTopicEvent(events.TopicUnSubscribeTickTargets, order.StockNum)
					return
				}
				o.waitingOrder = nil
				return
			} else if order.TradeTime.Add(o.cancelWaitTime).Before(time.Now()) {
				log.Warnf("Try Cancel Order Again -> Stock: %s, Action: %d, Price: %.2f, Qty: %d", order.StockNum, order.Action, order.Price, order.Quantity)
				go o.CheckPlaceOrderStatus(order)
				return
			}
		}
	}()
}

func (o *TradeAgent) checkNeededPost() (entity.OrderAction, *entity.StockOrder) {
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
	// lastClose := cc.GetHistoryClose(o.stockNum, basic.LastTradeDay)
	basic := cc.GetBasicInfo()
	tradeDay := basic.TradeDay
	for {
		if len(o.tickArr) > 1 && o.lastBidAsk != nil {
			firstTick := o.tickArr[0]
			cc.SetHistoryOpen(o.stockNum, tradeDay, firstTick.Open)

			// if firstTick.Open != lastClose {
			// 	bus.PublishTopicEvent(topicUnSubscribeTickTargets, o.stockNum)
			// 	log.Warnf("Not open from last close, unsubscribe %s", o.stockNum)
			// 	break
			// }

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
			position = total - 1
		}
	}
	return 100 * float64(total-position) / float64(total)
}

// func (o *TradeAgent) alreadyTrade() bool {
// 	defer o.orderMapLock.RUnlock()
// 	o.orderMapLock.RLock()
// 	return len(o.orderMap) != 0 && len(o.orderMap)%2 == 0
// }

// GetStockNum -/
func (o *TradeAgent) GetStockNum() string {
	return o.stockNum
}

// GetTickChan -.
func (o *TradeAgent) GetTickChan() chan *entity.RealTimeTick {
	return o.tickChan
}

// GetBidAskChan -.
func (o *TradeAgent) GetBidAskChan() chan *entity.RealTimeBidAsk {
	return o.bidAskChan
}

// ReceiveTick -.
func (o *TradeAgent) ReceiveTick(input *entity.RealTimeTick) {
	o.lastTick = input
	o.tickArr = append(o.tickArr, input)
}

// ReceiveBidAsk -.
func (o *TradeAgent) ReceiveBidAsk(input *entity.RealTimeBidAsk) {
	o.lastBidAsk = input
}

// IsReady -.
func (o *TradeAgent) IsReady() bool {
	if o.waitingOrder != nil || o.analyzeTickTime.IsZero() || !o.openPass {
		return false
	}
	return true
}

// WaitingOrder -.
func (o *TradeAgent) WaitingOrder(order *entity.StockOrder) {
	o.waitingOrder = order
}

// CancelWaitingOrder -.
func (o *TradeAgent) CancelWaitingOrder() {
	o.waitingOrder = nil
}
