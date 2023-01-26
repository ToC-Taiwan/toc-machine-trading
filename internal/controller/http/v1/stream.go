package v1

import (
	"net/http"

	"tmt/internal/controller/http/websocket/future"
	"tmt/internal/controller/http/websocket/pick"

	"tmt/internal/usecase"

	"github.com/gin-gonic/gin"
)

type streamRoutes struct {
	t usecase.Stream
	o usecase.Trade
	h usecase.History
}

func newStreamRoutes(handler *gin.RouterGroup, t usecase.Stream, o usecase.Trade, history usecase.History) {
	r := &streamRoutes{
		t: t,
		o: o,
		h: history,
	}

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
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, snapshot)
}

func (r *streamRoutes) servePickStockWS(c *gin.Context) {
	pick.StartWSPickStock(c, r.t)
}

func (r *streamRoutes) serveFutureWS(c *gin.Context) {
	future.StartWSFutureTrade(c, r.t, r.o, r.h)
}
