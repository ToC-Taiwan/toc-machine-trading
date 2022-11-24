package websocket

import (
	"time"
)

type socketPickStock struct {
	StockNum        string  `json:"stock_num"`
	StockName       string  `json:"stock_name"`
	IsTarget        bool    `json:"is_target"`
	PriceChange     float64 `json:"price_change"`
	PriceChangeRate float64 `json:"price_change_rate"`
	Price           float64 `json:"price"`
	Wrong           bool    `json:"wrong"`
}

func (w *WSRouter) updatePickStock(clientMsg clientMsg) {
	w.mutex.Lock()
	w.pickStockArr = clientMsg.PickStockList
	w.mutex.Unlock()
}

func (w *WSRouter) sendPickStockSnapShot() {
	for {
		select {
		case <-w.ctx.Done():
			return

		case <-time.After(time.Second):
			w.mutex.Lock()
			tmpStockArr := w.pickStockArr
			w.mutex.Unlock()

			if len(tmpStockArr) == 0 {
				continue
			}

			snapShot, err := w.s.GetStockSnapshotByNumArr(tmpStockArr)
			if err != nil {
				return
			}

			data := []socketPickStock{}
			for _, s := range snapShot {
				if s.StockName != "" {
					data = append(data, socketPickStock{
						StockNum:        s.StockNum,
						StockName:       s.StockName,
						IsTarget:        false,
						PriceChange:     s.PriceChg,
						PriceChangeRate: s.PctChg,
						Price:           s.Close,
					})
				} else {
					data = append(data, socketPickStock{
						StockNum: s.StockNum,
						Wrong:    true,
					})
				}
			}
			w.msgChan <- data
		}
	}
}
