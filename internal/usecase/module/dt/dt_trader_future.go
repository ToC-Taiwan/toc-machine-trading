package dt

import (
	"sync"

	"tmt/cmd/config"
	"tmt/internal/entity"
	"tmt/internal/usecase/grpcapi"
	"tmt/pkg/eventbus"

	"github.com/google/uuid"
)

type DTTraderFuture struct {
	id string

	tickChan chan *entity.RealTimeFutureTick

	sc  *grpcapi.TradegRPCAPI
	bus *eventbus.Bus

	tradeConfig  *config.TradeFuture
	baseOrder    *entity.FutureOrder
	waitingOrder *entity.FutureOrder

	finishOrderMap     map[string]*entity.FutureOrder
	finishOrderMapLock sync.RWMutex

	ready bool
	done  bool

	once sync.Once
}

func NewDTTraderFuture(orderWithCfg orderWithCfg, s *grpcapi.TradegRPCAPI, bus *eventbus.Bus) *DTTraderFuture {
	d := &DTTraderFuture{
		id:          uuid.NewString(),
		tickChan:    make(chan *entity.RealTimeFutureTick),
		sc:          s,
		bus:         bus,
		baseOrder:   orderWithCfg.order,
		tradeConfig: orderWithCfg.cfg,
	}

	if err := d.placeOrder(orderWithCfg.order); err != nil {
		logger.Errorf("NewDTTraderFuture err: %s", err.Error())
		return nil
	}

	d.bus.SubscribeTopic(topicUpdateOrder, d.updateOrder)
	d.processTick()

	return d
}

func (d *DTTraderFuture) processTick() {
	var tradeOutAction entity.OrderAction
	switch d.baseOrder.Action {
	case entity.ActionBuy:
		tradeOutAction = entity.ActionSell
	case entity.ActionSell:
		tradeOutAction = entity.ActionBuy
	}

	go func() {
		for {
			tick, ok := <-d.tickChan
			if !ok {
				return
			}

			if d.waitingOrder != nil {
				continue
			}

			if d.isTraderDone() {
				continue
			}

			d.checkByBalance(tick, tradeOutAction)
		}
	}()
}

func (d *DTTraderFuture) checkByBalance(tick *entity.RealTimeFutureTick, tradeOutAction entity.OrderAction) {
	var place bool
	switch tradeOutAction {
	case entity.ActionSell:
		if tick.Close >= d.baseOrder.Price+d.tradeConfig.TargetBalanceHigh || tick.Close <= d.baseOrder.Price+d.tradeConfig.TargetBalanceLow {
			place = true
		}

	case entity.ActionBuy:
		if tick.Close <= d.baseOrder.Price-d.tradeConfig.TargetBalanceHigh || tick.Close >= d.baseOrder.Price-d.tradeConfig.TargetBalanceLow {
			place = true
		}
	}

	if !place {
		return
	}

	o := &entity.FutureOrder{
		Code: tick.Code,
		BaseOrder: entity.BaseOrder{
			Action:   tradeOutAction,
			Price:    tick.Close,
			Quantity: d.baseOrder.Quantity,
		},
	}

	if err := d.placeOrder(o); err != nil {
		logger.Errorf("checkByBalance place order error: %s", err.Error())
		return
	}
	d.waitingOrder = o
}

func (d *DTTraderFuture) updateOrder(order *entity.FutureOrder) {
	if !d.ready && d.baseOrder.OrderID == order.OrderID {
		switch order.Status {
		case entity.StatusFilled:
			d.ready = true
		case entity.StatusCancelled, entity.StatusFailed:
			d.postDone()
		}
		return
	}

	d.finishOrderMapLock.Lock()
	defer d.finishOrderMapLock.Unlock()

	if _, ok := d.finishOrderMap[order.OrderID]; ok {
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

	if d.ready {
		d.finishOrderMapLock.Lock()
		d.finishOrderMap[o.OrderID] = o
		d.finishOrderMapLock.Unlock()
	}

	return nil
}

func (d *DTTraderFuture) sell(o *entity.FutureOrder) error {
	result, err := d.sc.SellFuture(o)
	if err != nil {
		return err
	}

	o.OrderID = result.GetOrderId()
	o.Status = entity.StringToOrderStatus(result.GetStatus())

	if d.ready {
		d.finishOrderMapLock.Lock()
		d.finishOrderMap[o.OrderID] = o
		d.finishOrderMapLock.Unlock()
	}

	return nil
}

func (d *DTTraderFuture) TickChan() chan *entity.RealTimeFutureTick {
	if !d.ready {
		return nil
	}
	return d.tickChan
}

func (d *DTTraderFuture) postDone() {
	d.once.Do(func() {
		d.bus.UnSubscribeTopic(topicUpdateOrder, d.updateOrder)
		d.bus.PublishTopicEvent(topicTraderDone, d.id)
		d.done = true
	})
}
