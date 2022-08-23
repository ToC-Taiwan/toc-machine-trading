package v1

import (
	"net/http"

	"tmt/internal/usecase"

	"github.com/gin-gonic/gin"
)

type targetRoutes struct {
	t usecase.Target
}

func newTargetRoutes(handler *gin.RouterGroup, t usecase.Target) {
	r := &targetRoutes{t}

	h := handler.Group("/targets")
	{
		h.GET("/", r.getTargets)
	}
}

// @Summary     getTargets
// @Description getTargets
// @ID          getTargets
// @Tags  	    targets
// @Accept      json
// @Produce     json
// @Success     200 {object} []entity.Target
// @Router      /targets [get]
func (r *targetRoutes) getTargets(c *gin.Context) {
	c.JSON(http.StatusOK, r.t.GetTargets(c.Request.Context()))
}
