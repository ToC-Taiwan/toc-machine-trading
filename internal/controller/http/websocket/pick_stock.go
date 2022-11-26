package websocket

import (
	"encoding/json"
	"sync"
	"time"

	"tmt/internal/usecase"

	"github.com/gin-gonic/gin"
)

type WSPickStock struct {
	*WSRouter

	s usecase.Stream

	pickStockArr []string
	mutex        sync.Mutex
}

// StartWSPickStock -.
func StartWSPickStock(c *gin.Context, s usecase.Stream) {
	w := &WSPickStock{
		s:        s,
		WSRouter: NewWSRouter(c),
	}

	forwardChan := make(chan []byte)
	go func() {
		for {
			msg, ok := <-forwardChan
			if !ok {
				return
			}

			var pMsg pickStockClientMsg
			if err := json.Unmarshal(msg, &pMsg); err != nil {
				w.msgChan <- errMsg{ErrMsg: err.Error()}
				continue
			}
			w.updatePickStock(pMsg)
		}
	}()
	w.read(forwardChan)
	w.sendPickStockSnapShot()
}

type pickStockClientMsg struct {
	PickStockList []string `json:"pick_stock_list"`
}

type socketPickStock struct {
	StockNum        string  `json:"stock_num"`
	StockName       string  `json:"stock_name"`
	IsTarget        bool    `json:"is_target"`
	PriceChange     float64 `json:"price_change"`
	PriceChangeRate float64 `json:"price_change_rate"`
	Price           float64 `json:"price"`
	Wrong           bool    `json:"wrong"`
}

func (w *WSPickStock) updatePickStock(clientMsg pickStockClientMsg) {
	w.mutex.Lock()
	w.pickStockArr = clientMsg.PickStockList
	w.mutex.Unlock()
}

func (w *WSPickStock) sendPickStockSnapShot() {
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
