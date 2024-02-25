// Package ginws package ginws
package ginws

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WSRouter -.
type WSRouter struct {
	binaryBytesChan chan []byte
	stringBytesChan chan []byte

	conn *websocket.Conn
	ctx  context.Context
}

// NewWSRouter -.
func NewWSRouter(c *gin.Context) *WSRouter {
	r := &WSRouter{
		stringBytesChan: make(chan []byte),
		binaryBytesChan: make(chan []byte),
		ctx:             c.Request.Context(),
	}
	r.upgrade(c)
	return r
}

// Upgrade -.
func (w *WSRouter) upgrade(gin *gin.Context) {
	upGrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	c, err := upGrader.Upgrade(gin.Writer, gin.Request, nil)
	if err != nil {
		return
	}

	w.conn = c
	go w.write()
}

func (w *WSRouter) write() {
	for {
		select {
		case <-w.Ctx().Done():
			return

		case cl := <-w.stringBytesChan:
			w.sendText(cl)

		case cl := <-w.binaryBytesChan:
			w.sendBinary(cl)
		}
	}
}

func (w *WSRouter) sendText(data []byte) {
	_ = w.conn.WriteMessage(websocket.TextMessage, data)
}

func (w *WSRouter) sendBinary(data []byte) {
	_ = w.conn.WriteMessage(websocket.BinaryMessage, data)
}

func (w *WSRouter) ReadFromClient(forwardChan chan []byte) {
	for {
		_, message, err := w.conn.ReadMessage()
		if err != nil {
			close(forwardChan)
			break
		}

		if string(message) == "ping" {
			w.stringBytesChan <- []byte("pong")
			continue
		}

		forwardChan <- message
	}
	_ = w.conn.Close()
}

func (w *WSRouter) Wait() {
	for {
		_, _, err := w.conn.ReadMessage()
		if err != nil {
			break
		}
	}
	_ = w.conn.Close()
}

func (w *WSRouter) SendStringBytesToClient(msg []byte) {
	if msg == nil {
		return
	}
	w.stringBytesChan <- msg
}

func (w *WSRouter) SendBinaryBytesToClient(msg []byte) {
	if msg == nil {
		return
	}
	w.binaryBytesChan <- msg
}

func (w *WSRouter) Ctx() context.Context {
	return w.ctx
}
