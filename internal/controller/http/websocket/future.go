package websocket

import (
	"context"
	"time"

	"tmt/internal/entity"
)

type futureOrder struct {
	Code   string             `json:"code"`
	Action entity.OrderAction `json:"action"`
	Price  float64            `json:"price"`
	Qty    int64              `json:"qty"`

	HalfAutomation bool           `json:"half_automation"`
	AutomationType AutomationType `json:"automation_type"`
}

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

type futurePosition struct {
	Position []*entity.FuturePosition `json:"position"`
}

func (w *WSRouter) cancelOverTimeOrder(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Second):
			w.orderLock.RLock()
			for _, order := range w.futureOrderMap {
				if order.Status != entity.StatusFilled && time.Since(order.TradeTime) > 10*time.Second {
					_, _, err := w.o.CancelFutureOrderID(order.OrderID)
					if err != nil {
						w.msgChan <- errMsg{err.Error()}
					}
				}
			}
			w.orderLock.RUnlock()
		}
	}
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

	var err error
	switch clientMsg.FutureOrder.Action {
	case entity.ActionBuy:
		order.OrderID, order.Status, err = w.o.BuyFuture(order)
		if err != nil {
			w.msgChan <- errMsg{ErrMsg: err.Error()}
			log.Error(err)
			return
		}

	case entity.ActionSell:
		order.OrderID, order.Status, err = w.o.SellFuture(order)
		if err != nil {
			w.msgChan <- errMsg{ErrMsg: err.Error()}
			log.Error(err)
			return
		}
	}

	w.orderLock.Lock()
	order.TradeTime = time.Now()
	w.futureOrderMap[order.OrderID] = order
	w.orderLock.Unlock()
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
	orderStatusChan := make(chan interface{})

	go w.processOrderStatus(orderStatusChan)
	go w.processTickArr(tickChan)
	go w.sendTradeIndex(ctx)
	go w.sendPosition(ctx)
	go w.cancelOverTimeOrder(ctx)

	w.s.NewFutureRealTimeConnection(tickChan, w.connectionID)
	w.s.NewOrderStatusConnection(orderStatusChan, w.connectionID)

	<-ctx.Done()
	w.s.DeleteFutureRealTimeConnection(w.connectionID)
	w.s.DeleteOrderStatusConnection(w.connectionID)
}

func (w *WSRouter) processOrderStatus(orderStatusChan chan interface{}) {
	for {
		order, ok := <-orderStatusChan
		if !ok {
			return
		}

		if o, ok := order.(*entity.FutureOrder); ok {
			w.orderLock.Lock()
			if cache := w.futureOrderMap[o.OrderID]; cache != nil && cache.Status != o.Status {
				w.msgChan <- o
				w.futureOrderMap[o.OrderID] = o
			}
			w.orderLock.Unlock()
		}
	}
}

func (w *WSRouter) processTickArr(tickChan chan *entity.RealTimeFutureTick) {
	var tickArr []*entity.RealTimeFutureTick
	for {
		tick, ok := <-tickChan
		if !ok {
			return
		}
		tickArr = append(tickArr, tick)

		var outVolume, inVolume int64
		for i := len(tickArr) - 1; i >= 0; i-- {
			if time.Since(tickArr[i].TickTime) > 15*time.Second {
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

func (w *WSRouter) sendTradeIndex(ctx context.Context) {
	if index, err := w.generateTradeIndex(ctx); err != nil {
		w.msgChan <- errMsg{ErrMsg: err.Error()}
	} else {
		w.msgChan <- index
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(5 * time.Second):
			if index, err := w.generateTradeIndex(ctx); err != nil {
				w.msgChan <- errMsg{ErrMsg: err.Error()}
			} else {
				w.msgChan <- index
			}
		}
	}
}

func (w *WSRouter) generateTradeIndex(ctx context.Context) (*tradeIndex, error) {
	nf, err := w.s.GetNasdaqFutureClose()
	if err != nil {
		return nil, err
	}

	tse, err := w.s.GetTSESnapshot(ctx)
	if err != nil {
		return nil, err
	}

	nasdaq, err := w.s.GetNasdaqClose()
	if err != nil {
		return nil, err
	}

	otc, err := w.s.GetOTCSnapshot(ctx)
	if err != nil {
		return nil, err
	}

	return &tradeIndex{
		TSE:    tse,
		OTC:    otc,
		Nasdaq: nasdaq,
		NF:     nf,
	}, nil
}

func (w *WSRouter) sendPosition(ctx context.Context) {
	if position, err := w.generatePosition(); err != nil {
		w.msgChan <- errMsg{ErrMsg: err.Error()}
	} else {
		w.msgChan <- &futurePosition{position}
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(1500 * time.Millisecond):
			if position, err := w.generatePosition(); err != nil {
				w.msgChan <- errMsg{ErrMsg: err.Error()}
			} else {
				w.msgChan <- &futurePosition{position}
			}
		}
	}
}

func (w *WSRouter) generatePosition() ([]*entity.FuturePosition, error) {
	position, err := w.o.GetFuturePosition()
	if err != nil {
		return nil, err
	}
	return position, nil
}

type tradeIndex struct {
	TSE    *entity.StockSnapShot `json:"tse"`
	OTC    *entity.StockSnapShot `json:"otc"`
	Nasdaq *entity.YahooPrice    `json:"nasdaq"`
	NF     *entity.YahooPrice    `json:"nf"`
}
