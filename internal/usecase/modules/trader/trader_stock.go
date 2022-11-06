// Package trader package trader
package trader

import (
	"sort"
	"sync"
	"time"

	"tmt/cmd/config"
	"tmt/internal/entity"
	"tmt/internal/usecase/modules/event"

	"github.com/google/uuid"
)

// StockTrader -.
type StockTrader struct {
	stockNum      string
	orderQuantity int64
	tickArr       realTimeStockTickArr
	periodTickArr realTimeStockTickArr

	stockAnalyzeCfg config.StockAnalyze

	orderMapLock sync.RWMutex
	orderMap     map[entity.OrderAction][]*entity.StockOrder
	waitingOrder *entity.StockOrder

	tickChan   chan *entity.RealTimeStockTick
	bidAskChan chan *entity.RealTimeStockBidAsk

	historyTickAnalyze []int64
	analyzeTickTime    time.Time
	lastTick           *entity.RealTimeStockTick
	lastBidAsk         *entity.RealTimeStockBidAsk

	tradeInWaitTime  time.Duration
	tradeOutWaitTime time.Duration
	cancelWaitTime   time.Duration

	openPass           bool
	stockTradeInSwitch bool
}

// NewStockTrader -.
func NewStockTrader(stockNum string, tradeSwitch config.StockTradeSwitch, stockAnalyzeCfg config.StockAnalyze) *StockTrader {
	var quantity int64 = 1
	if biasRate := cc.GetBiasRate(stockNum); biasRate > cc.GetHighBiasRate() || biasRate < cc.GetLowBiasRate() {
		quantity = 2
	}

	arr := cc.GetHistoryTickAnalyze(stockNum)
	sort.Slice(arr, func(i, j int) bool {
		return arr[i] > arr[j]
	})

	new := &StockTrader{
		stockNum:           stockNum,
		orderQuantity:      quantity,
		orderMap:           make(map[entity.OrderAction][]*entity.StockOrder),
		tickChan:           make(chan *entity.RealTimeStockTick),
		bidAskChan:         make(chan *entity.RealTimeStockBidAsk),
		historyTickAnalyze: arr,
		tradeInWaitTime:    time.Duration(tradeSwitch.TradeInWaitTime) * time.Second,
		tradeOutWaitTime:   time.Duration(tradeSwitch.TradeOutWaitTime) * time.Second,
		cancelWaitTime:     time.Duration(tradeSwitch.CancelWaitTime) * time.Second,
		stockAnalyzeCfg:    stockAnalyzeCfg,
	}

	go new.checkFirstTickArrive()

	bus.SubscribeTopic(event.TopicUpdateStockTradeSwitch, new.updateStockTradeSwitch)
	return new
}

func (o *StockTrader) updateStockTradeSwitch(allow bool) {
	o.stockTradeInSwitch = allow
}

// TradingRoom -.
func (o *StockTrader) TradingRoom() {
	go func() {
		for {
			tick := <-o.tickChan
			o.lastTick = tick
			o.tickArr = append(o.tickArr, tick)

			if o.waitingOrder != nil || o.analyzeTickTime.IsZero() || !o.openPass {
				continue
			}

			o.placeOrder(o.generateOrder())
		}
	}()

	go func() {
		for {
			o.lastBidAsk = <-o.bidAskChan
		}
	}()
}

func (o *StockTrader) placeOrder(order *entity.StockOrder) {
	if order == nil {
		return
	}

	if order.Price == 0 {
		log.Errorf("%s Order price is 0", order.StockNum)
		return
	}

	o.waitingOrder = order

	// if out of trade in time, return
	if !o.stockTradeInSwitch && (order.Action == entity.ActionBuy || order.Action == entity.ActionSellFirst) {
		// avoid stuck in the market
		o.waitingOrder = nil
		return
	}

	bus.PublishTopicEvent(event.TopicPlaceStockOrder, order)
	go o.checkPlaceOrderStatus(order)
}

func (o *StockTrader) generateOrder() *entity.StockOrder {
	if o.lastTick.TickTime.Sub(o.analyzeTickTime) > time.Duration(o.stockAnalyzeCfg.TickAnalyzePeriod*1.1)*time.Millisecond {
		o.analyzeTickTime = o.lastTick.TickTime
		o.periodTickArr = realTimeStockTickArr{o.lastTick}
		return nil
	}

	if o.lastTick.TickTime.Sub(o.analyzeTickTime) < time.Duration(o.stockAnalyzeCfg.TickAnalyzePeriod)*time.Millisecond {
		o.periodTickArr = append(o.periodTickArr, o.lastTick)
		return nil
	}
	// copy new arr before reset
	analyzeArr := o.periodTickArr
	// reset analyze tick time and arr
	o.analyzeTickTime = o.lastTick.TickTime
	o.periodTickArr = realTimeStockTickArr{o.lastTick}

	if postOrderAction, preOrder := o.checkNeededPost(); postOrderAction != entity.ActionNone {
		return o.generateTradeOutOrder(o.stockAnalyzeCfg, postOrderAction, preOrder)
	}

	if o.lastTick.PctChg < o.stockAnalyzeCfg.CloseChangeRatioLow || o.lastTick.PctChg > o.stockAnalyzeCfg.CloseChangeRatioHigh {
		return nil
	}

	if pr := o.getPRByVolume(analyzeArr.getTotalVolume()); pr < o.stockAnalyzeCfg.VolumePRLimit {
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
	case allOutInRation > o.stockAnalyzeCfg.AllOutInRatio && o.lastTick.Low < o.lastTick.Close:
		order.Action = entity.ActionBuy
		order.Price = o.lastBidAsk.AskPrice1
		return order
	case 100-allOutInRation > o.stockAnalyzeCfg.AllInOutRatio && o.lastTick.High > o.lastTick.Close:
		order.Action = entity.ActionSellFirst
		order.Price = o.lastBidAsk.BidPrice1
		return order
	default:
		return nil
	}
}

func (o *StockTrader) generateTradeOutOrder(cfg config.StockAnalyze, postOrderAction entity.OrderAction, preOrder *entity.StockOrder) *entity.StockOrder {
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

func (o *StockTrader) checkPlaceOrderStatus(order *entity.StockOrder) {
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

func (o *StockTrader) cancelOrder(order *entity.StockOrder) {
	order.TradeTime = time.Time{}
	bus.PublishTopicEvent(event.TopicCancelStockOrder, order)

	go func() {
		for {
			time.Sleep(time.Second)
			if order.TradeTime.IsZero() {
				continue
			}

			if order.Status == entity.StatusCancelled {
				log.Warnf("Order Canceled -> Stock: %s, Action: %d, Price: %.2f, Qty: %d", order.StockNum, order.Action, order.Price, order.Quantity)
				if order.Action == entity.ActionBuy || order.Action == entity.ActionSellFirst {
					bus.PublishTopicEvent(event.TopicUnSubscribeStockTickTargets, order.StockNum)
					return
				}
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

func (o *StockTrader) checkNeededPost() (entity.OrderAction, *entity.StockOrder) {
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

func (o *StockTrader) checkFirstTickArrive() {
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

func (o *StockTrader) getPRByVolume(volume int64) float64 {
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

// func (o *StockTrader) alreadyTrade() bool {
// 	defer o.orderMapLock.RUnlock()
// 	o.orderMapLock.RLock()
// 	return len(o.orderMap) != 0 && len(o.orderMap)%2 == 0
// }

// GetStockNum -/
func (o *StockTrader) GetStockNum() string {
	return o.stockNum
}

// GetTickChan -.
func (o *StockTrader) GetTickChan() chan *entity.RealTimeStockTick {
	return o.tickChan
}

// GetBidAskChan -.
func (o *StockTrader) GetBidAskChan() chan *entity.RealTimeStockBidAsk {
	return o.bidAskChan
}
