// Package websocket package websocket
package websocket

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// WSRouter -.
type WSRouter struct {
	connectionID string
	msgChan      chan interface{}
	conn         *websocket.Conn
	ctx          context.Context
}

// NewWSRouter -.
func NewWSRouter(c *gin.Context) *WSRouter {
	r := &WSRouter{
		connectionID: uuid.New().String(),
		msgChan:      make(chan interface{}),
	}
	r.Upgrade(c)
	return r
}

// Upgrade -.
func (w *WSRouter) Upgrade(gin *gin.Context) {
	upGrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	c, _ := upGrader.Upgrade(gin.Writer, gin.Request, nil)
	defer func() { _ = c.Close() }()

	w.conn = c
	w.ctx = gin.Request.Context()

	go w.write()
}

func (w *WSRouter) send(data []byte) error {
	if err := w.conn.WriteMessage(websocket.TextMessage, data); err != nil {
		return err
	}
	return nil
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
				if serveMsgStr, err := json.Marshal(v); err == nil {
					if err := w.send(serveMsgStr); err != nil {
						return
					}
				}
			}
		}
	}
}

func (w *WSRouter) read(forwardChan chan []byte) {
	go func() {
		for {
			_, message, err := w.conn.ReadMessage()
			if err != nil {
				close(forwardChan)
				return
			}

			if string(message) == "ping" {
				w.msgChan <- "pong"
				continue
			}

			forwardChan <- message
		}
	}()
}
