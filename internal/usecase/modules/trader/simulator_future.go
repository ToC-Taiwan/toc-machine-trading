// Package trader package trader
package trader

import (
	"sort"
	"time"

	"tmt/cmd/config"
	"tmt/global"
	"tmt/internal/entity"
	"tmt/internal/usecase/modules/quota"
	"tmt/internal/usecase/modules/tradeday"

	"github.com/google/uuid"
)

// FutureSimulator -.
type FutureSimulator struct {
	code          string
	orderQuantity int64

	orderMap map[entity.OrderAction][]*entity.FutureOrder

	quota *quota.Quota

	tickArr realTimeFutureTickArr
	kbarArr realTimeKbarArr

	lastTick *entity.RealTimeFutureTick
	tickChan chan *entity.RealTimeFutureTick

	analyzeCfg config.FutureAnalyze

	tradeOutRecord map[string]int
	simDone        bool

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
		tradeOutRecord:    make(map[string]int),
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
		o.lastTick = tick
		o.tickArr = append(o.tickArr, tick)
		break
	}

	for {
		tick, ok := <-o.tickChan
		if !ok {
			o.simDone = true
			break
		}
		if tick.TickTime.Minute() != o.lastTick.TickTime.Minute() {
			o.kbarArr = append(o.kbarArr, o.tickArr.getKbar())
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

	if !o.kbarArr.isStable(10) {
		return nil
	}

	splitBySecondArr := o.tickArr.splitBySecond(10)
	if splitBySecondArr == nil {
		return nil
	}

	base := splitBySecondArr[0].getTotalVolume()
	for i := 1; i <= len(splitBySecondArr)-1; i++ {
		if splitBySecondArr[i].getTotalVolume() > base*2 {
			return nil
		}
	}

	// get out in ration in period
	outInRation := splitBySecondArr[0].getOutInRatio()
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
	case outInRation >= o.analyzeCfg.AllOutInRatio:
		order.Action = entity.ActionBuy
		return order
	case 100-outInRation >= o.analyzeCfg.AllInOutRatio:
		order.Action = entity.ActionSellFirst
		return order
	default:
		return nil
	}
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

	switch order.Action {
	case entity.ActionSell:
		if order.Price-preOrder.Price < -2 {
			if o.tradeOutRecord[order.GroupID] >= 5 {
				return order
			}
			o.tradeOutRecord[order.GroupID]++
		}

		if order.Price-preOrder.Price > 10 {
			return order
		}

		// if order.Price-preOrder.Price > 5 || order.Price-preOrder.Price < -2 {
		// 	return order
		// }

	case entity.ActionBuyLater:
		if order.Price-preOrder.Price > 2 {
			if o.tradeOutRecord[order.GroupID] >= 5 {
				return order
			}
			o.tradeOutRecord[order.GroupID]++
		}

		if order.Price-preOrder.Price < -10 {
			return order
		}

		// if order.Price-preOrder.Price < -5 || order.Price-preOrder.Price > 2 {
		// 	return order
		// }
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
			log.Warnf("#%3d %s -> %5d (f:%s)", tradeCount, v.TickTime.Format(global.LongTimeLayout), balance, tradeOutOrderMap[v.GroupID].TickTime.Sub(v.TickTime).String())

		case entity.ActionSellFirst:
			tradeCount++
			balance := o.quota.GetFutureSellCost(v.Price, v.Quantity) - o.quota.GetFutureBuyCost(tradeOutOrderMap[v.GroupID].Price, tradeOutOrderMap[v.GroupID].Quantity)
			revereBalance += balance
			log.Warnf("#%3d %s -> %5d (r:%s)", tradeCount, v.TickTime.Format(global.LongTimeLayout), balance, tradeOutOrderMap[v.GroupID].TickTime.Sub(v.TickTime).String())
		}
	}

	log.Warnf("#  TradeCount: %d, Total: %d", tradeCount, forwardBalance+revereBalance)
	return TradeBalance{
		Count:   tradeCount,
		Balance: forwardBalance + revereBalance,
	}
}