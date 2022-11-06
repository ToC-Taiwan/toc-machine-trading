// Package trader package trader
package trader

import (
	"sort"
	"time"

	"tmt/cmd/config"
	"tmt/internal/entity"
	"tmt/internal/usecase/modules/quota"
	"tmt/internal/usecase/modules/tradeday"
	"tmt/pkg/common"

	"github.com/google/uuid"
)

// FutureSimulator -.
type FutureSimulator struct {
	code          string
	orderQuantity int64

	orderMap map[entity.OrderAction][]*entity.FutureOrder

	quota *quota.Quota

	tickArr realTimeFutureTickArr

	lastTick *entity.RealTimeFutureTick
	tickChan chan *entity.RealTimeFutureTick

	analyzeCfg config.FutureAnalyze

	simDone bool

	firstTradePeriod  []time.Time
	secondTradePeriod []time.Time
}

// NewFutureSimulator -.
func NewFutureSimulator(code string, analyzeCfg config.FutureAnalyze, period tradeday.TradePeriod) *FutureSimulator {
	cfg := config.GetConfig()
	t := &FutureSimulator{
		code:              code,
		orderQuantity:     2,
		orderMap:          make(map[entity.OrderAction][]*entity.FutureOrder),
		tickChan:          make(chan *entity.RealTimeFutureTick),
		analyzeCfg:        analyzeCfg,
		quota:             quota.NewQuota(cfg.Quota),
		firstTradePeriod:  []time.Time{period.StartTime, period.StartTime.Add(time.Duration(cfg.FutureTradeSwitch.TradeTimeRange.FirstPartDuration) * time.Minute)},
		secondTradePeriod: []time.Time{period.EndTime.Add(-300 * time.Minute), period.EndTime.Add(-300 * time.Minute).Add(time.Duration(cfg.FutureTradeSwitch.TradeTimeRange.SecondPartDuration) * time.Minute)},
	}

	go t.SimulateRoom()
	return t
}

// SimulateRoom -.
func (o *FutureSimulator) SimulateRoom() {
	for {
		tick := <-o.tickChan
		o.tickArr = append(o.tickArr, tick)
		if len(o.tickArr) == 2 {
			o.tickArr = o.tickArr[1:]
			o.lastTick = tick
			break
		}
	}

	for {
		tick, ok := <-o.tickChan
		if !ok {
			o.simDone = true
			break
		}

		if tick.TickTime.Minute() != o.lastTick.TickTime.Minute() {
			o.placeFutureOrder(o.generateOrder())
		}

		o.lastTick = tick
		o.tickArr = append(o.tickArr, tick)
	}
}

func (o *FutureSimulator) generateOrder() *entity.FutureOrder {
	if postOrderAction, preOrder := o.checkNeededPost(); postOrderAction != entity.ActionNone {
		return o.generateTradeOutOrder(postOrderAction, preOrder)
	}

	rsi := o.tickArr.getRSIByTickCount(300)
	if rsi == 0 {
		return nil
	}

	// get out in ration in period
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
	case rsi <= 40:
		order.Action = entity.ActionBuy
	case rsi >= 60:
		order.Action = entity.ActionSellFirst
	default:
		return nil
	}
	return order
}

func (o *FutureSimulator) generateTradeOutOrder(postOrderAction entity.OrderAction, preOrder *entity.FutureOrder) *entity.FutureOrder {
	order := &entity.FutureOrder{
		Code: o.code,
		BaseOrder: entity.BaseOrder{
			Action:   postOrderAction,
			Price:    o.lastTick.Close,
			Quantity: preOrder.Quantity,
			TickTime: o.lastTick.TickTime,
			GroupID:  preOrder.GroupID,
		},
	}

	if o.lastTick.TickTime.After(preOrder.TickTime.Add(time.Duration(o.analyzeCfg.MaxHoldTime) * time.Minute)) {
		return order
	}

	return nil
}

func (o *FutureSimulator) placeFutureOrder(order *entity.FutureOrder) {
	if order == nil {
		return
	}

	if order.Action == entity.ActionSell || order.Action == entity.ActionBuyLater {
		o.orderMap[order.Action] = append(o.orderMap[order.Action], order)
		return
	}

	if (order.TickTime.Before(o.firstTradePeriod[0]) || order.TickTime.After(o.firstTradePeriod[1])) && order.TickTime.Before(o.secondTradePeriod[0]) || order.TickTime.After(o.secondTradePeriod[1]) {
		return
	}

	o.orderMap[order.Action] = append(o.orderMap[order.Action], order)
}

// GetTickChan -.
func (o *FutureSimulator) GetTickChan() chan *entity.RealTimeFutureTick {
	return o.tickChan
}

func (o *FutureSimulator) checkNeededPost() (entity.OrderAction, *entity.FutureOrder) {
	if len(o.orderMap[entity.ActionBuy]) > len(o.orderMap[entity.ActionSell]) {
		return entity.ActionSell, o.orderMap[entity.ActionBuy][len(o.orderMap[entity.ActionSell])]
	}

	if len(o.orderMap[entity.ActionSellFirst]) > len(o.orderMap[entity.ActionBuyLater]) {
		return entity.ActionBuyLater, o.orderMap[entity.ActionSellFirst][len(o.orderMap[entity.ActionBuyLater])]
	}

	return entity.ActionNone, nil
}

// CalculateFutureTradeBalance -.
func (o *FutureSimulator) CalculateFutureTradeBalance() TradeBalance {
	for {
		if o.simDone {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	var orderList []*entity.FutureOrder
	for _, v := range o.orderMap {
		orderList = append(orderList, v...)
	}

	if len(orderList) == 0 {
		return TradeBalance{}
	}

	sort.Slice(orderList, func(i, j int) bool {
		return orderList[i].TickTime.Before(orderList[j].TickTime)
	})

	tradeOutOrderMap := make(map[string]*entity.FutureOrder)
	for _, v := range orderList {
		if v.Action == entity.ActionSell || v.Action == entity.ActionBuyLater {
			tradeOutOrderMap[v.GroupID] = v
		}
	}

	var forwardBalance, revereBalance, tradeCount int64
	for _, v := range orderList {
		switch v.Action {
		case entity.ActionBuy:
			tradeCount++
			balance := -o.quota.GetFutureBuyCost(v.Price, v.Quantity) + o.quota.GetFutureSellCost(tradeOutOrderMap[v.GroupID].Price, tradeOutOrderMap[v.GroupID].Quantity)
			forwardBalance += balance
			log.Warnf("#%3d %s -> %5d (f:%s)", tradeCount, v.TickTime.Format(common.LongTimeLayout), balance, tradeOutOrderMap[v.GroupID].TickTime.Sub(v.TickTime).String())

		case entity.ActionSellFirst:
			tradeCount++
			balance := o.quota.GetFutureSellCost(v.Price, v.Quantity) - o.quota.GetFutureBuyCost(tradeOutOrderMap[v.GroupID].Price, tradeOutOrderMap[v.GroupID].Quantity)
			revereBalance += balance
			log.Warnf("#%3d %s -> %5d (r:%s)", tradeCount, v.TickTime.Format(common.LongTimeLayout), balance, tradeOutOrderMap[v.GroupID].TickTime.Sub(v.TickTime).String())
		}
	}

	log.Warnf("#%3d Total: %d", tradeCount+1, forwardBalance+revereBalance)
	return TradeBalance{
		Count:   tradeCount,
		Balance: forwardBalance + revereBalance,
	}
}
