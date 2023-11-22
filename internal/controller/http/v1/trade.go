// Package v1 package v1
package v1

import (
	"net/http"
	"strconv"
	"time"

	"tmt/global"
	"tmt/internal/controller/http/resp"
	"tmt/internal/entity"
	"tmt/internal/usecase"
	"tmt/internal/utils"

	"github.com/gin-gonic/gin"
)

type tradeRoutes struct {
	t usecase.Trade
}

func NewTradeRoutes(handler *gin.RouterGroup, t usecase.Trade) {
	r := &tradeRoutes{t}

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

		h.GET("/day-trade/forward", r.calculateForwardDayTradeBalance)
		h.GET("/day-trade/reverse", r.calculateReverseDayTradeBalance)

		h.GET("/account/balance", r.getAccountBalance)
	}
}

type allOrder struct {
	Stock  []*entity.StockOrder  `json:"stock"`
	Future []*entity.FutureOrder `json:"future"`
}

// getAllOrder -.
//
//	@Summary		getAllOrder
//	@Description	getAllOrder
//	@ID				getAllOrder
//	@Tags			order
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	allOrder
//	@Failure		500	{object}	resp.Response{}
//	@Router			/order/all [get]
func (r *tradeRoutes) getAllOrder(c *gin.Context) {
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
//	@Summary		getAllOrderByTradeDay
//	@Description	getAllOrderByTradeDay
//	@ID				getAllOrderByTradeDay
//	@Tags			order
//	@Accept			json
//	@Produce		json
//	@param			tradeday	path		string	true	"tradeday"
//	@Success		200			{object}	futureOrders
//	@Failure		500			{object}	resp.Response{}
//	@Router			/order/date/{tradeday} [get]
func (r *tradeRoutes) getAllOrderByTradeDay(c *gin.Context) {
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
//	@Summary		updateTradeBalanceByTradeDay
//	@Description	updateTradeBalanceByTradeDay
//	@ID				updateTradeBalanceByTradeDay
//	@Tags			order
//	@Accept			json
//	@Produce		json
//	@param			tradeday	path		string	true	"tradeday"
//	@Success		200			{object}	futureOrders
//	@Failure		500			{object}	resp.Response{}
//	@Router			/order/date/{tradeday} [put]
func (r *tradeRoutes) updateTradeBalanceByTradeDay(c *gin.Context) {
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
//	@Summary		manualInsertFutureOrder
//	@Description	manualInsertFutureOrder
//	@ID				manualInsertFutureOrder
//	@Tags			order
//	@Accept			json
//	@Produce		json
//	@param			body	body	manualInsertFutureOrderRequest{}	true	"Body"
//	@Success		200
//	@Failure		500	{object}	resp.Response{}
//	@Router			/order [post]
func (r *tradeRoutes) manualInsertFutureOrder(c *gin.Context) {
	body := &manualInsertFutureOrderRequest{}
	if err := c.BindJSON(body); err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	orderTime, err := time.ParseInLocation(global.LongTimeLayout, body.OrderTime, time.Local)
	if err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	order := &entity.FutureOrder{
		BaseOrder: entity.BaseOrder{
			OrderID:   utils.RandomASCIILowerOctdigitsString(8),
			Status:    entity.StatusFilled,
			Action:    body.Action,
			Price:     body.Price,
			Quantity:  body.Quantity,
			OrderTime: orderTime,
		},
		Code: body.Code,
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
//	@Summary		getAllTradeBalance
//	@Description	getAllTradeBalance
//	@ID				getAllTradeBalance
//	@Tags			order
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	tradeBalance
//	@Failure		500	{object}	resp.Response{}
//	@Router			/order/balance [get]
func (r *tradeRoutes) getAllTradeBalance(c *gin.Context) {
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
//	@Summary		getLastStockTradeBalance
//	@Description	getLastStockTradeBalance
//	@ID				getLastStockTradeBalance
//	@Tags			order
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	entity.StockTradeBalance
//	@Failure		500	{object}	resp.Response{}
//	@Router			/order/balance/stock/last [get]
func (r *tradeRoutes) getLastStockTradeBalance(c *gin.Context) {
	balance, err := r.t.GetLastStockTradeBalance(c.Request.Context())
	if err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, balance)
}

// getLastFutureTradeBalance -.
//
//	@Summary		getLastFutureTradeBalance
//	@Description	getLastFutureTradeBalance
//	@ID				getLastFutureTradeBalance
//	@Tags			order
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	entity.FutureTradeBalance
//	@Failure		500	{object}	resp.Response{}
//	@Router			/order/balance/future/last [get]
func (r *tradeRoutes) getLastFutureTradeBalance(c *gin.Context) {
	balance, err := r.t.GetLastFutureTradeBalance(c.Request.Context())
	if err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, balance)
}

type dayTradeResult struct {
	Balance int64 `json:"balance"`
}

// calculateForwardDayTradeBalance -.
//
//	@Summary		calculateForwardDayTradeBalance
//	@Description	calculateForwardDayTradeBalance
//	@ID				calculateForwardDayTradeBalance
//	@Tags			order
//	@accept			json
//	@produce		json
//	@param			buy_price		header		string	true	"buy_price"
//	@param			buy_quantity	header		string	true	"buy_quantity"
//	@param			sell_price		header		string	true	"sell_price"
//	@param			sell_quantity	header		string	true	"sell_quantity"
//	@success		200				{object}	dayTradeResult
//	@failure		500				{object}	resp.Response{}
//	@Router			/order/day-trade/forward [get]
func (r *tradeRoutes) calculateForwardDayTradeBalance(c *gin.Context) {
	buyPriceString := c.Request.Header.Get("buy_price")
	buyPrice, err := strconv.ParseFloat(buyPriceString, 64)
	if err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	buyQuantityString := c.Request.Header.Get("buy_quantity")
	buyQuantity, err := strconv.ParseInt(buyQuantityString, 10, 64)
	if err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	sellPriceString := c.Request.Header.Get("sell_price")
	sellPrice, err := strconv.ParseFloat(sellPriceString, 64)
	if err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	sellQuantityString := c.Request.Header.Get("sell_quantity")
	sellQuantity, err := strconv.ParseInt(sellQuantityString, 10, 64)
	if err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	pay := r.t.CalculateBuyCost(buyPrice, buyQuantity)
	payDiscount := r.t.CalculateTradeDiscount(buyPrice, buyQuantity)
	earning := r.t.CalculateSellCost(sellPrice, sellQuantity)
	earningDiscount := r.t.CalculateTradeDiscount(sellPrice, sellQuantity)

	c.JSON(http.StatusOK, dayTradeResult{
		Balance: -pay + payDiscount + earning + earningDiscount,
	})
}

// calculateReverseDayTradeBalance -.
//
//	@Summary		calculateReverseDayTradeBalance
//	@Description	calculateReverseDayTradeBalance
//	@ID				calculateReverseDayTradeBalance
//	@Tags			order
//	@accept			json
//	@produce		json
//	@param			sell_first_price	header		string	true	"sell_first_price"
//	@param			sell_first_quantity	header		string	true	"sell_first_quantity"
//	@param			buy_later_price		header		string	true	"buy_later_price"
//	@param			buy_later_quantity	header		string	true	"buy_later_quantity"
//	@success		200					{object}	dayTradeResult
//	@failure		500					{object}	resp.Response{}
//	@Router			/order/day-trade/reverse [get]
func (r *tradeRoutes) calculateReverseDayTradeBalance(c *gin.Context) {
	sellFirstPriceString := c.Request.Header.Get("sell_first_price")
	sellFirstPrice, err := strconv.ParseFloat(sellFirstPriceString, 64)
	if err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	sellFirstQuantityString := c.Request.Header.Get("sell_first_quantity")
	sellFirstQuantity, err := strconv.ParseInt(sellFirstQuantityString, 10, 64)
	if err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	buyLaterPriceString := c.Request.Header.Get("buy_later_price")
	buyLaterPrice, err := strconv.ParseFloat(buyLaterPriceString, 64)
	if err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	buyLaterQuantityString := c.Request.Header.Get("buy_later_quantity")
	buyLaterQuantity, err := strconv.ParseInt(buyLaterQuantityString, 10, 64)
	if err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	firstIn := r.t.CalculateSellCost(sellFirstPrice, sellFirstQuantity)
	firstInDiscount := r.t.CalculateTradeDiscount(sellFirstPrice, sellFirstQuantity)
	payLater := r.t.CalculateBuyCost(buyLaterPrice, buyLaterQuantity)
	payLaterDiscount := r.t.CalculateTradeDiscount(buyLaterPrice, buyLaterQuantity)

	c.JSON(http.StatusOK, dayTradeResult{
		Balance: firstIn + firstInDiscount - payLater + payLaterDiscount,
	})
}

// moveFutureOrderToLatestTradeDay -.
//
//	@Summary		moveFutureOrderToLatestTradeDay
//	@Description	moveFutureOrderToLatestTradeDay
//	@ID				moveFutureOrderToLatestTradeDay
//	@Tags			order
//	@Accept			json
//	@Produce		json
//	@param			order-id	path	string	true	"order-id"
//	@Success		200
//	@Failure		500	{object}	resp.Response{}
//	@Router			/order/future/{order-id} [patch]
func (r *tradeRoutes) moveFutureOrderToLatestTradeDay(c *gin.Context) {
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
//	@Summary		moveStockOrderToLatestTradeDay
//	@Description	moveStockOrderToLatestTradeDay
//	@ID				moveStockOrderToLatestTradeDay
//	@Tags			order
//	@Accept			json
//	@Produce		json
//	@param			order-id	path	string	true	"order-id"
//	@Success		200
//	@Failure		500	{object}	resp.Response{}
//	@Router			/order/stock/{order-id} [patch]
func (r *tradeRoutes) moveStockOrderToLatestTradeDay(c *gin.Context) {
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
//	@Summary		getAccountBalance
//	@Description	getAccountBalance
//	@ID				getAccountBalance
//	@Tags			account
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	accountSummary
//	@Failure		500	{object}	resp.Response{}
//	@Router			/order/account/balance [get]
func (r *tradeRoutes) getAccountBalance(c *gin.Context) {
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
