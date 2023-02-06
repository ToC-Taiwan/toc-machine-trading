// Package dt package dt
package dt

import (
	"sync"

	"tmt/cmd/config"
	"tmt/internal/entity"
	"tmt/internal/usecase/grpcapi"
	"tmt/pkg/eventbus"
	"tmt/pkg/log"

	"github.com/google/uuid"
)

var logger = log.Get()

type DTFuture struct {
	code          string
	orderQuantity int64
	tickArr       []*entity.RealTimeFutureTick

	sc       *grpcapi.TradegRPCAPI
	localBus *eventbus.Bus

	traderMap     map[string]*DTTraderFuture
	traderMapLock sync.RWMutex

	tickChan chan *entity.RealTimeFutureTick
	notify   chan *entity.FutureOrder

	tradeConfig *config.TradeFuture
}

func NewDTFuture(code string, s *grpcapi.TradegRPCAPI, tradeConfig *config.TradeFuture) *DTFuture {
	d := &DTFuture{
		code:          code,
		orderQuantity: tradeConfig.Quantity,
		tickChan:      make(chan *entity.RealTimeFutureTick),
		notify:        make(chan *entity.FutureOrder),
		sc:            s,
		tickArr:       []*entity.RealTimeFutureTick{},
		localBus:      eventbus.Get(uuid.NewString()),
		tradeConfig:   tradeConfig,
		traderMap:     make(map[string]*DTTraderFuture),
	}

	d.localBus.SubscribeTopic(topicTraderDone, d.removeDoneTrader)

	d.processTick()
	d.processOrderStatus()

	return d
}

func (d *DTFuture) Notify() chan *entity.FutureOrder {
	return d.notify
}

func (d *DTFuture) TickChan() chan *entity.RealTimeFutureTick {
	return d.tickChan
}

func (d *DTFuture) processOrderStatus() {
	finishedOrderMap := make(map[string]*entity.FutureOrder)
	go func() {
		for {
			o := <-d.notify
			if _, ok := finishedOrderMap[o.OrderID]; ok {
				continue
			}

			if !o.Cancellable() {
				finishedOrderMap[o.OrderID] = o
			} else {
				d.cancelOverTimeOrder(o)
			}

			d.localBus.PublishTopicEvent(topicUpdateOrder, o)
		}
	}()
}

func (d *DTFuture) processTick() {
	go func() {
		for {
			tick := <-d.tickChan
			d.tickArr = append(d.tickArr, tick)

			if len(d.tickArr) > 1 {
				d.tickArr = d.tickArr[1:]
			}

			if o := d.generateOrder(); o != nil {
				for i := 0; i < int(d.orderQuantity); i++ {
					d.addTrader(o)
				}
			}

			d.sendTickToTrader(tick)
		}
	}()
}

func (d *DTFuture) sendTickToTrader(tick *entity.RealTimeFutureTick) {
	d.traderMapLock.RLock()
	defer d.traderMapLock.RUnlock()

	for _, trader := range d.traderMap {
		if ch := trader.TickChan(); ch != nil {
			ch <- tick
		}
	}
}

func (d *DTFuture) cancelOverTimeOrder(order *entity.FutureOrder) {
	result, err := d.sc.CancelFuture(order.OrderID)
	if err != nil {
		logger.Error(err)
		return
	}

	if entity.StringToOrderStatus(result.GetStatus()) != entity.StatusCancelled {
		logger.Error("Cancel future order failed", result.GetStatus())
		return
	}
}

func (d *DTFuture) generateOrder() *entity.FutureOrder {
	return &entity.FutureOrder{}
}

func (d *DTFuture) addTrader(order *entity.FutureOrder) {
	d.traderMapLock.Lock()
	defer d.traderMapLock.Unlock()

	orderWithCfg := orderWithCfg{
		order: order,
		cfg:   d.tradeConfig,
	}

	if trader := NewDTTraderFuture(orderWithCfg, d.sc, d.localBus); trader != nil {
		d.traderMap[trader.id] = trader
	}
}

func (d *DTFuture) removeDoneTrader(id string) {
	d.traderMapLock.Lock()
	defer d.traderMapLock.Unlock()

	if trader := d.traderMap[id]; trader != nil {
		close(trader.tickChan)
		delete(d.traderMap, id)
	}
}
