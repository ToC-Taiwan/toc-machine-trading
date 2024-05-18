// Package pick package pick
package pick

import (
	"context"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/toc-taiwan/toc-machine-trading/internal/controller/http/websocket/ginws"
	"github.com/toc-taiwan/toc-machine-trading/internal/entity"
	"github.com/toc-taiwan/toc-machine-trading/internal/usecase"
)

type TargetType int

const (
	Unknown TargetType = iota
	Stock
	Future
)

type WSTargetSearcher struct {
	*ginws.WSRouter
	s          usecase.Basic
	mapChan    chan string
	stockChan  chan []*entity.Stock
	futureChan chan []*entity.Future
}

// StartWSTargetSearcher -.
func StartWSTargetSearcher(c *gin.Context, s usecase.Basic, t TargetType) {
	w := &WSTargetSearcher{
		s:          s,
		WSRouter:   ginws.NewWSRouter(c),
		mapChan:    make(chan string),
		stockChan:  make(chan []*entity.Stock),
		futureChan: make(chan []*entity.Future),
	}
	forwardChan := make(chan []byte)
	go w.sendData(c.Request.Context())
	switch t {
	case Stock:
		go w.s.CreateStockSearchRoom(w.mapChan, w.stockChan)
	case Future:
		go w.s.CreateFutureSearchRoom(w.mapChan, w.futureChan)
	default:
		return
	}
	go func() {
		for {
			msg, ok := <-forwardChan
			if !ok {
				close(w.mapChan)
				return
			}
			w.mapChan <- string(msg)
		}
	}()
	w.ReadFromClient(forwardChan)
}

type stockResponse struct {
	Stocks []*entity.Stock `json:"stocks"`
	Total  int             `json:"total"`
}

type futureResponse struct {
	Futures []*entity.Future `json:"futures"`
	Total   int              `json:"total"`
}

func (w *WSTargetSearcher) sendData(ctx context.Context) {
	for {
		var response any
		select {
		case <-ctx.Done():
			return
		case data := <-w.stockChan:
			response = stockResponse{
				Stocks: data,
				Total:  len(data),
			}
		case data := <-w.futureChan:
			response = futureResponse{
				Futures: data,
				Total:   len(data),
			}
		default:
			continue
		}
		if response != nil {
			marshal, err := json.Marshal(response)
			if err != nil {
				continue
			}
			w.SendStringBytesToClient(marshal)
		}
	}
}
