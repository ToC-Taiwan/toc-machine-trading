// Package dt package dt
package dt

import (
	"sync"

	"tmt/cmd/config"
	"tmt/internal/entity"
	"tmt/internal/usecase/grpcapi"
	"tmt/pkg/eventbus"

	"github.com/google/uuid"
)

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
		orderQuantity: 1,
		tickChan:      make(chan *entity.RealTimeFutureTick),
		notify:        make(chan *entity.FutureOrder),
		sc:            s,
		tickArr:       []*entity.RealTimeFutureTick{},
		localBus:      eventbus.Get(uuid.NewString()),
		tradeConfig:   tradeConfig,
		traderMap:     make(map[string]*DTTraderFuture),
	}

	d.localBus.SubscribeTopic(topicTraderDone, d.removeDoneTrader)
	d.prepare()
	d.processOrderStatus()

	return d
}

func (d *DTFuture) prepare() {
	go func() {
		for {
			tick := <-d.tickChan
			d.tickArr = append(d.tickArr, tick)
			if len(d.tickArr) > 1 {
				d.tickArr = d.tickArr[1:]
			}

			if d.check() {
				d.traderMapLock.Lock()
				for _, trader := range d.traderMap {
					trader.TickChan() <- tick
				}
				d.traderMapLock.Unlock()
			}
		}
	}()
}

func (d *DTFuture) check() bool {
	return true
}

func (d *DTFuture) processOrderStatus() {}

func (d *DTFuture) removeDoneTrader() {
	d.traderMapLock.Lock()
	defer d.traderMapLock.Unlock()
}

func (d *DTFuture) Notify() chan *entity.FutureOrder {
	return d.notify
}

func (d *DTFuture) TickChan() chan *entity.RealTimeFutureTick {
	return d.tickChan
}
