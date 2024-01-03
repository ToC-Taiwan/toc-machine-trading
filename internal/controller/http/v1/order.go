// Package v1 package v1
package v1

import (
	"net/http"
	"time"

	"tmt/internal/controller/http/resp"
	"tmt/internal/entity"
	"tmt/internal/usecase"
	"tmt/pkg/utils"

	"github.com/gin-gonic/gin"
)

type orderRoutes struct {
	t usecase.Trade
}

func NewOrderRoutes(handler *gin.RouterGroup, t usecase.Trade) {
	r := &orderRoutes{t}

	h := handler.Group("/order")
	{
		h.POST("", r.manualInsertFutureOrder)

		h.GET("/all", r.getAllOrder)
		h.GET("/balance", r.getAllTradeBalance)
		h.GET("/balance/stock/last", r.getLastStockTradeBalance)
		h.GET("/balance/future/last", r.getLastFutureTradeBalance)

		h.GET("/date/:tradeday", r.getAllOrderByTradeDay)
		h.PUT("/date/:tradeday", r.updateTradeBalanceByTradeDay)

		h.PATCH("/stock/:order-id", r.moveStockOrderToLatestTradeDay)
		h.PATCH("/future/:order-id", r.moveFutureOrderToLatestTradeDay)

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
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	allOrder
//	@Failure	500	{object}	resp.Response{}
//	@Router		/v1/order/all [get]
func (r *orderRoutes) getAllOrder(c *gin.Context) {
	stockOrderArr, err := r.t.GetAllStockOrder(c.Request.Context())
	if err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	futureOrderArr, err := r.t.GetAllFutureOrder(c.Request.Context())
	if err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err.Error())
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
		resp.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, futureOrders{futureOrderArr})
}

// updateTradeBalanceByTradeDay -.
//
//	@Tags		Order V1
//	@Accept		json
//	@Produce	json
//	@param		tradeday	path		string	true	"tradeday"
//	@Success	200			{object}	futureOrders
//	@Failure	500			{object}	resp.Response{}
//	@Router		/v1/order/date/{tradeday} [put]
func (r *orderRoutes) updateTradeBalanceByTradeDay(c *gin.Context) {
	tradeDay := c.Param("tradeday")
	if tradeDay == "" {
		resp.ErrorResponse(c, http.StatusInternalServerError, "tradeday is empty")
		return
	}

	if err := r.t.UpdateTradeBalanceByTradeDay(c.Request.Context(), tradeDay); err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

type manualInsertFutureOrderRequest struct {
	Code      string             `json:"code"       binding:"required"`
	Price     float64            `json:"price"      binding:"required"`
	Quantity  int64              `json:"quantity"   binding:"required"`
	OrderTime string             `json:"order_time" binding:"required"`
	Action    entity.OrderAction `json:"action"     binding:"required"`
}

// manualInsertFutureOrder -.
//
//	@Tags		Order V1
//	@Accept		json
//	@Produce	json
//	@param		body	body	manualInsertFutureOrderRequest{}	true	"Body"
//	@Success	200
//	@Failure	500	{object}	resp.Response{}
//	@Router		/v1/order [post]
func (r *orderRoutes) manualInsertFutureOrder(c *gin.Context) {
	body := &manualInsertFutureOrderRequest{}
	if err := c.BindJSON(body); err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	orderTime, err := time.ParseInLocation(entity.LongTimeLayout, body.OrderTime, time.Local)
	if err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	order := &entity.FutureOrder{
		OrderDetail: entity.OrderDetail{
			OrderID:   utils.RandomASCIILowerOctdigitsString(8),
			Status:    entity.StatusFilled,
			Action:    body.Action,
			Price:     body.Price,
			OrderTime: orderTime,
		},
		Position: body.Quantity,
		Code:     body.Code,
	}

	if err := r.t.ManualInsertFutureOrder(c.Request.Context(), order); err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

type tradeBalance struct {
	Stock  []*entity.StockTradeBalance  `json:"stock"`
	Future []*entity.FutureTradeBalance `json:"future"`
}

// getAllTradeBalance -.
//
//	@Tags		Order V1
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	tradeBalance
//	@Failure	500	{object}	resp.Response{}
//	@Router		/v1/order/balance [get]
func (r *orderRoutes) getAllTradeBalance(c *gin.Context) {
	allStockArr, err := r.t.GetAllStockTradeBalance(c.Request.Context())
	if err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	allFutureArr, err := r.t.GetAllFutureTradeBalance(c.Request.Context())
	if err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err.Error())
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
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	entity.StockTradeBalance
//	@Failure	500	{object}	resp.Response{}
//	@Router		/v1/order/balance/stock/last [get]
func (r *orderRoutes) getLastStockTradeBalance(c *gin.Context) {
	balance, err := r.t.GetLastStockTradeBalance(c.Request.Context())
	if err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, balance)
}

// getLastFutureTradeBalance -.
//
//	@Tags		Order V1
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	entity.FutureTradeBalance
//	@Failure	500	{object}	resp.Response{}
//	@Router		/v1/order/balance/future/last [get]
func (r *orderRoutes) getLastFutureTradeBalance(c *gin.Context) {
	balance, err := r.t.GetLastFutureTradeBalance(c.Request.Context())
	if err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, balance)
}

// moveFutureOrderToLatestTradeDay -.
//
//	@Tags		Order V1
//	@Accept		json
//	@Produce	json
//	@param		order-id	path	string	true	"order-id"
//	@Success	200
//	@Failure	500	{object}	resp.Response{}
//	@Router		/v1/order/future/{order-id} [patch]
func (r *orderRoutes) moveFutureOrderToLatestTradeDay(c *gin.Context) {
	id := c.Param("order-id")
	if id == "" {
		resp.ErrorResponse(c, http.StatusInternalServerError, "order-id is empty")
		return
	}

	if e := r.t.MoveFutureOrderToLatestTradeDay(c.Request.Context(), id); e != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, e.Error())
		return
	}
	c.JSON(http.StatusOK, nil)
}

// moveStockOrderToLatestTradeDay -.
//
//	@Tags		Order V1
//	@Accept		json
//	@Produce	json
//	@param		order-id	path	string	true	"order-id"
//	@Success	200
//	@Failure	500	{object}	resp.Response{}
//	@Router		/v1/order/stock/{order-id} [patch]
func (r *orderRoutes) moveStockOrderToLatestTradeDay(c *gin.Context) {
	id := c.Param("order-id")
	if id == "" {
		resp.ErrorResponse(c, http.StatusInternalServerError, "order-id is empty")
		return
	}

	if e := r.t.MoveStockOrderToLatestTradeDay(c.Request.Context(), id); e != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, e.Error())
		return
	}
	c.JSON(http.StatusOK, nil)
}

type accountSummary struct {
	Balance []*entity.AccountBalance `json:"balance" yaml:"balance"`
	Total   float64                  `json:"total" yaml:"total"`
}

// getAccountBalance -.
//
//	@Tags		Account V1
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	accountSummary
//	@Failure	500	{object}	resp.Response{}
//	@Router		/v1/order/account/balance [get]
func (r *orderRoutes) getAccountBalance(c *gin.Context) {
	balance, err := r.t.GetAccountBalance(c.Request.Context())
	if err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var total float64
	for _, b := range balance {
		total += b.Balance
		total += b.TodayMargin
	}

	c.JSON(http.StatusOK, accountSummary{
		Balance: balance,
		Total:   total,
	})
}
