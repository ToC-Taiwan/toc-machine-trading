// Package history package history
package history

import (
	"encoding/json"
	"time"

	"tmt/internal/controller/http/websocket/ginws"
	"tmt/internal/entity"
	"tmt/internal/usecase"
	"tmt/pb"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
)

type WSHistory struct {
	*ginws.WSRouter
	s        usecase.History
	reqChan  chan []byte
	dataChan chan []byte
}

// StartWSHistory -.
func StartWSHistory(c *gin.Context, s usecase.History) {
	w := &WSHistory{
		s:        s,
		WSRouter: ginws.NewWSRouter(c),
		dataChan: make(chan []byte),
		reqChan:  make(chan []byte),
	}
	forwardChan := make(chan []byte)
	go w.sender()
	go w.getData()
	go func() {
		for {
			req, ok := <-forwardChan
			if !ok {
				close(w.reqChan)
				return
			}
			w.reqChan <- req
		}
	}()
	w.ReadFromClient(forwardChan)
}

func (w *WSHistory) sender() {
	for {
		select {
		case <-w.Ctx().Done():
			return

		case v := <-w.dataChan:
			w.SendBinaryToClient(v)
		}
	}
}

type kbarReq struct {
	StockNum  string `json:"stock_num"`
	StartDate string `json:"start_date"`
	Interval  int64  `json:"interval"`
}

func (w *WSHistory) getData() {
	for {
		select {
		case <-w.Ctx().Done():
			return

		case req := <-w.reqChan:
			var r kbarReq
			if err := json.Unmarshal(req, &r); err != nil {
				continue
			}
			startDateTime, err := time.ParseInLocation(entity.ShortTimeLayout, r.StartDate, time.Local)
			if err != nil {
				continue
			}

			data, err := w.s.GetDayKbarByStockNumMultiDate(r.StockNum, startDateTime, r.Interval)
			if err != nil {
				continue
			}
			result := &pb.WSHistoryKbarMessage{}
			for _, v := range data {
				result.Arr = append(result.Arr, &pb.Kbar{
					Open:     v.Open,
					Close:    v.Close,
					High:     v.High,
					Low:      v.Low,
					Volume:   v.Volume,
					KbarTime: v.KbarTime.Format(entity.ShortTimeLayout),
				})
			}
			b, err := proto.Marshal(result)
			if err != nil {
				continue
			}
			w.dataChan <- b
		}
	}
}
