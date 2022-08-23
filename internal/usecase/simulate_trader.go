package usecase

import (
	"sort"
	"sync"
	"time"

	"tmt/internal/entity"
	"tmt/pkg/config"

	"github.com/google/uuid"
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

	tradeSwitch  config.TradeSwitch
	simulateDone bool
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

func (o *SimulateTradeAgent) searchOrder(cfg config.Analyze, tickArr *[]*entity.HistoryTick, beforeLastTradeDayClose float64) {
	go func() {
		for {
			tick, ok := <-o.tickChan
			if !ok {
				break
			}

			o.lastTick = tick
			o.tickArr = append(o.tickArr, o.lastTick)

			order := o.generateSimulateOrder(cfg)
			if order == nil {
				continue
			}

			if (order.Action == entity.ActionBuy || order.Action == entity.ActionSellFirst) && o.lastTick.TickTime.After(cc.GetBasicInfo().LastTradeDay.Add(9*time.Hour).Add(time.Duration(o.tradeSwitch.TradeInEndTime)*time.Minute)) {
				o.simulateDone = true
				continue
			}

			o.orderMapLock.Lock()
			o.orderMap[order.Action] = append(o.orderMap[order.Action], order)
			o.orderMapLock.Unlock()
		}
	}()

	o.convertToRealTimeTick(tickArr, beforeLastTradeDayClose)
}

func (o *SimulateTradeAgent) convertToRealTimeTick(tickArr *[]*entity.HistoryTick, beforeLastTradeDayClose float64) {
	for i, tick := range *tickArr {
		realTimeTick := &entity.RealTimeTick{
			StockNum: tick.StockNum,
			TickTime: tick.TickTime,
			Close:    tick.Close,
			Volume:   tick.Volume,
			TickType: tick.TickType,
			PctChg:   100 * (tick.Close - beforeLastTradeDayClose) / beforeLastTradeDayClose,
		}

		if i == 0 && realTimeTick.Close != beforeLastTradeDayClose {
			o.simulateDone = true
		}

		if !o.simulateDone {
			o.tickChan <- realTimeTick
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

	if postOrderAction, preOrder := o.checkNeededPost(); postOrderAction != entity.ActionNone {
		return o.generateSimulateTradeOutOrder(cfg, postOrderAction, preOrder)
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
	allOutInRation := o.tickArr.getOutInRatio()

	// need to compare with all and period
	order := &entity.Order{
		StockNum: o.stockNum,
		TickTime: o.lastTick.TickTime,
		Quantity: o.orderQuantity,
		GroupID:  uuid.New().String(),

		TradeTime: o.lastTick.TickTime,
	}

	switch {
	case periodOutInRation-allOutInRation > cfg.AllOutInRatio*0.1 && allOutInRation > cfg.AllOutInRatio:
		order.Action = entity.ActionBuy
		order.Price = o.lastTick.Close
	case allOutInRation-periodOutInRation > cfg.AllInOutRatio*0.1 && 100-allOutInRation > cfg.AllInOutRatio:
		order.Action = entity.ActionSellFirst
		order.Price = o.lastTick.Close
	default:
		return nil
	}

	return order
}

func (o *SimulateTradeAgent) generateSimulateTradeOutOrder(cfg config.Analyze, postOrderAction entity.OrderAction, preOrder *entity.Order) *entity.Order {
	if o.lastTick.TickTime.After(preOrder.TradeTime.Add(time.Duration(cfg.MaxHoldTime) * time.Minute)) {
		return &entity.Order{
			StockNum:  o.stockNum,
			Action:    postOrderAction,
			Price:     o.lastTick.Close,
			Quantity:  preOrder.Quantity,
			TradeTime: o.lastTick.TickTime,
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

func (o *SimulateTradeAgent) checkNeededPost() (entity.OrderAction, *entity.Order) {
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

func (o *SimulateTradeAgent) getAllOrders() []*entity.Order {
	defer o.orderMapLock.RUnlock()
	o.orderMapLock.RLock()

	var orders []*entity.Order
	for _, v := range o.orderMap {
		orders = append(orders, v...)
	}

	if len(orders)%2 != 0 {
		log.Warnf("Orders are not enough %s", o.stockNum)
		return []*entity.Order{}
	}

	return orders
}
