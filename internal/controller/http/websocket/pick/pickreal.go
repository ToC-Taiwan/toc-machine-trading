// Package pick package pick
package pick

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/toc-taiwan/toc-machine-trading/internal/controller/http/websocket/ginws"
	"github.com/toc-taiwan/toc-machine-trading/internal/usecase"
	"github.com/toc-taiwan/toc-trade-protobuf/golang/pb"
	"google.golang.org/protobuf/proto"
)

type WSPickRealStock struct {
	*ginws.WSRouter
	s        usecase.RealTime
	mapChan  chan *pb.PickRealMap
	tickChan chan []byte
}

// StartWSPickStock -.
func StartWSPickStock(c *gin.Context, s usecase.RealTime, odd bool) {
	w := &WSPickRealStock{
		s:        s,
		WSRouter: ginws.NewWSRouter(c),
		mapChan:  make(chan *pb.PickRealMap),
		tickChan: make(chan []byte),
	}
	forwardChan := make(chan []byte)
	connectionID := uuid.New().String()
	go w.sendRealStock()
	go w.s.CreateRealTimePick(connectionID, odd, w.mapChan, w.tickChan)
	go func() {
		for {
			msg, ok := <-forwardChan
			if !ok {
				close(w.mapChan)
				return
			}
			var pickRequest pb.PickRealMap
			if err := proto.Unmarshal(msg, &pickRequest); err != nil {
				continue
			}
			w.mapChan <- &pickRequest
		}
	}()
	w.ReadFromClient(forwardChan)
	w.s.DeleteRealTimeClient(connectionID)
	close(w.tickChan)
}

func (w *WSPickRealStock) sendRealStock() {
	for {
		tick, ok := <-w.tickChan
		if !ok {
			return
		}
		w.SendBinaryBytesToClient(tick)
	}
}
