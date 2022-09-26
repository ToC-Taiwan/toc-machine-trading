package v1

import (
	"net/http"

	"tmt/internal/controller/http/websocket"
	"tmt/internal/usecase"

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
		h.GET("/ws/pick-stock", r.servePickStockWS)
		h.GET("/ws/future", r.serveFutureWS)
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

func (r *streamRoutes) servePickStockWS(c *gin.Context) {
	wsRouter := websocket.NewWSRouter(r.t)
	wsRouter.Run(c, websocket.WSPickStock)
}

func (r *streamRoutes) serveFutureWS(c *gin.Context) {
	wsRouter := websocket.NewWSRouter(r.t)
	wsRouter.Run(c, websocket.WSFuture)
}
