package v1

import (
	"net/http"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/internal/usecase"
	"toc-machine-trading/pkg/logger"

	"github.com/gin-gonic/gin"
)

type stockRoutes struct {
	t usecase.Stock
}

func newStockRoutes(handler *gin.RouterGroup, t usecase.Stock) {
	r := &stockRoutes{t}

	h := handler.Group("/basic")
	{
		h.GET("/stock", r.getAllStockDetail)
	}
}

type stockDetailResponse struct {
	StockDetail []*entity.Stock `json:"stock_detail"`
}

// @Summary     getAllStockDetail
// @Description getAllStockDetail
// @ID          stock_detail
// @Tags  	    basic
// @Accept      json
// @Produce     json
// @Success     200 {object} stockDetailResponse
// @Failure     500 {object} response
// @Router      /basic/stock [get]
func (r *stockRoutes) getAllStockDetail(c *gin.Context) {
	stockDetail, err := r.t.GetAllStockDetail(c.Request.Context())
	if err != nil {
		logger.Get().Error(err)
		errorResponse(c, http.StatusInternalServerError, "grpc sinopac problems")
		return
	}

	c.JSON(http.StatusOK, stockDetailResponse{
		StockDetail: stockDetail,
	})
}
