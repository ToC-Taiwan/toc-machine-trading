// Package v1 package v1
package v1

import (
	wsHistory "tmt/internal/controller/http/websocket/history"
	"tmt/internal/usecase"

	"github.com/gin-gonic/gin"
)

type historyRoutes struct {
	t usecase.History
}

func NewHistoryRoutes(handler *gin.RouterGroup, t usecase.History) {
	r := &historyRoutes{t}

	h := handler.Group("/history")
	{
		h.GET("/ws", r.serveWS)
	}
}

func (r *historyRoutes) serveWS(c *gin.Context) {
	wsHistory.StartWSHistory(c, r.t)
}
