// Package v1 package v1
package v1

import (
	"net/http"

	"github.com/toc-taiwan/toc-machine-trading/internal/controller/http/resp"
	"github.com/toc-taiwan/toc-machine-trading/internal/controller/http/websocket/pick"
	"github.com/toc-taiwan/toc-machine-trading/internal/usecase"

	"github.com/gin-gonic/gin"
)

type realTimeRoutes struct {
	t     usecase.RealTime
	h     usecase.History
	basic usecase.Basic
}

func NewRealTimeRoutes(handler *gin.RouterGroup, basic usecase.Basic, t usecase.RealTime, history usecase.History) {
	r := &realTimeRoutes{
		t:     t,
		h:     history,
		basic: basic,
	}

	h := handler.Group("/stream")
	{
		h.PUT("/snapshot", r.getSnapshots)
		h.GET("/ws/pick-future/:code", r.servePickFutureWS)
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

func (r *realTimeRoutes) servePickFutureWS(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		resp.ErrorResponse(c, http.StatusBadRequest, "code is empty")
		return
	}

	pick.StartWSPickRealFuture(c, code, r.t, r.h, r.basic)
}
