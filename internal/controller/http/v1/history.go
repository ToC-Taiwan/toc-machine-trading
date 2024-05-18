// Package v1 package v1
package v1

import (
	"github.com/toc-taiwan/toc-machine-trading/internal/controller/http/websocket/history"
	"github.com/toc-taiwan/toc-machine-trading/internal/usecase"

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
	history.StartWSHistory(c, r.t)
}
