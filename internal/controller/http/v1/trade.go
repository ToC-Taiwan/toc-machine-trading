// Package v1 package v1
package v1

import (
	"net/http"

	"tmt/internal/controller/http/resp"
	"tmt/internal/usecase"

	"github.com/gin-gonic/gin"
)

type tradeRoutes struct {
	t usecase.Trade
}

func NewTradeRoutes(handler *gin.RouterGroup, t usecase.Trade) {
	r := &tradeRoutes{t}

	h := handler.Group("/trade")
	{
		h.PUT("/stock/buy/odd", r.buyOddStock)
	}
}

type oddStockRequest struct {
	Num   string  `json:"num"`
	Price float64 `json:"price"`
	Share int64   `json:"share"`
}

type oddStockResponse struct {
	ID     string  `json:"id"`
	Status string  `json:"status"`
	Price  float64 `json:"price"`
	Share  int64   `json:"share"`
}

// buyOddStock -.
//
//	@Summary		buyOddStock
//	@Description	buyOddStock
//	@ID				buyOddStock
//	@Tags			trade
//	@Accept			json
//	@Produce		json
//	@param			body	body		oddStockRequest{}	true	"Body"
//	@Success		200		{object}	oddStockResponse{}
//	@Router			/v1/trade/stock/buy/odd [put]
func (r *tradeRoutes) buyOddStock(c *gin.Context) {
	p := oddStockRequest{}
	if err := c.ShouldBindJSON(&p); err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, status, err := r.t.BuyOddStock(p.Num, p.Price, p.Share)
	if err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, oddStockResponse{
		ID:     id,
		Status: status.String(),
		Price:  p.Price,
		Share:  p.Share,
	})
}
