package dt

import (
	"errors"
	"sync"
	"time"

	"tmt/cmd/config"
	"tmt/internal/entity"
	"tmt/internal/usecase/grpcapi"
	"tmt/pb"
	"tmt/pkg/eventbus"

	"github.com/google/uuid"
)

type DTTraderFuture struct {
	id    string // trader id
	ready bool   // if true, need tick
	done  bool   // if true, trader done, no need to check tick

	tradeOutAction   entity.OrderAction // trade out action
	lastTradeOutTime time.Time          // last trade out time

	tickChan       chan *entity.RealTimeFutureTick // tick chan
	finishOrderMap map[string]*entity.FutureOrder  // order id -> order

	sc  *grpcapi.TradegRPCAPI
	bus *eventbus.Bus

	tradeConfig  *config.TradeFuture
	baseOrder    *entity.FutureOrder // base order
	waitingOrder *entity.FutureOrder // waiting order

	once sync.Once // post done once

	sellFirst bool
	waitTimes int64
	lastTick  *entity.RealTimeFutureTick
}

// NewDTTraderFuture create a new DTTraderFuture, if quantity > orderQtyUnit, return nil or place order error, return nil
func NewDTTraderFuture(orderWithCfg orderWithCfg, s *grpcapi.TradegRPCAPI, bus *eventbus.Bus) *DTTraderFuture {
	if orderWithCfg.order.Quantity > orderQtyUnit {
		logger.Warnf("New DTTraderFuture quantity > %d", orderQtyUnit)
		return nil
	}

	d := &DTTraderFuture{
		id:               uuid.NewString(),
		tickChan:         make(chan *entity.RealTimeFutureTick),
		finishOrderMap:   make(map[string]*entity.FutureOrder),
		sc:               s,
		bus:              bus,
		baseOrder:        orderWithCfg.order,
		tradeConfig:      orderWithCfg.cfg,
		lastTradeOutTime: orderWithCfg.lastTradeOutTime,
		waitTimes:        orderWithCfg.cfg.TradeOutWaitTimes,
	}

	if orderWithCfg.order.Action == entity.ActionSell {
		d.sellFirst = true
	}

	if err := d.placeOrder(d.baseOrder); err != nil {
		logger.Errorf("New DTTraderFuture place order err: %s", err.Error())
		return nil
	}

	switch d.baseOrder.Action {
	case entity.ActionBuy:
		d.tradeOutAction = entity.ActionSell
	case entity.ActionSell:
		d.tradeOutAction = entity.ActionBuy
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

			if d.waitingOrder != nil {
				continue
			}

			if d.isTraderDone() {
				continue
			}

			d.checkByBalance(tick)
		}
	}()
}

func (d *DTTraderFuture) isTraderDone() bool {
	if d.done {
		return true
	}

	var endQty int64
	for _, o := range d.finishOrderMap {
		if o.Status == entity.StatusFilled {
			endQty += o.Quantity
		}
	}

	if endQty == d.baseOrder.Quantity {
		d.postDone()
		return true
	}

	return false
}

func (d *DTTraderFuture) checkWaitTimes(tick *entity.RealTimeFutureTick) bool {
	defer func() {
		d.lastTick = tick
	}()

	if d.lastTick == nil {
		return true
	}

	switch d.tradeOutAction {
	case entity.ActionSell:
		if d.waitTimes > 0 && tick.Close >= d.lastTick.Close {
			d.waitTimes--
			return true
		}

	case entity.ActionBuy:
		if d.waitTimes > 0 && tick.Close <= d.lastTick.Close {
			d.waitTimes--
			return true
		}
	}

	return false
}

func (d *DTTraderFuture) checkByBalance(tick *entity.RealTimeFutureTick) {
	if d.checkWaitTimes(tick) {
		return
	}

	var place bool
	switch d.tradeOutAction {
	case entity.ActionSell:
		if tick.Close >= d.baseOrder.Price+d.tradeConfig.TargetBalanceHigh || tick.Close <= d.baseOrder.Price+d.tradeConfig.TargetBalanceLow {
			place = true
		}

	case entity.ActionBuy:
		if tick.Close <= d.baseOrder.Price-d.tradeConfig.TargetBalanceHigh || tick.Close >= d.baseOrder.Price-d.tradeConfig.TargetBalanceLow {
			place = true
		}
	}

	if !place && time.Now().Before(d.lastTradeOutTime) {
		return
	}

	o := &entity.FutureOrder{
		Code: tick.Code,
		BaseOrder: entity.BaseOrder{
			Action:   d.tradeOutAction,
			Price:    tick.Close,
			Quantity: orderQtyUnit,
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

	if _, ok := d.finishOrderMap[order.OrderID]; ok {
		d.finishOrderMap[order.OrderID] = order
		if d.waitingOrder != nil && d.waitingOrder.OrderID == order.OrderID && !order.Cancellable() {
			d.waitingOrder = nil
		}
	}
}

func (d *DTTraderFuture) placeOrder(o *entity.FutureOrder) error {
	var fn func(order *entity.FutureOrder) (*pb.TradeResult, error)
	switch o.Action {
	case entity.ActionBuy:
		fn = d.sc.BuyFuture
	case entity.ActionSell:
		fn = d.sc.SellFuture
		if d.sellFirst {
			fn = d.sc.SellFirstFuture
		}
	default:
		return nil
	}

	result, err := fn(o)
	if err != nil {
		return err
	}

	if e := d.sc.NotifyToSlack(o.String()); e != nil {
		logger.Errorf("notify to slack error: %s", e.Error())
	}

	o.OrderID = result.GetOrderId()
	if o.OrderID == "" {
		return errors.New("order id is empty")
	}

	o.Status = entity.StringToOrderStatus(result.GetStatus())
	if o.Status == entity.StatusFailed {
		return errors.New("order status is failed")
	}

	logger.Infof("%s future %s at %.2f x %d", o.Action.String(), o.Code, o.Price, o.Quantity)
	if !d.ready {
		return nil
	}
	d.finishOrderMap[o.OrderID] = o
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
