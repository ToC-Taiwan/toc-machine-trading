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

type SimulatorFuture struct {
	// current cond
	currentCond *config.TradeFuture

	// given by target
	code           string
	quota          *quota.Quota
	tradePeriod    tradeday.TradePeriod
	historyTickArr entity.RealTimeFutureTickArr
	allCond        []*config.TradeFuture

	// given by cond
	quantity int64

	// clear before each cond
	waitTimes           map[string]int
	tickChan            chan *entity.RealTimeFutureTick
	allowTradeTimeRange [][]time.Time
	allOrder            []*entity.FutureOrder
	lastTickRate        float64
	waitingOrder        *entity.FutureOrder
	lastTick            *entity.RealTimeFutureTick
	tickArr             entity.RealTimeFutureTickArr
	lastPlaceOrderTime  time.Time
	maxTradeOutTime     time.Time

	logger *log.Log
}

func NewSimulatorFuture(target SimulatorFutureTarget) *SimulatorFuture {
	logger := log.Get()
	if len(target.TradeConfigArr) == 0 {
		logger.Fatal("trade config arr is empty")
	}

	return &SimulatorFuture{
		code:           target.Code,
		quota:          target.Quota,
		tradePeriod:    target.TradePeriod,
		historyTickArr: target.Ticks,
		allCond:        target.TradeConfigArr,
		logger:         logger,
	}
}

func (s *SimulatorFuture) clear() {
	s.waitTimes = make(map[string]int)
	s.tickChan = make(chan *entity.RealTimeFutureTick)
	s.allowTradeTimeRange = [][]time.Time{}
	s.allOrder = []*entity.FutureOrder{}
	s.lastTickRate = 0
	s.waitingOrder = nil
	s.lastTick = nil
	s.tickArr = entity.RealTimeFutureTickArr{}
	s.lastPlaceOrderTime = time.Time{}
	s.maxTradeOutTime = time.Time{}
}

func (s *SimulatorFuture) OneCond() *SimulateBalance {
	s.clear()

	s.currentCond = s.allCond[0]
	s.quantity = s.currentCond.Quantity
	s.allowTradeTimeRange = s.tradePeriod.ToTimeRange(s.currentCond.TradeTimeRange.FirstPartDuration, s.currentCond.TradeTimeRange.SecondPartDuration)

	return s.process().calculateFutureTradeBalance()
}

func (s *SimulatorFuture) AllConds(slackMsgChan chan string) {
	slackMsgChan <- fmt.Sprintf("Start SimulatorFuture AllConds: %d", len(s.allCond))
	workerCount := 10
	workQueue := make(chan *SimulatorFuture, workerCount)
	for i := 0; i < workerCount; i++ {
		clone := *s
		workQueue <- &clone
	}

	var best int64
	resultChan := make(chan *SimulateBalance)
	go func() {
		for {
			result, ok := <-resultChan
			if !ok {
				break
			}
			if result.TotalBalance > best {
				best = result.TotalBalance
				slackMsgChan <- result.String()
				slackMsgChan <- "------------------------------------"
			}
		}
	}()

	for _, cond := range s.allCond {
		c := *cond
		worker := <-workQueue
		go func() {
			worker.clear()
			worker.currentCond = &c
			worker.quantity = worker.currentCond.Quantity
			worker.allowTradeTimeRange = worker.tradePeriod.ToTimeRange(worker.currentCond.TradeTimeRange.FirstPartDuration, worker.currentCond.TradeTimeRange.SecondPartDuration)
			resultChan <- worker.process().calculateFutureTradeBalance()
			workQueue <- worker
		}()
	}

	recoverWorker := []*SimulatorFuture{}
	for {
		worker := <-workQueue
		recoverWorker = append(recoverWorker, worker)
		if len(recoverWorker) == workerCount {
			close(resultChan)
			close(workQueue)
			break
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

func (s *SimulatorFuture) process() *SimulatorFuture {
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
				s.maxTradeOutTime = tick.TickTime.Add(time.Duration(s.currentCond.MaxHoldTime) * time.Minute)
				s.waitTimes[o.OrderID] = int(s.currentCond.TradeOutWaitTimes)

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
	return s
}

func (s *SimulatorFuture) cutTickArr() {
	if len(s.tickArr) < 2 {
		return
	}

	if s.tickArr.GetLastTwoTickGapTime() > time.Second {
		s.tickArr = entity.RealTimeFutureTickArr{}
		s.lastTickRate = 0
		return
	}

	if s.tickArr.GetTotalTime() > time.Duration(2*s.currentCond.TickInterval)*time.Second {
		s.tickArr = s.tickArr[1:]
	}
}

func (s *SimulatorFuture) generateOrder() *entity.FutureOrder {
	if s.lastTick.TickTime.Sub(s.lastPlaceOrderTime) < 3*time.Minute || !s.isAllowTrade(s.lastTick.TickTime) {
		return nil
	}

	outInRatio, tickRate := s.tickArr.GetOutInRatioAndRate(time.Duration(s.currentCond.TickInterval) * time.Second)
	defer func() {
		s.lastTickRate = tickRate
	}()
	if s.lastTickRate == 0 {
		return nil
	}

	if s.lastTickRate < s.currentCond.RateLimit || 100*(tickRate-s.lastTickRate)/s.lastTickRate < s.currentCond.RateChangeRatio {
		return nil
	}

	switch {
	case outInRatio > s.currentCond.OutInRatio:
		return &entity.FutureOrder{
			Code: s.code,
			BaseOrder: entity.BaseOrder{
				Action:   entity.ActionBuy,
				Price:    s.lastTick.Close - 1,
				Quantity: s.currentCond.Quantity,
			},
		}
	case 100-outInRatio > s.currentCond.InOutRatio:
		return &entity.FutureOrder{
			Code: s.code,
			BaseOrder: entity.BaseOrder{
				Action:   entity.ActionSell,
				Price:    s.lastTick.Close + 1,
				Quantity: s.currentCond.Quantity,
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
		if tick.Close <= s.waitingOrder.Price-s.currentCond.TargetBalanceHigh || tick.Close >= s.waitingOrder.Price-s.currentCond.TargetBalanceLow {
			place = true
		}

	case entity.ActionBuy:
		if tick.Close >= s.waitingOrder.Price+s.currentCond.TargetBalanceHigh || tick.Close <= s.waitingOrder.Price+s.currentCond.TargetBalanceLow {
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
		Cond:         s.currentCond,
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
		s.logger.Error("forward qty not zero")
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
		s.logger.Error("forward qty not zero")
		return 0, tradeCount
	}

	return reverseBalance, tradeCount
}
