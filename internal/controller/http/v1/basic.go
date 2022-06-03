package v1

import (
	"net/http"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/internal/usecase"
	"toc-machine-trading/pkg/logger"

	"github.com/gin-gonic/gin"
)

type basicRoutes struct {
	t usecase.Basic
}

func newBasicRoutes(handler *gin.RouterGroup, t usecase.Basic) {
	r := &basicRoutes{t}

	h := handler.Group("/basic")
	{
		h.GET("/stock/sinopac-to-repo", r.getAllSinopacStockAndUpdateRepo)
		h.GET("/stock/repo", r.getAllRepoStock)
	}
}

type stockDetailResponse struct {
	StockDetail []*entity.Stock `json:"stock_detail"`
}

// @Summary     getAllSinopacStockAndUpdateRepo
// @Description getAllSinopacStockAndUpdateRepo
// @ID          getAllSinopacStockAndUpdateRepo
// @Tags  	    basic
// @Accept      json
// @Produce     json
// @Success     200 {object} stockDetailResponse
// @Failure     500 {object} response
// @Router      /basic/stock/sinopac-to-repo [get]
func (r *basicRoutes) getAllSinopacStockAndUpdateRepo(c *gin.Context) {
	stockDetail, err := r.t.GetAllSinopacStockAndUpdateRepo(c.Request.Context())
	if err != nil {
		logger.Get().Error(err)
		errorResponse(c, http.StatusInternalServerError, "sinopac problems")
		return
	}

	c.JSON(http.StatusOK, stockDetailResponse{
		StockDetail: stockDetail,
	})
}

// @Summary     getAllRepoStock
// @Description getAllRepoStock
// @ID          getAllRepoStock
// @Tags  	    basic
// @Accept      json
// @Produce     json
// @Success     200 {object} stockDetailResponse
// @Failure     500 {object} response
// @Router      /basic/stock/repo [get]
func (r *basicRoutes) getAllRepoStock(c *gin.Context) {
	stockDetail, err := r.t.GetAllRepoStock(c.Request.Context())
	if err != nil {
		logger.Get().Error(err)
		errorResponse(c, http.StatusInternalServerError, "repo problems")
		return
	}

	c.JSON(http.StatusOK, stockDetailResponse{
		StockDetail: stockDetail,
	})
}
