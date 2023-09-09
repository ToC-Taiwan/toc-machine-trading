package hadger

import (
	"tmt/internal/entity"
	"tmt/internal/usecase/grpc"
	"tmt/pkg/eventbus"

	"github.com/google/uuid"
)

type HadgeTraderStock struct {
	id string

	tickChan chan *entity.RealTimeStockTick
	tickArr  []*entity.RealTimeStockTick

	forwardTrader *grpc.TradegRPCAPI
	reverseTrader *grpc.TradegRPCAPI

	bus *eventbus.Bus
}

func NewHadgeTraderStock(s, f *grpc.TradegRPCAPI, bus *eventbus.Bus) *HadgeTraderStock {
	h := &HadgeTraderStock{
		id:            uuid.NewString(),
		tickChan:      make(chan *entity.RealTimeStockTick),
		tickArr:       []*entity.RealTimeStockTick{},
		forwardTrader: s,
		reverseTrader: f,
		bus:           bus,
	}

	h.bus.SubscribeAsync(topicUpdateOrder, true, h.updateOrder)
	h.prepare()
	return h
}

func (h *HadgeTraderStock) TickChan() chan *entity.RealTimeStockTick {
	return h.tickChan
}

func (h *HadgeTraderStock) prepare() {
	go func() {
		for {
			tick := <-h.tickChan
			h.tickArr = append(h.tickArr, tick)
			// o := h.generateOrder()
		}
	}()
}

func (h *HadgeTraderStock) updateOrder(order *entity.StockOrder) {
	// TODO:
}
