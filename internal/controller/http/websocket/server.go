// Package websocket package websocket
package websocket

import (
	"encoding/json"
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

	s       usecase.Stream
	o       usecase.Order
	conn    *websocket.Conn
	msgChan chan interface{}
}

type msg struct {
	Topic         string       `json:"topic"`
	PickStockList []string     `json:"pick_stock_list"`
	FutureOrder   *futureOrder `json:"future_order"`
}

type futureOrder struct {
	Code   string             `json:"code"`
	Action entity.OrderAction `json:"action"`
	Price  float64            `json:"price"`
	Qty    int64              `json:"qty"`
}

type errMsg struct {
	ErrMsg string `json:"err_msg"`
}

// NewWSRouter -.
func NewWSRouter(s usecase.Stream, o usecase.Order) *WSRouter {
	r := &WSRouter{
		s:       s,
		o:       o,
		msgChan: make(chan interface{}),
	}
	return r
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
		go w.sendSnapShotArr(ctx)
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

func (w *WSRouter) updatePickStock(clientMsg msg) {
	w.mutex.Lock()
	w.pickStockArr = clientMsg.PickStockList
	w.mutex.Unlock()
}

func (w *WSRouter) processTrade(clientMsg msg) {
	if clientMsg.FutureOrder == nil {
		return
	}

	order := &entity.FutureOrder{
		Code: clientMsg.FutureOrder.Code,
		BaseOrder: entity.BaseOrder{
			Action:   clientMsg.FutureOrder.Action,
			Quantity: clientMsg.FutureOrder.Qty,
			Price:    clientMsg.FutureOrder.Price,
		},
	}

	switch clientMsg.FutureOrder.Action {
	case entity.ActionBuy:
		if _, _, err := w.o.BuyFuture(order); err != nil {
			w.msgChan <- errMsg{ErrMsg: err.Error()}
		}
	case entity.ActionSell:
		if _, _, err := w.o.SellFuture(order); err != nil {
			w.msgChan <- errMsg{ErrMsg: err.Error()}
		}
	}
}

func (w *WSRouter) write() {
	for {
		cl, ok := <-w.msgChan
		if !ok {
			return
		}

		serveMsgStr, _ := json.Marshal(cl)
		_ = w.conn.WriteMessage(websocket.TextMessage, serveMsgStr)
	}
}
