package usecase

import (
	"sort"
	"sync"
	"time"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/pkg/config"
)

// SimulateTradeAgent -.
type SimulateTradeAgent struct {
	stockNum      string
	orderQuantity int64
	tickArr       RealTimeTickArr
	periodTickArr RealTimeTickArr

	orderMapLock sync.RWMutex
	orderMap     map[entity.OrderAction][]*entity.Order

	tickChan chan *entity.RealTimeTick

	historyTickAnalyze []int64
	analyzeTickTime    time.Time
	lastTick           *entity.RealTimeTick

	tradeSwitch config.TradeSwitch
}

// NewSimulateAgent -.
func NewSimulateAgent(stockNum string) *SimulateTradeAgent {
	var quantity int64 = 1
	if biasRate := cc.GetBiasRate(stockNum); biasRate > 4 || biasRate < -4 {
		quantity = 2
	}

	arr := cc.GetHistoryTickAnalyze(stockNum)
	sort.Slice(arr, func(i, j int) bool {
		return arr[i] > arr[j]
	})

	new := &SimulateTradeAgent{
		stockNum:           stockNum,
		orderQuantity:      quantity,
		orderMap:           make(map[entity.OrderAction][]*entity.Order),
		tickChan:           make(chan *entity.RealTimeTick),
		historyTickAnalyze: arr,
	}

	return new
}

// ResetAgent -.
func (o *SimulateTradeAgent) ResetAgent(stockNum string) {
	o.stockNum = stockNum
	o.orderQuantity = 1
	o.tickArr = []*entity.RealTimeTick{}
	o.periodTickArr = []*entity.RealTimeTick{}
	o.orderMap = make(map[entity.OrderAction][]*entity.Order)
	o.tickChan = make(chan *entity.RealTimeTick)
	o.historyTickAnalyze = []int64{}
	o.analyzeTickTime = time.Time{}
	o.lastTick = &entity.RealTimeTick{}
	o.tradeSwitch = config.TradeSwitch{}
}

func (o *SimulateTradeAgent) searchOrder(cfg config.Analyze, tickArr *[]*entity.HistoryTick, beforeLastTradeDayClose float64) {
	var finish bool
	go func() {
		for {
			tick, ok := <-o.tickChan
			if !ok {
				break
			}

			o.lastTick = tick
			o.tickArr = append(o.tickArr, o.lastTick)

			order := o.generateSimulateOrder(cfg)
			if order == nil || finish {
				continue
			}

			o.orderMapLock.Lock()
			o.orderMap[order.Action] = append(o.orderMap[order.Action], order)
			o.orderMapLock.Unlock()

			if order.Action == entity.ActionSell || order.Action == entity.ActionBuyLater {
				finish = true
			}
		}
	}()

	for _, tick := range *tickArr {
		if !finish {
			o.tickChan <- &entity.RealTimeTick{
				StockNum: tick.StockNum,
				TickTime: tick.TickTime,
				Close:    tick.Close,
				Volume:   tick.Volume,
				TickType: tick.TickType,
				PctChg:   100 * (tick.Close - beforeLastTradeDayClose) / beforeLastTradeDayClose,
			}
		} else {
			close(o.tickChan)
			break
		}
	}
}

func (o *SimulateTradeAgent) generateSimulateOrder(cfg config.Analyze) *entity.Order {
	if o.lastTick.TickTime.Sub(o.analyzeTickTime) < time.Duration(cfg.TickAnalyzePeriod)*time.Millisecond {
		o.periodTickArr = append(o.periodTickArr, o.lastTick)
		return nil
	}
	// copy new arr before reset
	analyzeArr := o.periodTickArr
	// reset analyze tick time and arr
	o.analyzeTickTime = o.lastTick.TickTime
	o.periodTickArr = RealTimeTickArr{o.lastTick}

	if o.lastTick.TickTime.Before(cc.GetBasicInfo().LastTradeDay.Add(9 * time.Hour).Add(time.Duration(o.tradeSwitch.HoldTimeFromOpen) * time.Second)) {
		return nil
	}

	if postOrderAction, qty, preTime := o.checkNeededPost(); postOrderAction != entity.ActionNone {
		return o.generateSimulateTradeOutOrder(cfg, postOrderAction, qty, preTime)
	}

	if o.lastTick.PctChg < cfg.CloseChangeRatioLow || o.lastTick.PctChg > cfg.CloseChangeRatioHigh {
		return nil
	}

	periodVolume := analyzeArr.getTotalVolume()
	if pr := o.getPRByVolume(periodVolume); pr < cfg.VolumePRLimit {
		return nil
	}

	// get out in ration in this period
	periodOutInRation := analyzeArr.getOutInRatio()

	// need to compare with all and period
	order := &entity.Order{
		StockNum:  o.stockNum,
		Quantity:  o.orderQuantity,
		TradeTime: o.lastTick.TickTime,
	}

	switch {
	case periodOutInRation > cfg.OutInRatio:
		order.Action = entity.ActionBuy
		order.Price = o.lastTick.Close
	case 100-periodOutInRation > cfg.InOutRatio:
		order.Action = entity.ActionSellFirst
		order.Price = o.lastTick.Close
	default:
		return nil
	}

	return order
}

func (o *SimulateTradeAgent) generateSimulateTradeOutOrder(cfg config.Analyze, postOrderAction entity.OrderAction, qty int64, preTime time.Time) *entity.Order {
	rsi := o.tickArr.getRSIByTickTime(preTime, cfg.RSIMinCount)
	if rsi != 0 {
		switch postOrderAction {
		case entity.ActionSell:
			if rsi >= cfg.RSIHigh {
				return &entity.Order{
					StockNum:  o.stockNum,
					Action:    postOrderAction,
					Price:     o.lastTick.Close,
					Quantity:  qty,
					TradeTime: o.lastTick.TickTime,
				}
			}
		case entity.ActionBuyLater:
			if rsi <= cfg.RSILow {
				return &entity.Order{
					StockNum:  o.stockNum,
					Action:    postOrderAction,
					Price:     o.lastTick.Close,
					Quantity:  qty,
					TradeTime: o.lastTick.TickTime,
				}
			}
		}
	}
	return nil
}

func (o *SimulateTradeAgent) getPRByVolume(volume int64) float64 {
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

func (o *SimulateTradeAgent) checkNeededPost() (entity.OrderAction, int64, time.Time) {
	defer o.orderMapLock.RUnlock()
	o.orderMapLock.RLock()

	if len(o.orderMap[entity.ActionBuy]) > len(o.orderMap[entity.ActionSell]) {
		order := o.orderMap[entity.ActionBuy][len(o.orderMap[entity.ActionSell])]
		return entity.ActionSell, order.Quantity, order.TradeTime
	}

	if len(o.orderMap[entity.ActionSellFirst]) > len(o.orderMap[entity.ActionBuyLater]) {
		order := o.orderMap[entity.ActionSellFirst][len(o.orderMap[entity.ActionBuyLater])]
		return entity.ActionBuyLater, order.Quantity, order.TradeTime
	}

	return entity.ActionNone, 0, time.Time{}
}

func (o *SimulateTradeAgent) getAllOrders() []*entity.Order {
	defer o.orderMapLock.RUnlock()
	o.orderMapLock.RLock()

	var orders []*entity.Order
	for _, v := range o.orderMap {
		orders = append(orders, v...)
	}

	if len(orders)%2 != 0 {
		return []*entity.Order{}
	}

	return orders
}
