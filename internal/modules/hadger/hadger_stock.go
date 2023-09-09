// Package hadger package hadger
package hadger

import (
	"sync"

	"tmt/cmd/config"
	"tmt/internal/entity"
	"tmt/internal/modules/quota"
	"tmt/internal/usecase/grpc"
	"tmt/pkg/eventbus"

	"github.com/google/uuid"
)

type HadgerStock struct {
	num           string
	orderQuantity int64
	tickArr       []*entity.RealTimeStockTick

	forwardTrader *grpc.TradegRPCAPI
	reverseTrader *grpc.TradegRPCAPI
	localBus      *eventbus.Bus
	quota         *quota.Quota

	traderMap     map[string]*HadgeTraderStock
	traderMapLock sync.RWMutex

	tickChan   chan *entity.RealTimeStockTick
	notify     chan *entity.StockOrder
	switchChan chan bool

	tradeConfig *config.TradeStock
	isTradeTime bool
}

func NewHadgerStock(num string, s, f *grpc.TradegRPCAPI, q *quota.Quota, tradeConfig *config.TradeStock) *HadgerStock {
	h := &HadgerStock{
		num:           num,
		orderQuantity: 1,
		tickChan:      make(chan *entity.RealTimeStockTick),
		notify:        make(chan *entity.StockOrder),
		forwardTrader: s,
		reverseTrader: f,
		quota:         q,
		tickArr:       []*entity.RealTimeStockTick{},
		localBus:      eventbus.New(uuid.NewString()),
		tradeConfig:   tradeConfig,
		traderMap:     make(map[string]*HadgeTraderStock),
		switchChan:    make(chan bool),
	}

	h.localBus.SubscribeAsync(topicTraderDone, true, h.removeDoneTrader)
	h.processTick()
	h.processOrderStatus()

	return h
}

func (h *HadgerStock) processTick() {
	go func() {
		for {
			tick := <-h.tickChan
			h.tickArr = append(h.tickArr, tick)

			if len(h.tickArr) > 1 {
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

func (h *HadgerStock) processOrderStatus() {
	finishedOrderMap := make(map[string]*entity.StockOrder)
	go func() {
		for {
			select {
			case o := <-h.notify:
				if _, ok := finishedOrderMap[o.OrderID]; ok {
					continue
				}

				if !o.Cancellable() {
					finishedOrderMap[o.OrderID] = o
				} else {
					h.cancelOverTimeOrder(o)
				}

				h.localBus.PublishTopicEvent(topicUpdateOrder, o)
			case ts := <-h.switchChan:
				h.isTradeTime = ts
			}
		}
	}()
}

func (h *HadgerStock) check() bool {
	return false
}

func (h *HadgerStock) sendTickToTrader(tick *entity.RealTimeStockTick) {
	h.traderMapLock.RLock()
	defer h.traderMapLock.RUnlock()

	for _, trader := range h.traderMap {
		trader.TickChan() <- tick
	}
}

func (h *HadgerStock) cancelOverTimeOrder(order *entity.StockOrder) {}

func (h *HadgerStock) removeDoneTrader(id string) {
	h.traderMapLock.Lock()
	defer h.traderMapLock.Unlock()

	delete(h.traderMap, id)
}

func (h *HadgerStock) TickChan() chan *entity.RealTimeStockTick {
	return h.tickChan
}

func (h *HadgerStock) Notify() chan *entity.StockOrder {
	return h.notify
}

func (h *HadgerStock) SwitchChan() chan bool {
	return h.switchChan
}
