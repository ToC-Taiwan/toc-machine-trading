// Package target package target
package target

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/toc-taiwan/toc-machine-trading/internal/controller/http/websocket/ginws"
	"github.com/toc-taiwan/toc-machine-trading/internal/usecase"
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
	ticker := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-w.Ctx().Done():
			return

		case v := <-w.dataChan:
			w.SendBinaryBytesToClient(v)

		case <-ticker.C:
			m, err := w.getRank()
			if err != nil {
				continue
			}
			w.SendBinaryBytesToClient(m)
		}
	}
}
