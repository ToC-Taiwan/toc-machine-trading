package v1

import (
	"net/http"
	"strconv"
	"time"

	"tmt/internal/entity"
	"tmt/internal/usecase"
	"tmt/pkg/common"
	"tmt/pkg/utils"

	"github.com/gin-gonic/gin"
)

type orderRoutes struct {
	t usecase.Order
}

func newOrderRoutes(handler *gin.RouterGroup, t usecase.Order) {
	r := &orderRoutes{t}

	h := handler.Group("/order")
	{
		h.POST("", r.manualInsertFutureOrder)
		h.GET("/all", r.getAllOrder)
		h.GET("/date/:tradeday", r.getAllOrderByTradeDay)
		h.GET("/balance", r.getAllTradeBalance)

		h.GET("/day-trade/forward", r.calculateForwardDayTradeBalance)
		h.GET("/day-trade/reverse", r.calculateReverseDayTradeBalance)

		h.PUT("/status/update", r.askOrderUpdate)
	}
}

type allOrder struct {
	Stock  []*entity.StockOrder  `json:"stock"`
	Future []*entity.FutureOrder `json:"future"`
}

// @Summary     getAllOrder
// @Description getAllOrder
// @ID          getAllOrder
// @Tags  	    order
// @Accept      json
// @Produce     json
// @Success     200 {object} allOrder
// @Failure     500 {object} response
// @Router      /order/all [get]
func (r *orderRoutes) getAllOrder(c *gin.Context) {
	stockOrderArr, err := r.t.GetAllStockOrder(c.Request.Context())
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	futureOrderArr, err := r.t.GetAllFutureOrder(c.Request.Context())
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
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

// @Summary     getAllOrderByTradeDay
// @Description getAllOrderByTradeDay
// @ID          getAllOrderByTradeDay
// @Tags  	    order
// @Accept      json
// @Produce     json
// @param tradeday path string true "tradeday"
// @Success     200 {object} futureOrders
// @Failure     500 {object} response
// @Router      /order/date/{tradeday} [get]
func (r *orderRoutes) getAllOrderByTradeDay(c *gin.Context) {
	tradeDay := c.Param("tradeday")
	if tradeDay == "" {
		errorResponse(c, http.StatusInternalServerError, "tradeday is empty")
		return
	}

	futureOrderArr, err := r.t.GetFutureOrderByTradeDay(c.Request.Context(), tradeDay)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, futureOrders{futureOrderArr})
}

type manualInsertFutureOrderRequest struct {
	Code      string             `json:"code"       binding:"required"`
	Price     float64            `json:"price"      binding:"required"`
	Quantity  int64              `json:"quantity"   binding:"required"`
	OrderTime string             `json:"order_time" binding:"required"`
	Action    entity.OrderAction `json:"action"     binding:"required"`
}

// @Summary     manualInsertFutureOrder
// @Description manualInsertFutureOrder
// @ID          manualInsertFutureOrder
// @Tags  	    order
// @Accept      json
// @Produce     json
// @param body body manualInsertFutureOrderRequest{} true "Body"
// @Success     200
// @Failure     500 {object} response
// @Router      /order [post]
func (r *orderRoutes) manualInsertFutureOrder(c *gin.Context) {
	body := &manualInsertFutureOrderRequest{}
	if err := c.BindJSON(body); err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	orderTime, err := time.ParseInLocation(common.LongTimeLayout, body.OrderTime, time.Local)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	order := &entity.FutureOrder{
		BaseOrder: entity.BaseOrder{
			GroupID:   "-",
			OrderID:   utils.RandomASCIILowerOctdigitsString(8),
			Status:    entity.StatusFilled,
			Action:    body.Action,
			Price:     body.Price,
			Quantity:  body.Quantity,
			TradeTime: orderTime,
			TickTime:  orderTime,
			OrderTime: orderTime,
		},
		Code:   body.Code,
		Manual: true,
	}

	if err := r.t.ManualInsertFutureOrder(c.Request.Context(), order); err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

type tradeBalance struct {
	Stock  []*entity.StockTradeBalance  `json:"stock"`
	Future []*entity.FutureTradeBalance `json:"future"`
}

// @Summary     getAllTradeBalance
// @Description getAllTradeBalance
// @ID          getAllTradeBalance
// @Tags  	    order
// @Accept      json
// @Produce     json
// @Success     200 {object} tradeBalance
// @Failure     500 {object} response
// @Router      /order/balance [get]
func (r *orderRoutes) getAllTradeBalance(c *gin.Context) {
	allStockArr, err := r.t.GetAllStockTradeBalance(c.Request.Context())
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	allFutureArr, err := r.t.GetAllFutureTradeBalance(c.Request.Context())
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, tradeBalance{
		Stock:  allStockArr,
		Future: allFutureArr,
	})
}

type dayTradeResult struct {
	Balance int64 `json:"balance"`
}

// @Summary     calculateForwardDayTradeBalance
// @Description calculateForwardDayTradeBalance
// @ID          calculateForwardDayTradeBalance
// @Tags  	    order
// @accept json
// @produce json
// @param buy_price header string true "buy_price"
// @param buy_quantity header string true "buy_quantity"
// @param sell_price header string true "sell_price"
// @param sell_quantity header string true "sell_quantity"
// @success 200 {object} dayTradeResult
// @failure 500 {object} response
// @Router /order/day-trade/forward [get]
func (r *orderRoutes) calculateForwardDayTradeBalance(c *gin.Context) {
	buyPriceString := c.Request.Header.Get("buy_price")
	buyPrice, err := strconv.ParseFloat(buyPriceString, 64)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	buyQuantityString := c.Request.Header.Get("buy_quantity")
	buyQuantity, err := strconv.ParseInt(buyQuantityString, 10, 64)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	sellPriceString := c.Request.Header.Get("sell_price")
	sellPrice, err := strconv.ParseFloat(sellPriceString, 64)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	sellQuantityString := c.Request.Header.Get("sell_quantity")
	sellQuantity, err := strconv.ParseInt(sellQuantityString, 10, 64)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
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

// @Summary     calculateReverseDayTradeBalance
// @Description calculateReverseDayTradeBalance
// @ID          calculateReverseDayTradeBalance
// @Tags  	    order
// @accept json
// @produce json
// @param sell_first_price header string true "sell_first_price"
// @param sell_first_quantity header string true "sell_first_quantity"
// @param buy_later_price header string true "buy_later_price"
// @param buy_later_quantity header string true "buy_later_quantity"
// @success 200 {object} dayTradeResult
// @failure 500 {object} response
// @Router /order/day-trade/reverse [get]
func (r *orderRoutes) calculateReverseDayTradeBalance(c *gin.Context) {
	sellFirstPriceString := c.Request.Header.Get("sell_first_price")
	sellFirstPrice, err := strconv.ParseFloat(sellFirstPriceString, 64)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	sellFirstQuantityString := c.Request.Header.Get("sell_first_quantity")
	sellFirstQuantity, err := strconv.ParseInt(sellFirstQuantityString, 10, 64)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	buyLaterPriceString := c.Request.Header.Get("buy_later_price")
	buyLaterPrice, err := strconv.ParseFloat(buyLaterPriceString, 64)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	buyLaterQuantityString := c.Request.Header.Get("buy_later_quantity")
	buyLaterQuantity, err := strconv.ParseInt(buyLaterQuantityString, 10, 64)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
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

// @Summary     askOrderUpdate
// @Description askOrderUpdate
// @ID          askOrderUpdate
// @Tags  	    order
// @Accept      json
// @Produce     json
// @Success     200
// @Failure     500 {object} response
// @Router      /order/status/update [put]
func (r *orderRoutes) askOrderUpdate(c *gin.Context) {
	if err := r.t.AskOrderUpdate(); err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, nil)
}
