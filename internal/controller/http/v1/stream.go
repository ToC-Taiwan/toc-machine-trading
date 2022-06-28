package v1

import (
	"net/http"

	"toc-machine-trading/internal/usecase"

	"github.com/gin-gonic/gin"
)

type streamRoutes struct {
	t usecase.Stream
}

func newStreamRoutes(handler *gin.RouterGroup, t usecase.Stream) {
	r := &streamRoutes{t}

	h := handler.Group("/stream")
	{
		h.GET("/", r.getTSESnapshot)
	}
}

// @Summary     getTSESnapshot
// @Description getTSESnapshot
// @ID          getTSESnapshot
// @Tags  	    stream
// @Accept      json
// @Produce     json
// @Router      /stream [get]
func (r *streamRoutes) getTSESnapshot(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}
