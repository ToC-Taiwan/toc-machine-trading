package httpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"tmt/pb"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

// WSRouter -.
type WSRouter struct {
	msgChan chan interface{}
	conn    *websocket.Conn
	ctx     context.Context
	logger  Logger
}

// NewWSRouter -.
func NewWSRouter(c *gin.Context, logger Logger) *WSRouter {
	r := &WSRouter{
		msgChan: make(chan interface{}),
		ctx:     c.Request.Context(),
		logger:  logger,
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
		w.logger.Errorf("upgrade websocket err: %v", err)
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

		case cl := <-w.msgChan:
			switch v := cl.(type) {
			case string:
				w.sendText([]byte(v))

			case *pb.WSMessage:
				if serveMsgStr, err := proto.Marshal(v); err == nil {
					w.sendBinary(serveMsgStr)
				}

			default:
				if serveMsgStr, err := json.Marshal(v); err == nil {
					w.sendText(serveMsgStr)
				}
			}
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
			w.msgChan <- "pong"
			continue
		}

		forwardChan <- message
	}

	if err := w.conn.Close(); err != nil {
		w.logger.Errorf("websocket close err: %v", err)
	}
}

func (w *WSRouter) SendToClient(msg interface{}) {
	w.msgChan <- msg
}

func (w *WSRouter) Ctx() context.Context {
	return w.ctx
}
