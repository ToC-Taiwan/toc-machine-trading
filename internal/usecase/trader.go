package usecase

import (
	"sort"
	"sync"
	"time"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/pkg/config"
	"toc-machine-trading/pkg/utils"
)

// TradeAgent -.
type TradeAgent struct {
	stockNum      string
	orderQuantity int64
	tickArr       RealTimeTickArr

	orderMapLock sync.RWMutex
	orderMap     map[entity.OrderAction][]*entity.Order
	waitingOrder *entity.Order

	tickChan   chan *entity.RealTimeTick
	bidAskChan chan *entity.RealTimeBidAsk

	historyTickAnalyze []int64
	analyzeTickTime    time.Time
	lastTick           *entity.RealTimeTick
	lastBidAsk         *entity.RealTimeBidAsk
}

// NewAgent -.
func NewAgent(stockNum string) *TradeAgent {
	var quantity int64 = 1
	if biasRate := cc.GetBiasRate(stockNum); biasRate > 4 || biasRate < -4 {
		quantity = 2
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
	}

	go new.checkFirstTickArrive()

	return new
}

func (o *TradeAgent) generateOrder(cfg config.Analyze, needClear bool) *entity.Order {
	if needClear || o.waitingOrder != nil || o.analyzeTickTime.IsZero() {
		return nil
	}

	if postOrderAction, preTime := o.checkNeededPost(); postOrderAction != entity.ActionNone {
		return o.generateTradeOutOrder(cfg, postOrderAction, preTime)
	}

	if o.lastTick.TickTime.Before(o.analyzeTickTime.Add(time.Duration(cfg.TickAnalyzeMinPeriod) * time.Millisecond)) {
		return nil
	}
	o.analyzeTickTime = o.lastTick.TickTime

	period := o.tickArr.getLastNMilliSecondArr(cfg.TickAnalyzeMinPeriod)
	periodVolume := period.getTotalVolume()
	if pr := o.getPRByVolume(periodVolume); pr < cfg.VolumePRLow || pr > cfg.VolumePRHigh {
		return nil
	}
	periodOutInRation := period.getOutInRatio()

	// need to compare with all and period
	order := &entity.Order{
		StockNum:  o.stockNum,
		Quantity:  o.orderQuantity,
		TradeTime: time.Now(),
	}

	switch {
	case periodOutInRation > cfg.OutInRatio && o.lastBidAsk != nil:
		order.Action = entity.ActionBuy
		order.Price = o.lastBidAsk.BidPrice1
	case 100-periodOutInRation < cfg.InOutRatio && o.lastBidAsk != nil:
		order.Action = entity.ActionSellFirst
		order.Price = o.lastBidAsk.AskPrice1
	default:
		return nil
	}

	return order
}

func (o *TradeAgent) generateTradeOutOrder(cfg config.Analyze, postOrderAction entity.OrderAction, preTime time.Time) *entity.Order {
	// calculate max loss here
	//
	rsi := o.tickArr.getRSIByTickTime(preTime, cfg.RSIMinCount)
	if rsi != 0 {
		switch postOrderAction {
		case entity.ActionSell:
			if rsi >= cfg.RSIHigh {
				return &entity.Order{
					StockNum:  o.stockNum,
					Action:    postOrderAction,
					Price:     o.lastBidAsk.BidPrice1,
					Quantity:  o.orderQuantity,
					TradeTime: time.Now(),
				}
			}
		case entity.ActionBuyLater:
			if rsi <= cfg.RSILow {
				return &entity.Order{
					StockNum:  o.stockNum,
					Action:    postOrderAction,
					Price:     o.lastBidAsk.AskPrice1,
					Quantity:  o.orderQuantity,
					TradeTime: time.Now(),
				}
			}
		}
	}
	return nil
}

func (o *TradeAgent) clearUnfinishedOrder() *entity.Order {
	if o.waitingOrder != nil {
		return nil
	}

	if action, _ := o.checkNeededPost(); action != entity.ActionNone {
		return &entity.Order{
			StockNum:  o.stockNum,
			Action:    action,
			Price:     o.lastTick.Close,
			Quantity:  o.orderQuantity,
			TradeTime: time.Now(),
		}
	}
	return nil
}

func (o *TradeAgent) checkPlaceOrderStatus(order *entity.Order, timeout time.Duration) {
	for {
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
		time.Sleep(time.Second)
	}

	if order.Status == entity.StatusAborted || order.Status == entity.StatusFailed {
		o.waitingOrder = nil
		return
	}

	if order.OrderID != "" && order.Status != entity.StatusCancelled && order.Status != entity.StatusFilled {
		log.Warnf("Place Cancel Order -> Stock: %s, Action: %d, Price: %.2f, Qty: %d", order.StockNum, order.Action, order.Price, order.Quantity)
		bus.PublishTopicEvent(topicCancelOrder, order.OrderID)
		go o.checkCancelOrder(order.OrderID, timeout)
		return
	}

	log.Error("check place order status raise unknown error")
}

func (o *TradeAgent) checkCancelOrder(orderID string, timeout time.Duration) {
	for {
		order := cc.GetOrderByOrderID(orderID)
		if order.Status == entity.StatusCancelled {
			log.Warnf("Order Canceled -> Stock: %s, Action: %d, Price: %.2f, Qty: %d", order.StockNum, order.Action, order.Price, order.Quantity)
			o.waitingOrder = nil
			break
		} else if order.TradeTime.Add(timeout).Before(time.Now()) {
			log.Errorf("Cancel Order Timeout -> Stock: %s, Action: %d, Price: %.2f, Qty: %d", order.StockNum, order.Action, order.Price, order.Quantity)
			go o.checkCancelOrder(orderID, timeout)
			return
		}
		time.Sleep(3 * time.Second)
	}
}

func (o *TradeAgent) checkNeededPost() (entity.OrderAction, time.Time) {
	defer o.orderMapLock.RUnlock()
	o.orderMapLock.RLock()

	if len(o.orderMap[entity.ActionBuy]) > len(o.orderMap[entity.ActionSell]) {
		return entity.ActionSell, o.orderMap[entity.ActionBuy][len(o.orderMap[entity.ActionSell])].TradeTime
	}

	if len(o.orderMap[entity.ActionSellFirst]) > len(o.orderMap[entity.ActionBuyLater]) {
		return entity.ActionBuyLater, o.orderMap[entity.ActionSellFirst][len(o.orderMap[entity.ActionBuyLater])].TradeTime
	}

	return entity.ActionNone, time.Time{}
}

func (o *TradeAgent) checkFirstTickArrive() {
	// calculate open change rate here
	//
	tradeDay := cc.GetBasicInfo().TradeDay
	for {
		if len(o.tickArr) != 0 {
			cc.SetHistoryOpen(o.stockNum, tradeDay, o.tickArr[0].Close)
			o.analyzeTickTime = o.tickArr[0].TickTime
			break
		}
		time.Sleep(200 * time.Millisecond)
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

// RealTimeTickArr -.
type RealTimeTickArr []*entity.RealTimeTick

func (c RealTimeTickArr) getLastNMilliSecondArr(n float64) RealTimeTickArr {
	if len(c) < 2 {
		return RealTimeTickArr{}
	}

	startTime := c[len(c)-1].TickTime

	// skip if i == 0, the volume will be too large
	// skip the first tick of the day
	var cut int
	for i := len(c) - 1; i > 0; i-- {
		if startTime.Sub(c[i].TickTime) < time.Duration(n)*time.Millisecond {
			continue
		} else {
			cut = i - 1
			break
		}
	}

	if cut == 0 {
		return RealTimeTickArr{}
	}
	return c[cut:]
}

func (c RealTimeTickArr) getTotalVolume() int64 {
	if len(c) == 0 {
		return 0
	}

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

	if len(tmp) < count {
		return 0
	}

	rsi, err := utils.GenerateRSI(tmp)
	if err != nil {
		log.Error(err)
		return 0
	}
	return rsi
}
