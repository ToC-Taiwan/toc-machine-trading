// Package v1 package v1
package v1

import (
	"net/http"

	"tmt/internal/controller/http/resp"
	"tmt/internal/entity"
	"tmt/internal/usecase/usecase/basic"

	"github.com/gin-gonic/gin"
)

type basicRoutes struct {
	t basic.Basic
}

func NewBasicRoutes(handler *gin.RouterGroup, t basic.Basic) {
	r := &basicRoutes{t}

	h := handler.Group("/basic")
	{
		h.GET("/stock", r.getAllRepoStock)
		h.GET("/config", r.getAllConfig)
		h.GET("/usage/shioaji", r.getShioajiUsage)
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
// @Failure     404 {object} resp.Response{}
// @Failure     500 {object} resp.Response{}
// @Router      /basic/stock [get]
func (r *basicRoutes) getAllRepoStock(c *gin.Context) {
	stockNum := c.Query("num")
	stockDetail, err := r.t.GetAllRepoStock(c.Request.Context())
	if err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err.Error())
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
			resp.ErrorResponse(c, http.StatusNotFound, "stock not found")
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
// @Failure     500 {object} resp.Response{}
// @Router      /basic/config [get]
func (r *basicRoutes) getAllConfig(c *gin.Context) {
	c.JSON(http.StatusOK, r.t.GetConfig())
}

// @Summary     getShioajiUsage
// @Description getShioajiUsage
// @ID          getShioajiUsage
// @Tags  	    system
// @Accept      json
// @Produce     json
// @Success     200 {object} entity.ShioajiUsage
// @Failure     500 {object} resp.Response{}
// @Router      /basic/usage/shioaji [get]
func (r *basicRoutes) getShioajiUsage(c *gin.Context) {
	usage, err := r.t.GetShioajiUsage()
	if err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, usage)
}
