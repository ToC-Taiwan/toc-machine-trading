// Package v1 package v1
package v1

import (
	"net/http"

	"tmt/internal/controller/http/resp"
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
		h.PUT("/snapshot", r.getSnapshots)
		h.GET("/ws/pick-stock", r.servePickStockWS)
		h.GET("/ws/pick-stock/odds", r.servePickStockOddsWS)
	}
}

type snapshotRequest struct {
	StockList []string `json:"stock_list"`
}

// getSnapshots -.
//
//	@Tags		Stream V1
//	@Summary	Get snapshots
//	@security	JWT
//	@Accept		json
//	@param		body	body	snapshotRequest{}	true	"Body"
//	@Produce	json
//	@Success	200	{object}	[]entity.StockSnapShot
//	@Failure	400	{object}	resp.Response{}
//	@Failure	500	{object}	resp.Response{}
//	@Router		/v1/stream/snapshot [put]
func (r *realTimeRoutes) getSnapshots(c *gin.Context) {
	p := snapshotRequest{}
	if err := c.ShouldBindJSON(&p); err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	if len(p.StockList) == 0 {
		resp.ErrorResponse(c, http.StatusBadRequest, "stock list is empty")
		return
	}
	snapshot, err := r.t.GetStockSnapshotByNumArr(p.StockList)
	if err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, snapshot)
}

func (r *realTimeRoutes) servePickStockWS(c *gin.Context) {
	pick.StartWSPickStock(c, r.t, false)
}

func (r *realTimeRoutes) servePickStockOddsWS(c *gin.Context) {
	pick.StartWSPickStock(c, r.t, true)
}
