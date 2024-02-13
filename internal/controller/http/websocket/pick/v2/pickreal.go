// Package pick package pick
package pick

import (
	"tmt/internal/controller/http/websocket/ginws"
	"tmt/internal/usecase"
	"tmt/pb"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

type WSPickRealStock struct {
	*ginws.WSRouter
	s        usecase.RealTime
	mapChan  chan *pb.PickRealMap
	tickChan chan []byte
}

// StartWSPickStock -.
func StartWSPickStock(c *gin.Context, s usecase.RealTime) {
	w := &WSPickRealStock{
		s:        s,
		WSRouter: ginws.NewWSRouter(c),
		mapChan:  make(chan *pb.PickRealMap),
		tickChan: make(chan []byte),
	}
	forwardChan := make(chan []byte)
	connectionID := uuid.New().String()
	go w.sendRealStock()
	go w.s.CreateRealTimePick(connectionID, w.mapChan, w.tickChan)
	go func() {
		for {
			msg, ok := <-forwardChan
			if !ok {
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
		w.SendBinaryToClient(tick)
	}
}
