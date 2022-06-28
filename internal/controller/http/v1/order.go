package v1

import (
	"net/http"

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
