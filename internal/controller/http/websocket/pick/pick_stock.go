// Package pick package pick
package pick

import (
	"encoding/json"
	"sync"
	"time"

	"tmt/internal/controller/http/websocket"
	"tmt/internal/usecase"

	"github.com/gin-gonic/gin"
)

type WSPickStock struct {
	*websocket.WSRouter

	s usecase.RealTime

	pickStockArr []string
	mutex        sync.Mutex
}

// StartWSPickStock -.
func StartWSPickStock(c *gin.Context, s usecase.RealTime) {
	w := &WSPickStock{
		s:        s,
		WSRouter: websocket.NewWSRouter(c),
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
				w.SendToClient(errMsg{ErrMsg: err.Error()})
				continue
			}
			w.updatePickStock(pMsg)
		}
	}()
	go w.sendPickStockSnapShot()
	w.ReadFromClient(forwardChan)
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
		case <-w.Ctx().Done():
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
			w.SendToClient(data)
		}
	}
}
