package v1

import (
	"net/http"

	"tmt/cmd/config"
	"tmt/internal/entity"
	"tmt/internal/usecase"

	"github.com/gin-gonic/gin"
)

type basicRoutes struct {
	t usecase.Basic
}

func newBasicRoutes(handler *gin.RouterGroup, t usecase.Basic) {
	r := &basicRoutes{t}

	h := handler.Group("/basic")
	{
		h.GET("/stock", r.getAllRepoStock)
		h.GET("/stock/sinopac-to-repo", r.getAllSinopacStockAndUpdateRepo)
		h.GET("/config", r.getAllConfig)

		h.PUT("/system/terminate", r.terminateSinopac)
	}
}

type stockDetailResponse struct {
	StockDetail []*entity.Stock `json:"stock_detail"`
}

// @Summary     getAllRepoStock
// @Description getAllRepoStock
// @ID          getAllRepoStock
// @Tags  	    basic
// @Accept      json
// @Produce     json
// @Success     200 {object} stockDetailResponse
// @Failure     500 {object} response
// @Router      /basic/stock [get]
func (r *basicRoutes) getAllRepoStock(c *gin.Context) {
	stockDetail, err := r.t.GetAllRepoStock(c.Request.Context())
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, stockDetailResponse{
		StockDetail: stockDetail,
	})
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
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, stockDetailResponse{
		StockDetail: stockDetail,
	})
}

// @Summary     getAllConfig
// @Description getAllConfig
// @ID          getAllConfig
// @Tags  	    system
// @Accept      json
// @Produce     json
// @Success     200 {object} config.Config
// @Failure     500 {object} response
// @Router      /basic/config [get]
func (r *basicRoutes) getAllConfig(c *gin.Context) {
	cfg := config.GetConfig()
	c.JSON(http.StatusOK, cfg)
}

// @Summary     terminateSinopac
// @Description terminateSinopac
// @ID          terminateSinopac
// @Tags  	    system
// @Accept      json
// @Produce     json
// @Success     200
// @Failure     500 {object} response
// @Router      /basic/system/terminate [put]
func (r *basicRoutes) terminateSinopac(c *gin.Context) {
	err := r.t.TerminateSinopac(c.Request.Context())
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}
