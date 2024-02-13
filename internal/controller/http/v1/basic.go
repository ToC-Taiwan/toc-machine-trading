// Package v1 package v1
package v1

import (
	"net/http"

	"tmt/internal/controller/http/resp"
	"tmt/internal/entity"
	"tmt/internal/usecase"

	"github.com/gin-gonic/gin"
)

type basicRoutes struct {
	t usecase.Basic
}

func NewBasicRoutes(handler *gin.RouterGroup, t usecase.Basic) {
	r := &basicRoutes{t}

	h := handler.Group("/basic")
	{
		h.PUT("/stock", r.getStockDetail)
		h.GET("/usage/shioaji", r.getShioajiUsage)
	}
}

type stockDetailRequest struct {
	StockList []string `json:"stock_list"`
}

type stockDetailResponse struct {
	StockDetail []*entity.Stock `json:"stock_detail"`
}

// getStockDetail -.
//
//	@Tags		Basic V1
//	@Summary	Get stock detail by stock number
//	@security	JWT
//	@Accept		json
//	@Produce	json
//	@param		body	body		stockDetailRequest{}	true	"Body"
//	@Success	200		{object}	stockDetailResponse
//	@Failure	404		{object}	resp.Response{}
//	@Failure	500		{object}	resp.Response{}
//	@Router		/v1/basic/stock [put]
func (r *basicRoutes) getStockDetail(c *gin.Context) {
	p := stockDetailRequest{}
	if err := c.ShouldBindJSON(&p); err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	if len(p.StockList) == 0 {
		resp.ErrorResponse(c, http.StatusBadRequest, "stock list is empty")
		return
	}
	result := []*entity.Stock{}
	for _, stockNum := range p.StockList {
		stockDetail := r.t.GetStockDetail(stockNum)
		if stockDetail != nil {
			result = append(result, stockDetail)
		} else {
			result = append(result, &entity.Stock{
				Number: stockNum,
			})
		}
	}
	c.JSON(http.StatusOK, stockDetailResponse{
		StockDetail: result,
	})
}

// getShioajiUsage -.
//
//	@Tags		Basic V1
//	@Summary	Get shioaji usage
//	@security	JWT
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	entity.ShioajiUsage
//	@Failure	500	{object}	resp.Response{}
//	@Router		/v1/basic/usage/shioaji [get]
func (r *basicRoutes) getShioajiUsage(c *gin.Context) {
	usage, err := r.t.GetShioajiUsage()
	if err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, usage)
}
