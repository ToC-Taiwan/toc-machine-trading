// Package v1 package v1
package v1

import (
	"net/http"

	"tmt/internal/controller/http/resp"
	"tmt/internal/entity"
	"tmt/internal/usecase"

	"github.com/gin-gonic/gin"
)

type orderRoutes struct {
	t usecase.Trade
}

func NewOrderRoutes(handler *gin.RouterGroup, t usecase.Trade) {
	r := &orderRoutes{t}

	h := handler.Group("/order")
	{
		h.GET("/balance", r.getAllTradeBalance)
		h.GET("/future/all", r.getAllFutureOrder)
		h.POST("/future/:tradeday", r.getAllFutureOrderByTradeDay)
	}
}

type allOrder struct {
	Stock  []*entity.StockOrder  `json:"stock"`
	Future []*entity.FutureOrder `json:"future"`
}

// getAllOrder -.
//
//	@Tags		Order V1
//	@Summary	Get all order
//	@security	JWT
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	allOrder
//	@Failure	500	{object}	resp.Response{}
//	@Router		/v1/order/future/all [get]
func (r *orderRoutes) getAllFutureOrder(c *gin.Context) {
	stockOrderArr, err := r.t.GetAllStockOrder(c.Request.Context())
	if err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	futureOrderArr, err := r.t.GetAllFutureOrder(c.Request.Context())
	if err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, allOrder{
		Stock:  stockOrderArr,
		Future: futureOrderArr,
	})
}

type futureOrders struct {
	Orders []*entity.FutureOrder `json:"orders"`
}

// getAllFutureOrderByTradeDay -.
//
//	@Tags		Order V1
//	@Summary	Get all future order by trade day
//	@security	JWT
//	@Accept		json
//	@Produce	json
//	@param		tradeday	path		string	true	"tradeday"
//	@Success	200			{object}	futureOrders
//	@Failure	500			{object}	resp.Response{}
//	@Router		/v1/order/future/{tradeday} [post]
func (r *orderRoutes) getAllFutureOrderByTradeDay(c *gin.Context) {
	tradeDay := c.Param("tradeday")
	if tradeDay == "" {
		resp.ErrorResponse(c, http.StatusInternalServerError, "tradeday is empty")
		return
	}
	futureOrderArr, err := r.t.GetFutureOrderByTradeDay(c.Request.Context(), tradeDay)
	if err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, futureOrders{futureOrderArr})
}

type tradeBalance struct {
	Stock  []*entity.StockTradeBalance  `json:"stock"`
	Future []*entity.FutureTradeBalance `json:"future"`
}

// getAllTradeBalance -.
//
//	@Tags		Order V1
//	@Summary	Get all trade balance
//	@security	JWT
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	tradeBalance
//	@Failure	500	{object}	resp.Response{}
//	@Router		/v1/order/balance [get]
func (r *orderRoutes) getAllTradeBalance(c *gin.Context) {
	allStockArr, err := r.t.GetAllStockTradeBalance(c.Request.Context())
	if err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	allFutureArr, err := r.t.GetAllFutureTradeBalance(c.Request.Context())
	if err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, tradeBalance{
		Stock:  allStockArr,
		Future: allFutureArr,
	})
}
