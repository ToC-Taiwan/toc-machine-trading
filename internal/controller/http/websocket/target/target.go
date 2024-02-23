// Package target package target
package target

import (
	"time"

	"tmt/internal/controller/http/websocket/ginws"
	"tmt/internal/usecase"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
)

type WSTargetStock struct {
	*ginws.WSRouter
	s        usecase.Target
	dataChan chan []byte
}

// StartWSTargetStock -.
func StartWSTargetStock(c *gin.Context, s usecase.Target) {
	w := &WSTargetStock{
		s:        s,
		WSRouter: ginws.NewWSRouter(c),
		dataChan: make(chan []byte),
	}
	forwardChan := make(chan []byte)
	go w.sender()
	go func() {
		for {
			_, ok := <-forwardChan
			if !ok {
				close(w.dataChan)
				return
			}
		}
	}()
	data, _ := w.getRank()
	w.dataChan <- data
	w.ReadFromClient(forwardChan)
}

func (w *WSTargetStock) getRank() ([]byte, error) {
	data, err := w.s.GetCurrentVolumeRank()
	if err != nil {
		return nil, err
	}
	m, err := proto.Marshal(data)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (w *WSTargetStock) sender() {
	for {
		select {
		case <-w.Ctx().Done():
			return

		case v := <-w.dataChan:
			w.SendBinaryBytesToClient(v)

		case <-time.After(time.Second * 10):
			m, err := w.getRank()
			if err != nil {
				continue
			}
			w.SendBinaryBytesToClient(m)
		}
	}
}
