// Package dt package dt
package dt

import (
	"sync"
	"time"

	"tmt/cmd/config"
	"tmt/internal/entity"
	"tmt/internal/usecase/grpcapi"
	"tmt/pkg/eventbus"
	"tmt/pkg/log"

	"github.com/google/uuid"
)

type DTFuture struct {
	code          string
	orderQuantity int64

	lastTick *entity.RealTimeFutureTick
	tickArr  entity.RealTimeFutureTickArr

	sc          *grpcapi.TradegRPCAPI
	localBus    *eventbus.Bus
	tradeConfig *config.TradeFuture

	traderMap     map[string]*DTTraderFuture
	traderMapLock sync.RWMutex

	tickChan   chan *entity.RealTimeFutureTick
	notify     chan *entity.FutureOrder
	cancelChan chan *entity.FutureOrder
	switchChan chan bool

	lastTickRate float64
	isTradeTime  bool

	lastPlaceOrderTime time.Time
	logger             *log.Log
}

func NewDTFuture(code string, s *grpcapi.TradegRPCAPI, tradeConfig *config.TradeFuture) *DTFuture {
	d := &DTFuture{
		code:          code,
		orderQuantity: tradeConfig.Quantity,
		tickChan:      make(chan *entity.RealTimeFutureTick),
		notify:        make(chan *entity.FutureOrder),
		switchChan:    make(chan bool),
		sc:            s,
		tickArr:       []*entity.RealTimeFutureTick{},
		localBus:      eventbus.Get(uuid.NewString()),
		tradeConfig:   tradeConfig,
		traderMap:     make(map[string]*DTTraderFuture),
		cancelChan:    make(chan *entity.FutureOrder),
		logger:        log.Get(),
	}

	d.localBus.SubscribeAsync(topicTraderDone, true, d.removeDoneTrader)

	d.cancelOverTimeOrder()
	d.processOrderStatusAndTradeSwitch()
	d.processTick()

	return d
}

func (d *DTFuture) cancelOverTimeOrder() {
	cancelledIDMap := make(map[string]*entity.FutureOrder)
	go func() {
		for {
			order := <-d.cancelChan
			if _, ok := cancelledIDMap[order.OrderID]; ok {
				continue
			}

			result, err := d.sc.CancelFuture(order.OrderID)
			if err != nil {
				d.logger.Error(err)
				continue
			}

			if s := entity.StringToOrderStatus(result.GetStatus()); s != entity.StatusCancelled {
				d.logger.Errorf("Cancel order failed: %s %s", order.OrderID, result.GetError())
				continue
			}

			cancelledIDMap[order.OrderID] = order
		}
	}()
}

func (d *DTFuture) processOrderStatusAndTradeSwitch() {
	notifiedMap := make(map[string]*entity.FutureOrder)
	go func() {
		for {
			select {
			case o := <-d.notify:
				if !d.isTradeTime {
					continue
				}

				switch {
				case !o.Cancellable() && notifiedMap[o.OrderID] == nil:
					notifiedMap[o.OrderID] = o
					d.logger.Warnf("%s %s", o.Status.String(), o.String())
				case o.Cancellable() && time.Since(o.OrderTime) > time.Duration(d.tradeConfig.BuySellWaitTime)*time.Second:
					d.cancelChan <- o
				}

				d.localBus.PublishTopicEvent(topicUpdateOrder, o)

			case ts := <-d.switchChan:
				d.isTradeTime = ts
			}
		}
	}()
}

func (d *DTFuture) processTick() {
	go func() {
		for {
			d.lastTick = <-d.tickChan
			d.tickArr = append(d.tickArr, d.lastTick)

			d.cutTickArr()
			if d.sendTickToTrader(d.lastTick) > 0 {
				continue
			}

			if o := d.generateOrder(); o != nil {
				d.lastPlaceOrderTime = d.lastTick.TickTime
				var wg sync.WaitGroup
				for i := 0; i < int(d.orderQuantity); i++ {
					wg.Add(1)
					order := *o
					go d.addTrader(&order, &wg)
				}
				wg.Wait()
			}
		}
	}()
}

func (d *DTFuture) cutTickArr() {
	if len(d.tickArr) < 2 {
		return
	}

	if d.tickArr.GetLastTwoTickGapTime() > time.Second {
		d.tickArr = entity.RealTimeFutureTickArr{}
		d.lastTickRate = 0
		return
	}

	if d.tickArr.GetTotalTime() > time.Duration(2*d.tradeConfig.TickInterval)*time.Second {
		d.tickArr = d.tickArr[1:]
	}
}

func (d *DTFuture) sendTickToTrader(tick *entity.RealTimeFutureTick) int {
	d.traderMapLock.RLock()
	defer d.traderMapLock.RUnlock()
	var sent int
	for _, trader := range d.traderMap {
		sent++
		if ch := trader.TickChan(); ch != nil {
			ch <- tick
		}
	}
	return sent
}

func (d *DTFuture) generateOrder() *entity.FutureOrder {
	if !d.tradeConfig.AllowTrade || !d.isTradeTime {
		return nil
	}

	if time.Since(d.lastPlaceOrderTime) < 3*time.Minute {
		return nil
	}

	outInRatio, tickRate := d.tickArr.GetOutInRatioAndRate(time.Duration(d.tradeConfig.TickInterval) * time.Second)
	defer func() {
		d.lastTickRate = tickRate
	}()
	if d.lastTickRate == 0 {
		return nil
	}

	if d.lastTickRate < d.tradeConfig.RateLimit || 100*(tickRate-d.lastTickRate)/d.lastTickRate < d.tradeConfig.RateChangeRatio {
		return nil
	}

	switch {
	case outInRatio > d.tradeConfig.OutInRatio:
		return &entity.FutureOrder{
			Code: d.code,
			BaseOrder: entity.BaseOrder{
				Action:   entity.ActionBuy,
				Price:    d.lastTick.Close - 1,
				Quantity: orderQtyUnit,
			},
		}
	case 100-outInRatio > d.tradeConfig.InOutRatio:
		return &entity.FutureOrder{
			Code: d.code,
			BaseOrder: entity.BaseOrder{
				Action:   entity.ActionSell,
				Price:    d.lastTick.Close + 1,
				Quantity: orderQtyUnit,
			},
		}
	default:
		return nil
	}
}

func (d *DTFuture) addTrader(order *entity.FutureOrder, wg *sync.WaitGroup) {
	defer wg.Done()

	orderWithCfg := orderWithCfg{
		order:           order,
		cfg:             d.tradeConfig,
		maxTradeOutTime: time.Now().Add(time.Duration(d.tradeConfig.MaxHoldTime) * time.Minute),
	}

	if trader := NewDTTraderFuture(orderWithCfg, d.sc, d.localBus); trader != nil {
		d.traderMapLock.Lock()
		d.traderMap[trader.id] = trader
		d.traderMapLock.Unlock()
	}
}

func (d *DTFuture) removeDoneTrader(id string) {
	if trader := d.traderMap[id]; trader != nil {
		d.traderMapLock.Lock()
		close(trader.tickChan)
		delete(d.traderMap, id)
		d.traderMapLock.Unlock()
	}
}

func (d *DTFuture) Notify() chan *entity.FutureOrder {
	return d.notify
}

func (d *DTFuture) TickChan() chan *entity.RealTimeFutureTick {
	return d.tickChan
}

func (d *DTFuture) SwitchChan() chan bool {
	return d.switchChan
}
