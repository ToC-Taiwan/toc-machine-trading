package v1

import (
	"net/http"

	"tmt/internal/controller/http/websocket"
	"tmt/internal/usecase"

	"github.com/gin-gonic/gin"
)

type streamRoutes struct {
	t usecase.Stream
	o usecase.Order
}

func newStreamRoutes(handler *gin.RouterGroup, t usecase.Stream, o usecase.Order) {
	r := &streamRoutes{t, o}

	h := handler.Group("/stream")
	{
		h.GET("/tse/snapshot", r.getTSESnapshot)

		h.GET("/future/switch", r.getFutureSwitchStatus)
		h.POST("/future/switch", r.modifyFutureSwitch)

		h.GET("/ws/pick-stock", r.servePickStockWS)
		h.GET("/ws/future", r.serveFutureWS)
	}
}

// @Summary     getTSESnapshot
// @Description getTSESnapshot
// @ID          getTSESnapshot
// @Tags  	    stream
// @Accept      json
// @Produce     json
// @Success     200 {object} entity.StockSnapShot
// @Failure     500 {object} response
// @Router      /stream/tse/snapshot [get]
func (r *streamRoutes) getTSESnapshot(c *gin.Context) {
	snapshot, err := r.t.GetTSESnapshot(c.Request.Context())
	if err != nil {
		log.Error(err)
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, snapshot)
}

// @Summary     getFutureSwitchStatus
// @Description getFutureSwitchStatus
// @ID          getFutureSwitchStatus
// @Tags  	    stream
// @Accept      json
// @Produce     json
// @Success     200 {object} bool
// @Failure     500 {object} response
// @Router      /stream/future/switch [get]
func (r *streamRoutes) getFutureSwitchStatus(c *gin.Context) {
	c.JSON(http.StatusOK, r.t.GetFutureTradeSwitchStatus(c.Request.Context()))
}

type futureSwitch struct {
	Switch bool `json:"switch"`
}

// @Summary     modifyFutureSwitch
// @Description modifyFutureSwitch
// @ID          modifyFutureSwitch
// @Tags  	    stream
// @Accept      json
// @Produce     json
// @Param body body futureSwitch{} true "Body"
// @Success     200
// @Failure     500 {object} response
// @Router      /stream/future/switch [post]
func (r *streamRoutes) modifyFutureSwitch(c *gin.Context) {
	body := &futureSwitch{}
	if err := c.BindJSON(body); err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	r.t.TurnFutureTradeSwitch(c.Request.Context(), body.Switch)
	c.JSON(http.StatusOK, nil)
}

func (r *streamRoutes) servePickStockWS(c *gin.Context) {
	wsRouter := websocket.NewWSRouter(r.t, r.o)
	wsRouter.Run(c, websocket.WSPickStock)
}

func (r *streamRoutes) serveFutureWS(c *gin.Context) {
	wsRouter := websocket.NewWSRouter(r.t, r.o)
	wsRouter.Run(c, websocket.WSFuture)
}
