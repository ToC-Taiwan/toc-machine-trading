// Package v1 package v1
package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/toc-taiwan/toc-machine-trading/internal/controller/http/websocket/target"
	"github.com/toc-taiwan/toc-machine-trading/internal/usecase"
)

type targetRoutes struct {
	t usecase.Target
}

func NewTargetRoutes(handler *gin.RouterGroup, t usecase.Target) {
	r := &targetRoutes{t}

	h := handler.Group("/targets")
	{
		h.GET("", r.getTargets)
		h.GET("/ws", r.serveWS)
	}
}

// getTargets -.
//
//	@Tags		Targets V1
//	@Summary	Get targets
//	@security	JWT
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	[]entity.StockTarget
//	@Router		/v1/targets [get]
func (r *targetRoutes) getTargets(c *gin.Context) {
	c.JSON(http.StatusOK, r.t.GetTargets(c.Request.Context()))
}

func (r *targetRoutes) serveWS(c *gin.Context) {
	target.StartWSTargetStock(c, r.t)
}
