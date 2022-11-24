package v1

import (
	"net/http"
	"strconv"
	"time"

	"tmt/internal/entity"
	"tmt/internal/usecase"
	"tmt/pkg/common"

	"github.com/gin-gonic/gin"
)

type historyRoutes struct {
	t usecase.History
}

func newHistoryRoutes(handler *gin.RouterGroup, t usecase.History) {
	r := &historyRoutes{t}

	h := handler.Group("/history")
	{
		h.GET("/day-kbar/:stock/:start_date/:interval", r.getKbarData)
		h.GET("/simulate/:interval", r.simulateFuture)
	}
}

// @Summary     getKbarData
// @Description getKbarData
// @ID          getKbarData
// @Tags  	    history
// @Accept      json
// @Produce     json
// @param stock path string true "stock"
// @param start_date path string true "start_date"
// @param interval path string true "interval"
// @success 200 {object} []entity.StockHistoryKbar
// @Failure     500 {object} response
// @Router      /history/day-kbar/{stock}/{start_date}/{interval} [get]
func (r *historyRoutes) getKbarData(c *gin.Context) {
	stockNum := c.Param("stock")
	interval, err := strconv.Atoi(c.Param("interval"))
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	startDate := c.Param("start_date")
	startDateTime, err := time.Parse(common.ShortTimeLayout, startDate)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if startDateTime.Equal(r.t.GetTradeDay()) {
		startDateTime = startDateTime.AddDate(0, 0, -1)
	}

	var result []entity.StockHistoryKbar
	for i := 0; i < interval; i++ {
		tmp := r.t.GetDayKbarByStockNumDate(stockNum, startDateTime)
		startDateTime = startDateTime.AddDate(0, 0, -1)
		if tmp == nil {
			continue
		}
		result = append(result, *tmp)
	}

	c.JSON(http.StatusOK, result)
}

// @Summary     simulateFuture
// @Description simulateFuture
// @ID          simulateFuture
// @Tags  	    history
// @Accept      json
// @Produce     json
// @param interval path string true "interval"
// @success 200 {object} trader.TradeBalance
// @Failure     500 {object} response
// @Router      /history/simulate/{interval} [get]
func (r *historyRoutes) simulateFuture(c *gin.Context) {
	interval, err := strconv.Atoi(c.Param("interval"))
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	cond := r.t.GetFutureTradeCond(interval)
	c.JSON(http.StatusOK, cond)
}
