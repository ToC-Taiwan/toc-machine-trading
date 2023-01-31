// Package hadger package hadger
package hadger

import (
	"sync"

	"tmt/internal/entity"
	"tmt/internal/usecase/grpcapi"
	"tmt/internal/usecase/module/quota"
	"tmt/pkg/eventbus"

	"github.com/google/uuid"
)

type HadgerStock struct {
	num           string
	orderQuantity int64
	tickArr       []*entity.RealTimeStockTick

	forwardTrader *grpcapi.TradegRPCAPI
	reverseTrader *grpcapi.TradegRPCAPI
	localBus      *eventbus.Bus
	quota         *quota.Quota

	traderMap     map[string]*HadgeTraderStock
	traderMapLock sync.RWMutex

	tickChan chan *entity.RealTimeStockTick
	notify   chan *entity.StockOrder
}

func NewHadgerStock(num string, s, f *grpcapi.TradegRPCAPI, q *quota.Quota) *HadgerStock {
	h := &HadgerStock{
		num:           num,
		orderQuantity: 1,
		tickChan:      make(chan *entity.RealTimeStockTick),
		notify:        make(chan *entity.StockOrder),
		forwardTrader: s,
		reverseTrader: f,
		quota:         q,
		tickArr:       []*entity.RealTimeStockTick{},
		localBus:      eventbus.Get(uuid.NewString()),
	}

	h.localBus.SubscribeTopic(topicTraderDone, h.removeDoneTrader)
	h.prepare()
	h.processOrderStatus()

	return h
}

func (h *HadgerStock) TickChan() chan *entity.RealTimeStockTick {
	return h.tickChan
}

func (h *HadgerStock) check() bool {
	return false
}

func (h *HadgerStock) Notify() chan *entity.StockOrder {
	return h.notify
}

func (h *HadgerStock) prepare() {
	go func() {
		for {
			tick := <-h.tickChan
			h.tickArr = append(h.tickArr, tick)

			if len(h.tickArr) > 10 {
				h.tickArr = h.tickArr[1:]
			}

			if h.check() {
				h.traderMapLock.Lock()
				f := NewHadgeTraderStock(h.forwardTrader, h.reverseTrader, h.localBus)
				r := NewHadgeTraderStock(h.reverseTrader, h.forwardTrader, h.localBus)
				h.traderMap[f.id] = f
				h.traderMap[r.id] = r
				h.traderMapLock.Unlock()
			}

			h.sendTickToTrader(tick)
		}
	}()
}

func (h *HadgerStock) removeDoneTrader(id string) {
	h.traderMapLock.Lock()
	defer h.traderMapLock.Unlock()

	delete(h.traderMap, id)
}

func (h *HadgerStock) sendTickToTrader(tick *entity.RealTimeStockTick) {
	h.traderMapLock.RLock()
	defer h.traderMapLock.RUnlock()

	for _, trader := range h.traderMap {
		trader.TickChan() <- tick
	}
}

func (h *HadgerStock) processOrderStatus() {
	// finishedOrderMap := make(map[string]*entity.FutureOrder)
	go func() {
		for {
			order, ok := <-h.notify
			if !ok {
				return
			}

			h.localBus.PublishTopicEvent(topicUpdateOrder, order)
			// if o, ok := order.(*entity.FutureOrder); ok {
			// 	if finishedOrderMap[o.OrderID] != nil {
			// 		continue
			// 	}

			// 	w.updateCacheOrder(o)
			// 	if !o.Cancellable() {
			// 		finishedOrderMap[o.OrderID] = o
			// 		if w.waitingList.orderIDExist(o.OrderID) {
			// 			w.waitingList.remove(o.OrderID)
			// 		}
			// 	} else {
			// 		w.cancelOverTimeOrder(o)
			// 	}
			// }
		}
	}()
}
