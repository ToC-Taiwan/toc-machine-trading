package v1

import (
	"net/http"
	"strconv"

	"toc-machine-trading/internal/usecase"

	"github.com/gin-gonic/gin"
)

type orderRoutes struct {
	t usecase.Order
}

func newOrderRoutes(handler *gin.RouterGroup, t usecase.Order) {
	r := &orderRoutes{t}

	h := handler.Group("/order")
	{
		h.GET("/all", r.getAllOrder)
		h.GET("/balance", r.getAllTradeBalance)
		h.GET("/day-trade/forward", r.calculateDayTradeBalance)
	}
}

// @Summary     getAllOrder
// @Description getAllOrder
// @ID          getAllOrder
// @Tags  	    order
// @Accept      json
// @Produce     json
// @Success     200 {object} []entity.Order
// @Failure     500 {object} response
// @Router      /order/all [get]
func (r *orderRoutes) getAllOrder(c *gin.Context) {
	orderArr, err := r.t.GetAllOrder(c.Request.Context())
	if err != nil {
		log.Error(err)
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, orderArr)
}

// @Summary     getAllTradeBalance
// @Description getAllTradeBalance
// @ID          getAllTradeBalance
// @Tags  	    order
// @Accept      json
// @Produce     json
// @Success     200 {object} []entity.TradeBalance
// @Failure     500 {object} response
// @Router      /order/balance [get]
func (r *orderRoutes) getAllTradeBalance(c *gin.Context) {
	orderArr, err := r.t.GetAllTradeBalance(c.Request.Context())
	if err != nil {
		log.Error(err)
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, orderArr)
}

type dayTradeResult struct {
	Balance int64 `json:"balance"`
}

// @Summary     calculateDayTradeBalance
// @Description calculateDayTradeBalance
// @ID          calculateDayTradeBalance
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
func (r *orderRoutes) calculateDayTradeBalance(c *gin.Context) {
	buyPriceString := c.Request.Header.Get("buy_price")
	buyPrice, err := strconv.ParseFloat(buyPriceString, 64)
	if err != nil {
		log.Error(err)
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	buyQuantityString := c.Request.Header.Get("buy_quantity")
	buyQuantity, err := strconv.ParseInt(buyQuantityString, 10, 64)
	if err != nil {
		log.Error(err)
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	sellPriceString := c.Request.Header.Get("sell_price")
	sellPrice, err := strconv.ParseFloat(sellPriceString, 64)
	if err != nil {
		log.Error(err)
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	sellQuantityString := c.Request.Header.Get("sell_quantity")
	sellQuantity, err := strconv.ParseInt(sellQuantityString, 10, 64)
	if err != nil {
		log.Error(err)
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	pay := r.t.CalculateBuyCost(buyPrice, buyQuantity)
	payDiscount := r.t.CalculateTradeFee(buyPrice, buyQuantity)
	earning := r.t.CalculateSellCost(sellPrice, sellQuantity)
	earningDiscount := r.t.CalculateTradeFee(sellPrice, sellQuantity)

	c.JSON(http.StatusOK, dayTradeResult{
		Balance: -pay + payDiscount + earning + earningDiscount,
	})
}
