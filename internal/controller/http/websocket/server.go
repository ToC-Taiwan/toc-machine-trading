// Package websocket package websocket
package websocket

import (
	"encoding/json"
	"net/http"
	"sync"

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
		go w.sendFuture(ctx)
	}

	w.read(c)
}

func (w *WSRouter) read(c *websocket.Conn) {
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			return
		}

		if string(message) == "ping" {
			_ = c.WriteMessage(websocket.TextMessage, []byte("pong"))
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

		serveMsgStr, err := json.Marshal(cl)
		if err != nil {
			log.Error(err)
			return
		}

		err = w.conn.WriteMessage(websocket.TextMessage, serveMsgStr)
		if err != nil {
			log.Error(err)
			return
		}
	}
}
