// Package v1 package v1
package v1

import (
	"net/http"

	"tmt/internal/usecase"

	"github.com/gin-gonic/gin"
)

type targetRoutes struct {
	t usecase.Target
}

func NewTargetRoutes(handler *gin.RouterGroup, t usecase.Target) {
	r := &targetRoutes{t}

	h := handler.Group("/targets")
	{
		h.GET("", r.getTargets)
	}
}

// getTargets -.
//
//	@Summary	getTargets
//	@Tags		Targets V1
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	[]entity.StockTarget
//	@Router		/v1/targets [get]
func (r *targetRoutes) getTargets(c *gin.Context) {
	c.JSON(http.StatusOK, r.t.GetTargets(c.Request.Context()))
}
