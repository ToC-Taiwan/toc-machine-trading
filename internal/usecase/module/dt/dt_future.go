// Package dt package dt
package dt

import (
	"sync"
	"time"

	"tmt/cmd/config"
	"tmt/internal/entity"
	"tmt/internal/usecase/grpcapi"
	"tmt/pkg/eventbus"

	"github.com/google/uuid"
)

type DTFuture struct {
	code          string
	orderQuantity int64
	tickArr       entity.RealTimeFutureTickArr

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

	d.processOrderStatus()
	d.processTick()

	return d
}

func (d *DTFuture) processOrderStatus() {
	notCancellableOrderMap := make(map[string]*entity.FutureOrder)
	go func() {
		for {
			o := <-d.notify
			if _, ok := notCancellableOrderMap[o.OrderID]; ok {
				continue
			}

			if !o.Cancellable() {
				notCancellableOrderMap[o.OrderID] = o
			} else {
				d.cancelOverTimeOrder(o)
			}

			d.localBus.PublishTopicEvent(topicUpdateOrder, o)
		}
	}()
}

func (d *DTFuture) cancelOverTimeOrder(order *entity.FutureOrder) {
	if time.Since(order.OrderTime) < time.Duration(d.tradeConfig.BuySellWaitTime)*time.Second {
		return
	}

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
					go d.addTrader(o)
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

func (d *DTFuture) generateOrder() *entity.FutureOrder {
	d.traderMapLock.RLock()
	defer d.traderMapLock.RUnlock()
	if len(d.traderMap) > 0 {
		return nil
	}

	return &entity.FutureOrder{}
}

func (d *DTFuture) addTrader(order *entity.FutureOrder) {
	orderWithCfg := orderWithCfg{
		order: order,
		cfg:   d.tradeConfig,
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
