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

// WSRouter -.
type WSRouter struct {
	pickStockArr []string
	mutex        sync.Mutex

	s       usecase.Stream
	conn    *websocket.Conn
	msgChan chan interface{}
}

type msg struct {
	Data interface{} `json:"data"`
}

// NewWSRouter -.
func NewWSRouter(s usecase.Stream) *WSRouter {
	r := &WSRouter{
		s:       s,
		msgChan: make(chan interface{}),
	}
	return r
}

// Run -.
func (w *WSRouter) Run(gin *gin.Context) {
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
	go w.sendSnapShotArr(ctx)

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

		var pickStock []string
		if arr, ok := clientMsg.Data.(map[string]interface{})["pick_stock_list"].([]interface{}); ok {
			for _, v := range arr {
				if stock, ok := v.(string); ok {
					pickStock = append(pickStock, stock)
				}
			}
			w.mutex.Lock()
			w.pickStockArr = pickStock
			w.mutex.Unlock()
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
