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
		h.GET("/", r.getAllOrder)
	}
}

// @Summary     getAllOrder
// @Description getAllOrder
// @ID          getAllOrder
// @Tags  	    order
// @Accept      json
// @Produce     json
// @Router      /order [get]
func (r *orderRoutes) getAllOrder(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}
