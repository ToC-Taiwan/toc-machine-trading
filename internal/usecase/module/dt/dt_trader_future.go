package dt

import (
	"sync"

	"tmt/internal/entity"
	"tmt/internal/usecase/grpcapi"
	"tmt/pkg/eventbus"

	"github.com/google/uuid"
)

type DTTraderFuture struct {
	id string

	tickChan chan *entity.RealTimeFutureTick
	tickArr  []*entity.RealTimeFutureTick

	sc  *grpcapi.TradegRPCAPI
	bus *eventbus.Bus

	baseOrder    *entity.FutureOrder
	waitingOrder *entity.FutureOrder

	finishOrderMap     map[string]*entity.FutureOrder
	finishOrderMapLock sync.RWMutex

	done     bool
	needTick bool
}

func NewDTTraderFuture(order *entity.FutureOrder, s *grpcapi.TradegRPCAPI, bus *eventbus.Bus) *DTTraderFuture {
	d := &DTTraderFuture{
		id:        uuid.NewString(),
		tickChan:  make(chan *entity.RealTimeFutureTick),
		tickArr:   []*entity.RealTimeFutureTick{},
		sc:        s,
		bus:       bus,
		baseOrder: order,
	}

	if err := d.placeOrder(order); err != nil {
		logger.Error(err)
		return nil
	}

	d.bus.SubscribeTopic(topicUpdateOrder, d.updateOrder)
	d.processTick()

	return d
}

func (d *DTTraderFuture) processTick() {
	go func() {
		for {
			tick, ok := <-d.tickChan
			if !ok {
				return
			}

			d.tickArr = append(d.tickArr, tick)

			if d.waitingOrder != nil {
				continue
			}

			if d.isTraderDone() {
				continue
			}
		}
	}()
}

func (d *DTTraderFuture) updateOrder(order *entity.FutureOrder) {
	d.finishOrderMapLock.Lock()
	defer d.finishOrderMapLock.Unlock()

	if _, ok := d.finishOrderMap[order.OrderID]; ok {
		if d.baseOrder.OrderID == order.OrderID && order.Status == entity.StatusFilled {
			d.done = true
		}

		d.finishOrderMap[order.OrderID] = order
		if d.waitingOrder != nil && !order.Cancellable() && d.waitingOrder.OrderID == order.OrderID {
			d.waitingOrder = nil
		}
	}
}

func (d *DTTraderFuture) isTraderDone() bool {
	if d.done {
		return true
	}

	var endQty int64
	d.finishOrderMapLock.RLock()
	for _, o := range d.finishOrderMap {
		if o.OrderID == d.baseOrder.OrderID {
			continue
		}

		if o.Status == entity.StatusFilled {
			endQty += o.Quantity
		}
	}
	d.finishOrderMapLock.RUnlock()

	if endQty == d.baseOrder.Quantity {
		d.postDone()
		return true
	}
	return false
}

func (d *DTTraderFuture) placeOrder(o *entity.FutureOrder) error {
	switch o.Action {
	case entity.ActionBuy:
		return d.buy(o)
	case entity.ActionSell:
		return d.sell(o)
	default:
		return nil
	}
}

func (d *DTTraderFuture) buy(o *entity.FutureOrder) error {
	result, err := d.sc.BuyFuture(o)
	if err != nil {
		return err
	}

	o.OrderID = result.GetOrderId()
	o.Status = entity.StringToOrderStatus(result.GetStatus())

	d.finishOrderMapLock.Lock()
	d.finishOrderMap[o.OrderID] = o
	d.finishOrderMapLock.Unlock()

	return nil
}

func (d *DTTraderFuture) sell(o *entity.FutureOrder) error {
	result, err := d.sc.SellFuture(o)
	if err != nil {
		return err
	}

	o.OrderID = result.GetOrderId()
	o.Status = entity.StringToOrderStatus(result.GetStatus())

	d.finishOrderMapLock.Lock()
	d.finishOrderMap[o.OrderID] = o
	d.finishOrderMapLock.Unlock()

	return nil
}

func (d *DTTraderFuture) TickChan() chan *entity.RealTimeFutureTick {
	if !d.needTick {
		return nil
	}
	return d.tickChan
}

func (d *DTTraderFuture) postDone() {
	d.bus.UnSubscribeTopic(topicUpdateOrder, d.updateOrder)
	d.bus.PublishTopicEvent(topicTraderDone, d.id)
	d.done = true
}
