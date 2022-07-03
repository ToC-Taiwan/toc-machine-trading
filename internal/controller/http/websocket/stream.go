package websocket

import (
	"context"
	"time"
)

// SocketPickStock -.
type SocketPickStock struct {
	StockNum        string  `json:"stock_num"`
	StockName       string  `json:"stock_name"`
	IsTarget        bool    `json:"is_target"`
	PriceChange     float64 `json:"price_change"`
	PriceChangeRate float64 `json:"price_change_rate"`
	Price           float64 `json:"price"`
	Wrong           bool    `json:"wrong"`
}

func (w *WSRouter) sendSnapShotArr(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			close(w.msgChan)
			return
		default:
			time.Sleep(time.Second)
			w.mutex.Lock()
			tmpStockArr := w.pickStockArr
			w.mutex.Unlock()

			if len(tmpStockArr) == 0 {
				continue
			}

			snapShot, err := w.s.GetStockSnapshotByNumArr(tmpStockArr)
			if err != nil {
				log.Error(err)
				return
			}

			data := []SocketPickStock{}
			for _, s := range snapShot {
				if s.StockName != "" {
					data = append(data, SocketPickStock{
						StockNum:        s.StockNum,
						StockName:       s.StockName,
						IsTarget:        false,
						PriceChange:     s.PriceChg,
						PriceChangeRate: s.PctChg,
						Price:           s.Close,
					})
				} else {
					data = append(data, SocketPickStock{
						StockNum: s.StockNum,
						Wrong:    true,
					})
				}
			}

			w.msgChan <- data
		}
	}
}
