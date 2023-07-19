// Package v1 package v1
package v1

import (
	"net/http"
	"strconv"
	"time"

	"tmt/cmd/config"
	"tmt/global"
	"tmt/internal/controller/http/resp"
	"tmt/internal/entity"
	"tmt/internal/usecase"
	"tmt/internal/usecase/module/simulator"

	"github.com/gin-gonic/gin"
)

type historyRoutes struct {
	t usecase.History
}

func NewHistoryRoutes(handler *gin.RouterGroup, t usecase.History) {
	r := &historyRoutes{t}

	h := handler.Group("/history")
	{
		h.GET("/day-kbar/:stock/:start_date/:interval", r.getKbarData)
		h.POST("/simulate/future", r.simulateFuture)
		h.POST("/simulate/future/auto", r.simulateFutureAuto)
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
// @Failure     500 {object} resp.Response{}
// @Router      /history/day-kbar/{stock}/{start_date}/{interval} [get]
func (r *historyRoutes) getKbarData(c *gin.Context) {
	stockNum := c.Param("stock")
	interval, err := strconv.Atoi(c.Param("interval"))
	if err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	startDate := c.Param("start_date")
	startDateTime, err := time.Parse(global.ShortTimeLayout, startDate)
	if err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, err.Error())
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
// @param 		need_detail header bool false "It accepts 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False."
// @param       body body config.TradeFuture{} true "TradeFuture"
// @success     200 {object} simulator.SimulateBalance
// @Failure     500 {object} resp.Response{}
// @Router      /history/simulate/future [post]
func (r *historyRoutes) simulateFuture(c *gin.Context) {
	body := &config.TradeFuture{}
	if err := c.BindJSON(body); err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var needDetail bool
	result := r.t.SimulateOne(body)
	if h := c.GetHeader("need_detail"); h != "" {
		if need, err := strconv.ParseBool(h); err != nil {
			resp.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		} else if need {
			needDetail = true
		}
	}

	if !needDetail {
		result.ForwardOrder = nil
		result.ReverseOrder = nil
	}
	c.JSON(http.StatusOK, result)
}

// @Summary     simulateFutureAuto
// @Description simulateFutureAuto
// @ID          simulateFutureAuto
// @Tags  	    history
// @Accept      json
// @Produce     json
// @success     200
// @Failure     500 {object} resp.Response{}
// @Router      /history/simulate/future/auto [post]
func (r *historyRoutes) simulateFutureAuto(c *gin.Context) {
	go r.t.SimulateMulti(simulator.GenerateCond())
	c.JSON(http.StatusOK, nil)
}
