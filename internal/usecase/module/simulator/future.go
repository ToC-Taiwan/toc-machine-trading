// Package simulator package simulator
package simulator

import (
	"fmt"
	"time"

	"tmt/cmd/config"
	"tmt/internal/entity"
	"tmt/internal/usecase/module/quota"
	"tmt/internal/usecase/module/tradeday"
	"tmt/pkg/log"

	"github.com/google/uuid"
)

var logger = log.Get()

type SimulatorFuture struct {
	conds       []*config.TradeFuture
	tradeConfig *config.TradeFuture

	code                string
	quantity            int64
	tradePeriod         tradeday.TradePeriod
	allowTradeTimeRange [][]time.Time
	quota               *quota.Quota

	lastTickRate       float64
	lastPlaceOrderTime time.Time
	maxTradeOutTime    time.Time

	allOrder     []*entity.FutureOrder
	waitingOrder *entity.FutureOrder

	historyTickArr entity.RealTimeFutureTickArr
	tickArr        entity.RealTimeFutureTickArr
	lastTick       *entity.RealTimeFutureTick
	tickChan       chan *entity.RealTimeFutureTick
	waitTimes      map[string]int
}

func NewSimulatorFuture(target SimulatorFutureTarget) *SimulatorFuture {
	if len(target.TradeConfigArr) == 0 {
		logger.Fatal("trade config arr is empty")
	}

	return &SimulatorFuture{
		code:           target.Code,
		quota:          target.Quota,
		tradePeriod:    target.TradePeriod,
		historyTickArr: target.Ticks,
		conds:          target.TradeConfigArr,
	}
}

func (s *SimulatorFuture) clear() {
	s.waitTimes = make(map[string]int)
	s.tickChan = make(chan *entity.RealTimeFutureTick)

	s.allowTradeTimeRange = [][]time.Time{}
	s.allOrder = []*entity.FutureOrder{}

	s.lastTickRate = 0
	s.waitingOrder = nil

	s.tickArr = entity.RealTimeFutureTickArr{}
	s.lastPlaceOrderTime = time.Time{}
}

func (s *SimulatorFuture) OneCond() *SimulateBalance {
	s.clear()

	s.tradeConfig = s.conds[0]
	s.quantity = s.tradeConfig.Quantity
	s.allowTradeTimeRange = s.tradePeriod.ToTimeRange(s.tradeConfig.TradeTimeRange.FirstPartDuration, s.tradeConfig.TradeTimeRange.SecondPartDuration)

	s.process()
	return s.calculateFutureTradeBalance()
}

func (s *SimulatorFuture) AllConds(slackMsgChan chan string) {
	slackMsgChan <- fmt.Sprintf("Start SimulatorFuture AllConds: %d", len(s.conds))
	var best int64
	for _, cond := range s.conds {
		s.clear()

		s.tradeConfig = cond
		s.quantity = s.tradeConfig.Quantity
		s.allowTradeTimeRange = s.tradePeriod.ToTimeRange(s.tradeConfig.TradeTimeRange.FirstPartDuration, s.tradeConfig.TradeTimeRange.SecondPartDuration)

		s.process()
		result := s.calculateFutureTradeBalance()
		if result.TotalBalance > best {
			best = result.TotalBalance
			slackMsgChan <- s.calculateFutureTradeBalance().String()
			slackMsgChan <- "------------------------------------"
		}
	}
	slackMsgChan <- "SimulatorFuture AllConds Done"
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

func (s *SimulatorFuture) process() {
	stuck := make(chan struct{})
	go func() {
		for {
			tick, ok := <-s.tickChan
			if !ok {
				break
			}

			s.tickArr = append(s.tickArr, tick)
			s.cutTickArr()

			if s.waitingOrder != nil {
				s.checkByBalance(tick)
				continue
			}

			s.lastTick = tick
			if o := s.generateOrder(); o != nil {
				o.OrderTime = tick.TickTime
				o.OrderID = uuid.NewString()
				o.Status = entity.StatusFilled

				s.lastPlaceOrderTime = tick.TickTime
				s.maxTradeOutTime = tick.TickTime.Add(time.Duration(s.tradeConfig.MaxHoldTime) * time.Minute)
				s.waitTimes[o.OrderID] = int(s.tradeConfig.TradeOutWaitTimes)

				s.allOrder = append(s.allOrder, o)
				s.waitingOrder = o
			}
		}
		close(stuck)
	}()

	for _, tick := range s.historyTickArr {
		s.tickChan <- tick
	}

	close(s.tickChan)
	<-stuck
}

func (s *SimulatorFuture) cutTickArr() {
	if len(s.tickArr) < 2 {
		return
	}

	if s.tickArr.GetLastTwoTickGapTime() > 3*time.Second {
		s.tickArr = entity.RealTimeFutureTickArr{}
		s.lastTickRate = 0
		return
	}

	if s.tickArr.GetTotalTime() > time.Duration(2*s.tradeConfig.TickInterval)*time.Second {
		s.tickArr = s.tickArr[1:]
	}
}

func (s *SimulatorFuture) generateOrder() *entity.FutureOrder {
	if s.lastTick.TickTime.Sub(s.lastPlaceOrderTime) < 3*time.Minute || !s.isAllowTrade(s.lastTick.TickTime) {
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
	case 100-outInRatio > s.tradeConfig.InOutRatio:
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

	if times := s.waitTimes[s.waitingOrder.OrderID]; times <= 0 {
		return false
	}

	switch s.waitingOrder.Action {
	case entity.ActionBuy:
		if tick.Close < s.lastTick.Close {
			return false
		}
		s.waitTimes[s.waitingOrder.OrderID]--

	case entity.ActionSell:
		if tick.Close > s.lastTick.Close {
			return false
		}
		s.waitTimes[s.waitingOrder.OrderID]--
	}
	return true
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

	if !place && tick.TickTime.Before(s.maxTradeOutTime) {
		return
	}

	o := &entity.FutureOrder{
		Code: tick.Code,
		BaseOrder: entity.BaseOrder{
			Price:     tick.Close,
			Quantity:  s.waitingOrder.Quantity,
			OrderTime: tick.TickTime,
			OrderID:   uuid.NewString(),
			Status:    entity.StatusFilled,
		},
	}

	switch s.waitingOrder.Action {
	case entity.ActionSell:
		o.Action = entity.ActionBuy
	case entity.ActionBuy:
		o.Action = entity.ActionSell
	}

	s.allOrder = append(s.allOrder, o)
	s.waitingOrder = nil
}

func (s *SimulatorFuture) calculateFutureTradeBalance() *SimulateBalance {
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

	forwardBalance, fCount := s.calculateForwardFutureBalance(forward)
	revereBalance, rCount := s.calculateReverseFutureBalance(reverse)

	return &SimulateBalance{
		TotalBalance: forwardBalance + revereBalance,
		Forward:      forwardBalance,
		ForwardCount: fCount,
		ForwardOrder: forward,
		Reverse:      revereBalance,
		ReverseCount: rCount,
		ReverseOrder: reverse,
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
		logger.Error("forward qty not zero")
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
		logger.Error("forward qty not zero")
		return 0, tradeCount
	}

	return reverseBalance, tradeCount
}
