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
		h.PUT("/stock/buy/lot", r.buyLotStock)
	}
}

type oddStockRequest struct {
	Num   string  `json:"num"`
	Price float64 `json:"price"`
	Share int64   `json:"share"`
}

type lotStockRequest struct {
	Num   string  `json:"num"`
	Price float64 `json:"price"`
	Lot   int64   `json:"lot"`
}

type tradeResponse struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
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
//	@Success		200		{object}	tradeResponse{}
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
	c.JSON(http.StatusOK, tradeResponse{
		OrderID: id,
		Status:  status.String(),
	})
}

// buyLotStock -.
//
//	@Summary		buyLotStock
//	@Description	buyLotStock
//	@ID				buyLotStock
//	@Tags			trade
//	@Accept			json
//	@Produce		json
//	@param			body	body		lotStockRequest{}	true	"Body"
//	@Success		200		{object}	tradeResponse{}
//	@Router			/v1/trade/stock/buy/lot [put]
func (r *tradeRoutes) buyLotStock(c *gin.Context) {
	p := lotStockRequest{}
	if err := c.ShouldBindJSON(&p); err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, status, err := r.t.BuyLotStock(p.Num, p.Price, p.Lot)
	if err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, tradeResponse{
		OrderID: id,
		Status:  status.String(),
	})
}
