// Package v1 package v1
package v1

import (
	"net/http"

	"tmt/internal/controller/http/resp"
	"tmt/internal/controller/http/websocket/future"
	"tmt/internal/controller/http/websocket/pick"

	"tmt/internal/usecase/cases/history"
	"tmt/internal/usecase/cases/realtime"
	"tmt/internal/usecase/cases/trade"

	"github.com/gin-gonic/gin"
)

type realTimeRoutes struct {
	t realtime.RealTime
	o trade.Trade
	h history.History
}

func NewRealTimeRoutes(handler *gin.RouterGroup, t realtime.RealTime, o trade.Trade, history history.History) {
	r := &realTimeRoutes{
		t: t,
		o: o,
		h: history,
	}

	h := handler.Group("/stream")
	{
		h.GET("/tse/snapshot", r.getTSESnapshot)
		h.GET("/index", r.getIndex)

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
// @Failure     500 {object} resp.Response{}
// @Router      /stream/tse/snapshot [get]
func (r *realTimeRoutes) getTSESnapshot(c *gin.Context) {
	snapshot, err := r.t.GetTSESnapshot(c.Request.Context())
	if err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, snapshot)
}

// @Summary     getIndex
// @Description getIndex
// @ID          getIndex
// @Tags  	    stream
// @Accept      json
// @Produce     json
// @Success     200 {object} entity.TradeIndex
// @Router      /stream/index [get]
func (r *realTimeRoutes) getIndex(c *gin.Context) {
	c.JSON(http.StatusOK, r.t.GetTradeIndex())
}

func (r *realTimeRoutes) servePickStockWS(c *gin.Context) {
	pick.StartWSPickStock(c, r.t)
}

func (r *realTimeRoutes) serveFutureWS(c *gin.Context) {
	future.StartWSFutureTrade(c, r.t, r.o, r.h)
}
