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
		h.GET("/tse/snapshot", r.getTSESnapshot)
	}
}

// @Summary     getTSESnapshot
// @Description getTSESnapshot
// @ID          getTSESnapshot
// @Tags  	    stream
// @Accept      json
// @Produce     json
// @Success     200 {object} entity.StockSnapShot
// @Failure     500 {object} response
// @Router      /stream/tse/snapshot [get]
func (r *streamRoutes) getTSESnapshot(c *gin.Context) {
	snapshot, err := r.t.GetTSESnapshot(c.Request.Context())
	if err != nil {
		log.Error(err)
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, snapshot)
}
