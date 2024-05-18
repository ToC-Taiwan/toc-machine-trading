// Package v1 package v1
package v1

import (
	"net/http"

	"github.com/toc-taiwan/toc-machine-trading/internal/controller/http/auth"
	"github.com/toc-taiwan/toc-machine-trading/internal/controller/http/resp"
	"github.com/toc-taiwan/toc-machine-trading/internal/usecase"

	"github.com/gin-gonic/gin"
)

type tradeRoutes struct {
	t usecase.Trade
}

func NewTradeRoutes(handler *gin.RouterGroup, t usecase.Trade) {
	r := &tradeRoutes{t}

	h := handler.Group("/trade")
	{
		h.PUT("/stock/buy/odd", r.checkUserAuth, r.buyOddStock)
		h.PUT("/stock/sell/odd", r.checkUserAuth, r.sellOddStock)
		h.PUT("/cancel", r.checkUserAuth, r.cancelOrder)
		h.GET("/inventory/stock", r.getLatestInventoryStock)
	}
}

type cancelRequest struct {
	OrderID string `json:"order_id"`
}

type oddStockRequest struct {
	Num   string  `json:"num"`
	Price float64 `json:"price"`
	Share int64   `json:"share"`
}

type tradeResponse struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}

func (r *tradeRoutes) checkUserAuth(c *gin.Context) {
	if !r.t.IsAuthUser(auth.ExtractUsername(c)) {
		resp.ErrorResponse(c, http.StatusBadRequest, "user is not auth trader")
		return
	}
	c.Next()
}

// buyOddStock -.
//
//	@Tags		Trade V1
//	@Summary	Buy odd stock
//	@security	JWT
//	@Accept		json
//	@Produce	json
//	@param		body	body		oddStockRequest{}	true	"Body"
//	@Success	200		{object}	tradeResponse{}
//	@failure	401		{object}	resp.Response{}
//	@Router		/v1/trade/stock/buy/odd [put]
func (r *tradeRoutes) buyOddStock(c *gin.Context) {
	p := oddStockRequest{}
	if err := c.ShouldBindJSON(&p); err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	id, status, err := r.t.BuyOddStock(p.Num, p.Price, p.Share)
	if err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, tradeResponse{
		OrderID: id,
		Status:  status.String(),
	})
}

// sellOddStock -.
//
//	@Tags		Trade V1
//	@Summary	Sell odd stock
//	@security	JWT
//	@Accept		json
//	@Produce	json
//	@param		body	body		oddStockRequest{}	true	"Body"
//	@Success	200		{object}	tradeResponse{}
//	@failure	401		{object}	resp.Response{}
//	@Router		/v1/trade/stock/sell/odd [put]
func (r *tradeRoutes) sellOddStock(c *gin.Context) {
	p := oddStockRequest{}
	if err := c.ShouldBindJSON(&p); err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	id, status, err := r.t.SelloddStock(p.Num, p.Price, p.Share)
	if err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, tradeResponse{
		OrderID: id,
		Status:  status.String(),
	})
}

// cancelOrder -.
//
//	@Tags		Trade V1
//	@Summary	Cancel order
//	@security	JWT
//	@Accept		json
//	@Produce	json
//	@param		body	body	cancelRequest{}	true	"Body"
//	@Success	200
//	@failure	401	{object}	resp.Response{}
//	@Router		/v1/trade/cancel [put]
func (r *tradeRoutes) cancelOrder(c *gin.Context) {
	p := cancelRequest{}
	if err := c.ShouldBindJSON(&p); err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	_, _, err := r.t.CancelOrderByID(p.OrderID)
	if err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, nil)
}

// getLatestInventoryStock -.
//
//	@Tags		Trade V1
//	@Summary	Get latest inventory stock
//	@security	JWT
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	[]entity.InventoryStock{}
//	@failure	401	{object}	resp.Response{}
//	@Router		/v1/trade/inventory/stock [get]
func (r *tradeRoutes) getLatestInventoryStock(c *gin.Context) {
	stocks, err := r.t.GetLatestInventoryStock()
	if err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, stocks)
}
