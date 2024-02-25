// Package pick package pick
package pick

import (
	"time"

	"tmt/internal/controller/http/websocket/ginws"
	"tmt/internal/entity"
	"tmt/internal/usecase"
	"tmt/pb"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
)

type WSPickRealFuture struct {
	*ginws.WSRouter
	b        usecase.Basic
	s        usecase.RealTime
	h        usecase.History
	code     string
	tickChan chan *pb.FutureRealTimeTickMessage
}

// StartWSPickRealFuture -.
func StartWSPickRealFuture(c *gin.Context, code string, s usecase.RealTime, h usecase.History, b usecase.Basic) {
	w := &WSPickRealFuture{
		code:     code,
		s:        s,
		h:        h,
		b:        b,
		WSRouter: ginws.NewWSRouter(c),
		tickChan: make(chan *pb.FutureRealTimeTickMessage),
	}
	go w.sendData()
	go w.s.CreateRealTimePickFuture(w.Ctx(), w.code, w.tickChan)
	w.Wait()
}

func (w *WSPickRealFuture) generateTradeIndex() *pb.WSMessage {
	index := w.s.GetTradeIndex()
	data := &pb.TradeIndex{
		Tse: &pb.IndexStatus{
			BreakCount: index.TSE.BreakCount,
			PriceChg:   index.TSE.PriceChg,
		},
		Otc: &pb.IndexStatus{
			BreakCount: index.OTC.BreakCount,
			PriceChg:   index.OTC.PriceChg,
		},
		Nasdaq: &pb.IndexStatus{
			BreakCount: index.Nasdaq.BreakCount,
			PriceChg:   index.Nasdaq.PriceChg,
		},
		Nf: &pb.IndexStatus{
			BreakCount: index.NF.BreakCount,
			PriceChg:   index.NF.PriceChg,
		},
	}
	return &pb.WSMessage{
		Data: &pb.WSMessage_TradeIndex{
			TradeIndex: data,
		},
	}
}

func (w *WSPickRealFuture) fetchKbar() *pb.WSMessage {
	fetchTime := time.Now()
	for {
		fetch, err := w.h.GetFutureHistoryPBKbarByDate(w.code, fetchTime)
		if err != nil {
			return nil
		} else if len(fetch.Data) == 0 {
			fetchTime = fetchTime.AddDate(0, 0, -1)
			continue
		}
		return &pb.WSMessage{
			Data: &pb.WSMessage_HistoryKbar{
				HistoryKbar: fetch,
			},
		}
	}
}

func (w *WSPickRealFuture) getFutureDetail() *pb.WSMessage {
	future := w.b.GetFutureDetail(w.code)
	if future == nil {
		return nil
	}
	data := &pb.FutureDetailMessage{
		Code:           future.Code,
		Symbol:         future.Symbol,
		Name:           future.Name,
		Category:       future.Category,
		DeliveryMonth:  future.DeliveryMonth,
		DeliveryDate:   future.DeliveryDate.Format(entity.ShortTimeLayout),
		UnderlyingKind: future.UnderlyingKind,
		Unit:           future.Unit,
		LimitUp:        future.LimitUp,
		LimitDown:      future.LimitDown,
		Reference:      future.Reference,
		UpdateDate:     future.UpdateDate.Format(entity.ShortTimeLayout),
	}
	return &pb.WSMessage{
		Data: &pb.WSMessage_FutureDetail{
			FutureDetail: data,
		},
	}
}

func (w *WSPickRealFuture) sendData() {
	if future := w.getFutureDetail(); future != nil {
		w.sendMessage(future)
	}

	if data := w.generateTradeIndex(); data != nil {
		w.sendMessage(data)
	}

	if data := w.fetchKbar(); data != nil {
		w.sendMessage(data)
	}

	for {
		select {
		case <-w.Ctx().Done():
			return

		case tick := <-w.tickChan:
			w.sendMessage(&pb.WSMessage{
				Data: &pb.WSMessage_FutureTick{
					FutureTick: tick,
				},
			})

		case <-time.After(time.Second * 5):
			if data := w.generateTradeIndex(); data != nil {
				w.sendMessage(data)
			}

		case <-time.After(time.Minute):
			if data := w.fetchKbar(); data != nil {
				w.sendMessage(data)
			}
		}
	}
}

func (w *WSPickRealFuture) sendMessage(msg *pb.WSMessage) {
	content, err := proto.Marshal(msg)
	if err != nil {
		return
	}
	w.SendBinaryBytesToClient(content)
}
