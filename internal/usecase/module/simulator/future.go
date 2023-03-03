// Package simulator package simulator
package simulator

import (
	"time"

	"tmt/cmd/config"
	"tmt/internal/entity"
	"tmt/internal/usecase/module/quota"
	"tmt/internal/usecase/module/tradeday"

	"github.com/google/uuid"
)

type SimulatorFuture struct {
	code           string
	quantity       int64
	historyTickArr entity.RealTimeFutureTickArr

	tradeConfig *config.TradeFuture
	quota       *quota.Quota

	lastTick *entity.RealTimeFutureTick
	tickArr  entity.RealTimeFutureTickArr

	tickChan     chan *entity.RealTimeFutureTick
	lastTickRate float64

	lastPlaceOrderTime time.Time

	lastTradeOutTime time.Time
	waitingOrder     *entity.FutureOrder
	waitTimes        map[string]int

	allOrder []*entity.FutureOrder

	tradePeriod         tradeday.TradePeriod
	allowTradeTimeRange [][]time.Time
}

func NewSimulatorFuture(target SimulatorFutureTarget) *SimulatorFuture {
	s := &SimulatorFuture{
		code:        target.Code,
		quantity:    target.TradeConfig.Quantity,
		quota:       target.Quota,
		tradeConfig: target.TradeConfig,
		tradePeriod: target.TradePeriod,
		tickChan:    make(chan *entity.RealTimeFutureTick),
		waitTimes:   make(map[string]int),
	}

	firstStart := s.tradePeriod.StartTime
	secondStart := s.tradePeriod.EndTime.Add(-300 * time.Minute)

	s.allowTradeTimeRange = append(s.allowTradeTimeRange, []time.Time{firstStart, firstStart.Add(time.Duration(s.tradeConfig.TradeTimeRange.FirstPartDuration) * time.Minute)})
	s.allowTradeTimeRange = append(s.allowTradeTimeRange, []time.Time{secondStart, secondStart.Add(time.Duration(s.tradeConfig.TradeTimeRange.SecondPartDuration) * time.Minute)})

	for _, tick := range target.Ticks {
		s.historyTickArr = append(s.historyTickArr, tick.ToRealTimeTick())
	}

	s.ProcessTick()

	return s
}

func (s *SimulatorFuture) Start() {
	for _, tick := range s.historyTickArr {
		s.tickChan <- tick
	}
}

func (s *SimulatorFuture) ProcessTick() {
	go func() {
		for {
			s.lastTick = <-s.tickChan
			s.tickArr = append(s.tickArr, s.lastTick)

			s.cutTickArr()
			if s.waitingOrder != nil {
				s.checkByBalance(s.lastTick)
				continue
			}

			if o := s.generateOrder(); o != nil {
				s.lastPlaceOrderTime = s.lastTick.TickTime
				s.lastTradeOutTime = s.lastTick.TickTime.Add(time.Duration(s.tradeConfig.MaxHoldTime) * time.Minute)

				o.OrderTime = s.lastTick.TickTime
				o.OrderID = uuid.NewString()

				s.waitTimes[o.OrderID] = int(s.tradeConfig.TradeOutWaitTimes)
				s.waitingOrder = o
				s.allOrder = append(s.allOrder, o)
			}
		}
	}()
}

func (s *SimulatorFuture) cutTickArr() {
	if len(s.tickArr) < 2 {
		return
	}

	if s.tickArr[len(s.tickArr)-1].TickTime.Sub(s.tickArr[len(s.tickArr)-2].TickTime) > 3*time.Second {
		s.tickArr = entity.RealTimeFutureTickArr{}
		s.lastTickRate = 0
		return
	}

	if s.lastTick.TickTime.Sub(s.tickArr[0].TickTime) > time.Duration(2*s.tradeConfig.TickInterval)*time.Second {
		s.tickArr = s.tickArr[1:]
	}
}

func (s *SimulatorFuture) generateOrder() *entity.FutureOrder {
	if s.lastTick.TickTime.Sub(s.lastTradeOutTime) < 3*time.Minute || !s.isAllowTrade(s.lastTick.TickTime) {
		return nil
	}

	outInRatio, tickRate := s.tickArr.GetOutInRatioAndRate(time.Duration(s.tradeConfig.TickInterval) * time.Second)
	defer func() {
		s.lastTickRate = tickRate
	}()
	if s.lastTickRate == 0 {
		return nil
	}

	if s.lastTickRate < s.tradeConfig.RateLimit || 100*(tickRate-s.lastTickRate)/s.lastTickRate < s.tradeConfig.RateChangeRatio {
		return nil
	}

	switch {
	case outInRatio > s.tradeConfig.OutInRatio:
		return &entity.FutureOrder{
			Code: s.code,
			BaseOrder: entity.BaseOrder{
				Action:   entity.ActionBuy,
				Price:    s.lastTick.Close - 1,
				Quantity: s.tradeConfig.Quantity,
			},
		}
	case 100-outInRatio < s.tradeConfig.OutInRatio:
		return &entity.FutureOrder{
			Code: s.code,
			BaseOrder: entity.BaseOrder{
				Action:   entity.ActionSell,
				Price:    s.lastTick.Close + 1,
				Quantity: s.tradeConfig.Quantity,
			},
		}
	default:
		return nil
	}
}

func (s *SimulatorFuture) checkWaitTimes(tick *entity.RealTimeFutureTick) bool {
	defer func() {
		s.lastTick = tick
	}()

	if s.lastTick == nil {
		return true
	}

	switch s.waitingOrder.Action {
	case entity.ActionSell:
		times := s.waitTimes[s.waitingOrder.OrderID]
		if times > 0 && tick.Close <= s.lastTick.Close {
			s.waitTimes[s.waitingOrder.OrderID]--
			return true
		}

	case entity.ActionBuy:
		times := s.waitTimes[s.waitingOrder.OrderID]
		if times > 0 && tick.Close >= s.lastTick.Close {
			s.waitTimes[s.waitingOrder.OrderID]--
			return true
		}
	}

	return false
}

func (s *SimulatorFuture) checkByBalance(tick *entity.RealTimeFutureTick) {
	if s.checkWaitTimes(tick) {
		return
	}

	var place bool
	switch s.waitingOrder.Action {
	case entity.ActionSell:
		if tick.Close <= s.waitingOrder.Price-s.tradeConfig.TargetBalanceHigh || tick.Close >= s.waitingOrder.Price-s.tradeConfig.TargetBalanceLow {
			place = true
		}

	case entity.ActionBuy:
		if tick.Close >= s.waitingOrder.Price+s.tradeConfig.TargetBalanceHigh || tick.Close <= s.waitingOrder.Price+s.tradeConfig.TargetBalanceLow {
			place = true
		}
	}

	if !place && time.Now().Before(s.lastTradeOutTime) {
		return
	}

	o := &entity.FutureOrder{
		Code: tick.Code,
		BaseOrder: entity.BaseOrder{
			Price:    tick.Close,
			Quantity: s.waitingOrder.Quantity,
		},
	}

	switch s.waitingOrder.Action {
	case entity.ActionSell:
		o.Action = entity.ActionBuy
	case entity.ActionBuy:
		o.Action = entity.ActionSell
	}

	o.OrderTime = tick.TickTime
	o.OrderID = uuid.NewString()
	s.allOrder = append(s.allOrder, o)
	s.waitingOrder = nil
}

func (s *SimulatorFuture) CalculateFutureTradeBalance(notifyChan chan SimulateBalance) {
	var forward, reverse []*entity.FutureOrder
	qtyMap := make(map[string]int64)
	for _, v := range s.allOrder {
		switch v.Action {
		case entity.ActionBuy:
			if qtyMap[v.Code] >= 0 {
				forward = append(forward, v)
			} else {
				reverse = append(reverse, v)
			}
			qtyMap[v.Code] += v.Quantity
		case entity.ActionSell:
			if qtyMap[v.Code] > 0 {
				forward = append(forward, v)
			} else {
				reverse = append(reverse, v)
			}
			qtyMap[v.Code] -= v.Quantity
		}
	}

	forwardBalance, _ := s.calculateForwardFutureBalance(forward)
	revereBalance, _ := s.calculateReverseFutureBalance(reverse)

	notifyChan <- SimulateBalance{
		TotalBalance: forwardBalance + revereBalance,
		Forward:      forwardBalance,
		Reverse:      revereBalance,
		Cond:         s.tradeConfig,
	}
}

func (s *SimulatorFuture) calculateForwardFutureBalance(forward []*entity.FutureOrder) (int64, int64) {
	var forwardBalance, tradeCount int64
	var qty int64
	for _, v := range forward {
		tradeCount++

		switch v.Action {
		case entity.ActionBuy:
			qty += v.Quantity
			forwardBalance -= s.quota.GetFutureBuyCost(v.Price, v.Quantity)
		case entity.ActionSell:
			qty -= v.Quantity
			forwardBalance += s.quota.GetFutureSellCost(v.Price, v.Quantity)
		}
	}

	if qty != 0 {
		return 0, tradeCount
	}

	return forwardBalance, tradeCount
}

func (s *SimulatorFuture) calculateReverseFutureBalance(reverse []*entity.FutureOrder) (int64, int64) {
	var reverseBalance, tradeCount int64
	var qty int64
	for _, v := range reverse {
		tradeCount++

		switch v.Action {
		case entity.ActionSell:
			qty -= v.Quantity
			reverseBalance += s.quota.GetFutureSellCost(v.Price, v.Quantity)
		case entity.ActionBuy:
			qty += v.Quantity
			reverseBalance -= s.quota.GetFutureBuyCost(v.Price, v.Quantity)
		}
	}

	if qty != 0 {
		return 0, tradeCount
	}

	return reverseBalance, tradeCount
}

func (s *SimulatorFuture) isAllowTrade(tickTime time.Time) bool {
	var tempSwitch bool
	for _, rangeTime := range s.allowTradeTimeRange {
		if tickTime.After(rangeTime[0]) && tickTime.Before(rangeTime[1]) {
			tempSwitch = true
		}
	}
	return tempSwitch
}
