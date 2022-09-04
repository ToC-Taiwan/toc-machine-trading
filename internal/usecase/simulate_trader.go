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
	orderMap     map[entity.OrderAction][]*entity.StockOrder

	tickChan chan *entity.RealTimeTick

	historyTickAnalyze []int64
	analyzeTickTime    time.Time
	lastTick           *entity.RealTimeTick

	tradeSwitch  config.TradeSwitch
	simulateDone bool

	allowForward bool
	allowReverse bool
}

// NewSimulateAgent -.
func NewSimulateAgent(stockNum string) *SimulateTradeAgent {
	var quantity int64 = 2
	if biasRate := cc.GetBiasRate(stockNum); biasRate > cc.GetHighBiasRate() || biasRate < cc.GetLowBiasRate() {
		quantity = 1
	} else if biasRate == 0 {
		log.Errorf("%s BiasRate is 0", stockNum)
	}

	arr := cc.GetHistoryTickAnalyze(stockNum)
	sort.Slice(arr, func(i, j int) bool {
		return arr[i] > arr[j]
	})

	new := &SimulateTradeAgent{
		stockNum:           stockNum,
		orderQuantity:      quantity,
		orderMap:           make(map[entity.OrderAction][]*entity.StockOrder),
		tickChan:           make(chan *entity.RealTimeTick),
		historyTickAnalyze: arr,
	}

	if gap := cc.GetFutureGap(cc.GetBasicInfo().LastTradeDay); gap > 0 {
		new.allowForward = true
	} else if gap != 0 {
		new.allowReverse = true
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
	arr := *tickArr
	open := arr[0].Close
	high := arr[0].Close
	low := arr[0].Close

	// if open != beforeLastTradeDayClose {
	// 	o.simulateDone = true
	// }

	for _, tick := range arr[1:] {
		if tick.Close > high {
			high = tick.Close
		}

		if tick.Close < low {
			low = tick.Close
		}

		realTimeTick := &entity.RealTimeTick{
			StockNum: tick.StockNum,
			TickTime: tick.TickTime,
			Close:    tick.Close,
			Volume:   tick.Volume,
			TickType: tick.TickType,
			PctChg:   100 * (tick.Close - beforeLastTradeDayClose) / beforeLastTradeDayClose,
			High:     high,
			Low:      low,
			Open:     open,
		}

		if !o.simulateDone {
			o.tickChan <- realTimeTick
		} else {
			close(o.tickChan)
			break
		}
	}
}

func (o *SimulateTradeAgent) generateSimulateOrder(cfg config.Analyze) *entity.StockOrder {
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

	if o.lastTick.TickTime.Before(cc.GetBasicInfo().LastTradeDay.Add(9 * time.Hour).Add(time.Duration(o.tradeSwitch.HoldTimeFromOpen) * time.Second)) {
		return nil
	}

	if postOrderAction, preOrder := o.checkNeededPost(); postOrderAction != entity.ActionNone {
		return o.generateSimulateTradeOutOrder(cfg, postOrderAction, preOrder)
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
			TickTime:  o.lastTick.TickTime,
			Quantity:  o.orderQuantity,
			GroupID:   uuid.New().String(),
			TradeTime: o.lastTick.TickTime,
		},
	}

	switch {
	case allOutInRation > cfg.AllOutInRatio && o.lastTick.Low < o.lastTick.Close && o.allowForward:
		order.Action = entity.ActionBuy
		order.Price = o.lastTick.Close
		return order
	case 100-allOutInRation > cfg.AllInOutRatio && o.lastTick.High > o.lastTick.Close && o.allowReverse:
		order.Action = entity.ActionSellFirst
		order.Price = o.lastTick.Close
		return order
	default:
		return nil
	}
}

func (o *SimulateTradeAgent) generateSimulateTradeOutOrder(cfg config.Analyze, postOrderAction entity.OrderAction, preOrder *entity.StockOrder) *entity.StockOrder {
	order := &entity.StockOrder{
		StockNum: o.stockNum,
		BaseOrder: entity.BaseOrder{
			Action:    postOrderAction,
			Price:     o.lastTick.Close,
			Quantity:  preOrder.Quantity,
			TradeTime: o.lastTick.TickTime,
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
			position = total - 1
		}
	}
	return 100 * float64(total-position) / float64(total)
}

func (o *SimulateTradeAgent) checkNeededPost() (entity.OrderAction, *entity.StockOrder) {
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

func (o *SimulateTradeAgent) getAllOrders() []*entity.StockOrder {
	defer o.orderMapLock.RUnlock()
	o.orderMapLock.RLock()

	var orders []*entity.StockOrder
	for _, v := range o.orderMap {
		orders = append(orders, v...)
	}

	if len(orders)%2 != 0 {
		log.Warnf("Orders are not enough %s", o.stockNum)
		return []*entity.StockOrder{}
	}

	return orders
}
