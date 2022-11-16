package websocket

import (
	"context"
	"time"

	"tmt/internal/entity"

	"github.com/google/uuid"
)

type futureOrder struct {
	Code   string             `json:"code"`
	Action entity.OrderAction `json:"action"`
	Price  float64            `json:"price"`
	Qty    int64              `json:"qty"`

	HalfAutomation bool           `json:"half_automation"`
	AutomationType AutomationType `json:"automation_type"`
}

// type orderResult struct {
// 	OrderID string             `json:"order_id"`
// 	Status  entity.OrderStatus `json:"status"`
// }

type AutomationType int

const (
	// WSPickStock -
	AutomationByTime AutomationType = iota + 1
	// WSFuture -
	AutomationByBalance
)

type tradeRate struct {
	OutRate int64 `json:"out_rate"`
	InRate  int64 `json:"in_rate"`
}

type orderIDWithStatus struct {
	OrderID string             `json:"order_id"`
	Status  entity.OrderStatus `json:"status"`
}

func (w *WSRouter) processTrade(clientMsg msg) {
	if clientMsg.FutureOrder == nil {
		return
	}

	order := &entity.FutureOrder{
		Code: clientMsg.FutureOrder.Code,
		BaseOrder: entity.BaseOrder{
			Action:   clientMsg.FutureOrder.Action,
			Quantity: clientMsg.FutureOrder.Qty,
			Price:    clientMsg.FutureOrder.Price,
		},
	}

	switch clientMsg.FutureOrder.Action {
	case entity.ActionBuy:
		orderID, s, err := w.o.BuyFuture(order)
		if err != nil {
			w.msgChan <- errMsg{ErrMsg: err.Error()}
			log.Error(err)
		}
		w.checkOrderChan <- orderIDWithStatus{OrderID: orderID, Status: s}

	case entity.ActionSell:
		orderID, s, err := w.o.SellFuture(order)
		if err != nil {
			w.msgChan <- errMsg{ErrMsg: err.Error()}
			log.Error(err)
		}
		w.checkOrderChan <- orderIDWithStatus{OrderID: orderID, Status: s}
	}
}

func (w *WSRouter) checkOrderStatus(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		case o := <-w.checkOrderChan:
			if o.OrderID == "" {
				w.msgChan <- errMsg{ErrMsg: "order id is empty"}
			}

			if o.Status == entity.StatusFailed {
				w.msgChan <- errMsg{ErrMsg: "order failed"}
			}
			go w.queryOrderStatus(o.OrderID, ctx)
		}
	}
}

func (w *WSRouter) queryOrderStatus(orderID string, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		case <-time.After(1 * time.Second):
			order, err := w.o.GetFutureOrderStatusByID(orderID)
			if err != nil {
				w.msgChan <- errMsg{ErrMsg: err.Error()}
				log.Error(err)
				return
			}

			if order.Status == entity.StatusFilled || order.Status == entity.StatusFilling {
				w.msgChan <- order
				return
			}
		}
	}
}

func (w *WSRouter) sendFuture(ctx context.Context) {
	snapshot, err := w.s.GetFutureSnapshotByCode(w.s.GetMainFutureCode())
	if err != nil {
		w.msgChan <- errMsg{ErrMsg: err.Error()}
	}

	var tickType, chgType int64
	switch snapshot.TickType {
	case "Sell":
		tickType = 1
	case "Buy":
		tickType = 2
	}

	switch snapshot.ChgType {
	case "LimitUp":
		chgType = 1
	case "Up":
		chgType = 2
	case "Unchanged":
		chgType = 3
	case "Dowm":
		chgType = 4
	case "LimitDown":
		chgType = 5
	}

	w.msgChan <- &entity.RealTimeFutureTick{
		Code:        snapshot.Code,
		TickTime:    snapshot.SnapTime,
		Open:        snapshot.Open,
		Close:       snapshot.Close,
		High:        snapshot.High,
		Low:         snapshot.Low,
		Amount:      float64(snapshot.Amount),
		TotalAmount: float64(snapshot.AmountSum),
		Volume:      snapshot.Volume,
		TotalVolume: snapshot.VolumeSum,
		TickType:    tickType,
		ChgType:     chgType,
		PriceChg:    snapshot.PriceChg,
		PctChg:      snapshot.PctChg,
	}

	tickChan := make(chan *entity.RealTimeFutureTick)
	go w.processTickArr(tickChan)

	connectionID := uuid.New().String()
	w.s.NewFutureRealTimeConnection(tickChan, connectionID)

	<-ctx.Done()
	w.s.DeleteFutureRealTimeConnection(connectionID)
}

func (w *WSRouter) processTickArr(tickChan chan *entity.RealTimeFutureTick) {
	var tickArr []*entity.RealTimeFutureTick
	for {
		tick, ok := <-tickChan
		if !ok {
			close(w.msgChan)
			return
		}
		tickArr = append(tickArr, tick)

		var outVolume, inVolume int64
		for i := len(tickArr) - 1; i >= 0; i-- {
			if time.Since(tickArr[i].TickTime) > 30*time.Second {
				tickArr = tickArr[i+1:]
				break
			}

			switch tickArr[i].TickType {
			case 1:
				outVolume += tickArr[i].Volume
			case 2:
				inVolume += tickArr[i].Volume
			}
		}

		w.msgChan <- &tradeRate{
			OutRate: outVolume,
			InRate:  inVolume,
		}
		w.msgChan <- tick
	}
}
