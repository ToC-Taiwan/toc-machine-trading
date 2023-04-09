package v1

import (
	"net/http"

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
		h.GET("/config", r.getAllConfig)
		h.PUT("/system/terminate/sinopac", r.terminateSinopac)
		h.PUT("/system/terminate/fugle", r.terminateFugle)
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
// @Param 		num query string false "num"
// @Success     200 {object} stockDetailResponse
// @Failure     404 {object} response
// @Failure     500 {object} response
// @Router      /basic/stock [get]
func (r *basicRoutes) getAllRepoStock(c *gin.Context) {
	stockNum := c.Query("num")
	stockDetail, err := r.t.GetAllRepoStock(c.Request.Context())
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if stockNum != "" {
		result := []*entity.Stock{}
		for _, stock := range stockDetail {
			if stockNum == stock.Number {
				result = append(result, stock)
				break
			}
		}
		if len(result) == 0 {
			errorResponse(c, http.StatusNotFound, "stock not found")
			return
		}
		c.JSON(http.StatusOK, stockDetailResponse{
			StockDetail: result,
		})
	} else {
		c.JSON(http.StatusOK, stockDetailResponse{
			StockDetail: stockDetail,
		})
	}
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
	c.JSON(http.StatusOK, r.t.GetConfig())
}

// @Summary     terminateSinopac
// @Description terminateSinopac
// @ID          terminateSinopac
// @Tags  	    system
// @Accept      json
// @Produce     json
// @Success     200
// @Failure     500 {object} response
// @Router      /basic/system/terminate/sinopac [put]
func (r *basicRoutes) terminateSinopac(c *gin.Context) {
	err := r.t.TerminateSinopac()
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

// @Summary     terminateFugle
// @Description terminateFugle
// @ID          terminateFugle
// @Tags  	    system
// @Accept      json
// @Produce     json
// @Success     200
// @Failure     500 {object} response
// @Router      /basic/system/terminate/fugle [put]
func (r *basicRoutes) terminateFugle(c *gin.Context) {
	err := r.t.TerminateFugle()
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}
