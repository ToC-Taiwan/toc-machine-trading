// Package websocket package websocket
package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sync"

	"tmt/internal/entity"
	"tmt/internal/usecase"
	"tmt/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var log = logger.Get()

// WSType -
type WSType int

const (
	// WSPickStock -
	WSPickStock WSType = iota + 1
	// WSFuture -
	WSFuture
)

// WSRouter -.
type WSRouter struct {
	connectionID string
	msgChan      chan interface{}
	conn         *websocket.Conn
	s            usecase.Stream
	o            usecase.Order
	ctx          context.Context

	pickStockArr []string
	mutex        sync.Mutex

	futureOrderMap map[string]*entity.FutureOrder
	orderLock      sync.Mutex
}

type clientMsg struct {
	Topic string `json:"topic"`

	PickStockList []string     `json:"pick_stock_list"`
	FutureOrder   *futureOrder `json:"future_order"`
}

type errMsg struct {
	ErrMsg string `json:"err_msg"`
}

// NewWSRouter -.
func NewWSRouter(s usecase.Stream, o usecase.Order) *WSRouter {
	return &WSRouter{
		s:              s,
		o:              o,
		connectionID:   uuid.New().String(),
		msgChan:        make(chan interface{}),
		futureOrderMap: make(map[string]*entity.FutureOrder),
	}
}

// Run -.
func (w *WSRouter) Run(gin *gin.Context, wsType WSType) {
	upGrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	c, _ := upGrader.Upgrade(gin.Writer, gin.Request, nil)
	defer func() {
		if err := c.Close(); err != nil {
			log.Errorf("Websocket Close error: %s", err)
		}
	}()

	w.conn = c
	w.ctx = gin.Request.Context()

	go w.write()

	switch wsType {
	case WSPickStock:
		go w.sendPickStockSnapShot()
	case WSFuture:
		go w.sendFuture()
	}

	w.read()
}

func (w *WSRouter) read() {
	for {
		_, message, err := w.conn.ReadMessage()
		if err != nil {
			return
		}

		if string(message) == "ping" {
			w.msgChan <- "pong"
			continue
		}

		var msg clientMsg
		if err := json.Unmarshal(message, &msg); err != nil {
			w.msgChan <- errMsg{ErrMsg: err.Error()}
			log.Error(err)
			continue
		}

		switch msg.Topic {
		case "pick_stock":
			w.updatePickStock(msg)
		case "future_trade":
			w.processTrade(msg)
		}
	}
}

func (w *WSRouter) write() {
	for {
		select {
		case <-w.ctx.Done():
			return

		case cl := <-w.msgChan:
			switch v := cl.(type) {
			case string:
				if err := w.send([]byte(v)); err != nil {
					return
				}

			default:
				if serveMsgStr, err := json.Marshal(v); err != nil {
					log.Error(err)
				} else if err := w.send(serveMsgStr); err != nil {
					return
				}
			}
		}
	}
}

func (w *WSRouter) send(data []byte) error {
	if err := w.conn.WriteMessage(websocket.TextMessage, data); err != nil {
		if !errors.Is(err, websocket.ErrCloseSent) {
			log.Error(err)
		}
		return err
	}
	return nil
}
