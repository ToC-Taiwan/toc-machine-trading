// Package v1 package v1
package v1

import (
	"net/http"

	"tmt/internal/controller/http/resp"
	"tmt/internal/controller/http/websocket/future"
	"tmt/internal/controller/http/websocket/pick"
	"tmt/internal/usecase"

	"github.com/gin-gonic/gin"
)

type realTimeRoutes struct {
	t usecase.RealTime
	o usecase.Trade
	h usecase.History
}

func NewRealTimeRoutes(handler *gin.RouterGroup, t usecase.RealTime, o usecase.Trade, history usecase.History) {
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

// getTSESnapshot -.
//
//	@Tags		Stream V1
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	entity.StockSnapShot
//	@Failure	500	{object}	resp.Response{}
//	@Router		/v1/stream/tse/snapshot [get]
func (r *realTimeRoutes) getTSESnapshot(c *gin.Context) {
	snapshot, err := r.t.GetTSESnapshot(c.Request.Context())
	if err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, snapshot)
}

// getIndex -.
//
//	@Tags		Stream V1
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	entity.TradeIndex
//	@Router		/v1/stream/index [get]
func (r *realTimeRoutes) getIndex(c *gin.Context) {
	c.JSON(http.StatusOK, r.t.GetTradeIndex())
}

func (r *realTimeRoutes) servePickStockWS(c *gin.Context) {
	pick.StartWSPickStock(c, r.t)
}

func (r *realTimeRoutes) serveFutureWS(c *gin.Context) {
	future.StartWSFutureTrade(c, r.t, r.o, r.h)
}
