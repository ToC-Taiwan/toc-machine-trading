package dt

import (
	"tmt/internal/entity"
	"tmt/internal/usecase/grpcapi"
	"tmt/pkg/eventbus"

	"github.com/google/uuid"
)

type DTTraderFuture struct {
	id string

	tickChan chan *entity.RealTimeFutureTick
	tickArr  []*entity.RealTimeFutureTick

	sc *grpcapi.TradegRPCAPI

	bus *eventbus.Bus
}

func NewDTTraderFuture(s *grpcapi.TradegRPCAPI, bus *eventbus.Bus) *DTTraderFuture {
	h := &DTTraderFuture{
		id:       uuid.NewString(),
		tickChan: make(chan *entity.RealTimeFutureTick),
		tickArr:  []*entity.RealTimeFutureTick{},
		sc:       s,
		bus:      bus,
	}

	h.bus.SubscribeTopic(topicUpdateOrder, h.updateOrder)
	h.prepare()
	return h
}

func (h *DTTraderFuture) TickChan() chan *entity.RealTimeFutureTick {
	return h.tickChan
}

func (h *DTTraderFuture) prepare() {
	go func() {
		for {
			tick := <-h.tickChan
			h.tickArr = append(h.tickArr, tick)
			// o := h.generateOrder()
		}
	}()
}

func (h *DTTraderFuture) updateOrder(order *entity.FutureOrder) {
	// TODO:
}
