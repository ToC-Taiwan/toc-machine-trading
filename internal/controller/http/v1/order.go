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
		h.GET("/all", r.getAllOrder)
		h.GET("/balance", r.getAllTradeBalance)
		h.GET("/balance/stock/last", r.getLastStockTradeBalance)
		h.GET("/balance/future/last", r.getLastFutureTradeBalance)

		h.GET("/date/:tradeday", r.getAllOrderByTradeDay)
		h.GET("/account/balance", r.getAccountBalance)
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
//	@Router		/v1/order/all [get]
func (r *orderRoutes) getAllOrder(c *gin.Context) {
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

// getAllOrderByTradeDay -.
//
//	@Tags		Order V1
//	@Summary	Get all order by trade day
//	@security	JWT
//	@Accept		json
//	@Produce	json
//	@param		tradeday	path		string	true	"tradeday"
//	@Success	200			{object}	futureOrders
//	@Failure	500			{object}	resp.Response{}
//	@Router		/v1/order/date/{tradeday} [get]
func (r *orderRoutes) getAllOrderByTradeDay(c *gin.Context) {
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

// getLastStockTradeBalance -.
//
//	@Tags		Order V1
//	@Summary	Get last stock trade balance
//	@security	JWT
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	entity.StockTradeBalance
//	@Failure	500	{object}	resp.Response{}
//	@Router		/v1/order/balance/stock/last [get]
func (r *orderRoutes) getLastStockTradeBalance(c *gin.Context) {
	balance, err := r.t.GetLastStockTradeBalance(c.Request.Context())
	if err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, balance)
}

// getLastFutureTradeBalance -.
//
//	@Tags		Order V1
//	@Summary	Get last future trade balance
//	@security	JWT
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	entity.FutureTradeBalance
//	@Failure	500	{object}	resp.Response{}
//	@Router		/v1/order/balance/future/last [get]
func (r *orderRoutes) getLastFutureTradeBalance(c *gin.Context) {
	balance, err := r.t.GetLastFutureTradeBalance(c.Request.Context())
	if err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, balance)
}

// getAccountBalance -.
//
//	@Tags		Account V1
//	@Summary	Get account balance
//	@security	JWT
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	entity.AccountBalance{}
//	@Failure	500	{object}	resp.Response{}
//	@Router		/v1/order/account/balance [get]
func (r *orderRoutes) getAccountBalance(c *gin.Context) {
	balance, err := r.t.GetAccountBalance(c.Request.Context())
	if err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, balance)
}
