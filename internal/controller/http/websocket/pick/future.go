// Package pick package pick
package pick

import (
	"time"

	"github.com/toc-taiwan/toc-machine-trading/internal/controller/http/websocket/ginws"
	"github.com/toc-taiwan/toc-machine-trading/internal/entity"
	"github.com/toc-taiwan/toc-machine-trading/internal/usecase"

	"github.com/toc-taiwan/toc-trade-protobuf/src/golang/pb"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
)

type WSPickRealFuture struct {
	*ginws.WSRouter
	b         usecase.Basic
	s         usecase.RealTime
	h         usecase.History
	code      string
	tickChan  chan *pb.FutureRealTimeTickMessage
	fetchTime time.Time
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

	w.sendInitData()

	go w.sendRealTimeData()
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
	if w.fetchTime.IsZero() {
		w.fetchTime = time.Now()
	}
	for {
		fetch, err := w.h.GetFutureHistoryPBKbarByDate(w.code, w.fetchTime)
		if err != nil {
			return nil
		} else if len(fetch.Data) == 0 {
			w.fetchTime = w.fetchTime.AddDate(0, 0, -1)
			continue
		}

		// for i := len(fetch.Data) - 1; i >= 0; i-- {
		// 	if i == 0 {
		// 		break
		// 	}

		// 	if fetch.Data[i].Ts-fetch.Data[i-1].Ts > 60*1000*1000*1000 {
		// 		fetch.Data = fetch.Data[i:]
		// 		break
		// 	}
		// }

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

func (w *WSPickRealFuture) getFutureSnapshot() *pb.WSMessage {
	snapshot, err := w.s.GetFutureSnapshotByCode(w.code)
	if err != nil {
		return nil
	}
	return &pb.WSMessage{
		Data: &pb.WSMessage_Snapshot{
			Snapshot: snapshot,
		},
	}
}

func (w *WSPickRealFuture) sendRealTimeData() {
	ticker := time.NewTicker(time.Second * 5)
	minuteTicker := time.NewTicker(time.Minute)
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

		case <-ticker.C:
			if data := w.generateTradeIndex(); data != nil {
				w.sendMessage(data)
			}

		case <-minuteTicker.C:
			if data := w.fetchKbar(); data != nil {
				w.sendMessage(data)
			}
		}
	}
}

func (w *WSPickRealFuture) sendInitData() {
	if data := w.getFutureSnapshot(); data != nil {
		w.sendMessage(data)
	}

	if data := w.generateTradeIndex(); data != nil {
		w.sendMessage(data)
	}

	if data := w.fetchKbar(); data != nil {
		w.sendMessage(data)
	}

	if future := w.getFutureDetail(); future != nil {
		w.sendMessage(future)
	}
}

func (w *WSPickRealFuture) sendMessage(msg *pb.WSMessage) {
	content, err := proto.Marshal(msg)
	if err != nil {
		return
	}
	w.SendBinaryBytesToClient(content)
}
