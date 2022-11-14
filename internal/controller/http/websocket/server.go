// Package websocket package websocket
package websocket

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"

	"tmt/internal/entity"
	"tmt/internal/usecase"
	"tmt/pkg/logger"

	"github.com/gin-gonic/gin"
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
	pickStockArr []string
	mutex        sync.Mutex

	s usecase.Stream
	o usecase.Order

	conn    *websocket.Conn
	msgChan chan interface{}
}

type msg struct {
	Topic         string       `json:"topic"`
	PickStockList []string     `json:"pick_stock_list"`
	FutureOrder   *futureOrder `json:"future_order"`
}

type errMsg struct {
	ErrMsg string `json:"err_msg"`
}

// NewWSRouter -.
func NewWSRouter(s usecase.Stream, o usecase.Order) *WSRouter {
	return &WSRouter{
		s:       s,
		o:       o,
		msgChan: make(chan interface{}),
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
	ctx := gin.Request.Context()

	go w.write()

	switch wsType {
	case WSPickStock:
		go w.sendPickStockSnapShot(ctx)
	case WSFuture:
		log.Info("New WSFuture")
		go w.sendFuture(ctx)
	}

	w.read(c)
}

func (w *WSRouter) read(c *websocket.Conn) {
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			if !websocket.IsCloseError(err, websocket.CloseNoStatusReceived) {
				log.Error(err)
			}
			return
		}

		if string(message) == "ping" {
			w.msgChan <- "pong"
			continue
		}

		var clientMsg msg
		err = json.Unmarshal(message, &clientMsg)
		if err != nil {
			log.Error(err)
			return
		}

		switch clientMsg.Topic {
		case "pick_stock":
			w.updatePickStock(clientMsg)
		case "future_trade":
			w.processTrade(clientMsg)
		}
	}
}

func (w *WSRouter) write() {
	for {
		cl, ok := <-w.msgChan
		if !ok {
			return
		}

		switch v := cl.(type) {
		case string:
			if err := w.send([]byte(v)); err != nil {
				return
			}

		case *entity.RealTimeFutureTick:
			serveMsgStr, err := json.Marshal(v)
			if err != nil {
				log.Error(err)
				return
			}

			if err := w.send(serveMsgStr); err != nil {
				return
			}

		default:
			log.Warn("Unknown socket message type")
			continue
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
