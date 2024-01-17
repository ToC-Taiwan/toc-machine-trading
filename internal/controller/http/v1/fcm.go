package v1

import (
	"net/http"

	"tmt/internal/controller/http/resp"
	"tmt/internal/usecase"

	"github.com/gin-gonic/gin"
)

type fcmRoutes struct {
	t usecase.FCM
}

func NewFCMRoutes(handler *gin.RouterGroup, t usecase.FCM) {
	r := &fcmRoutes{t}

	h := handler.Group("/fcm")
	{
		h.POST("/announcement", r.announceMessage)
		h.POST("/push", r.pushMessage)
	}
}

type announceRequest struct {
	Message string `json:"message"`
}

// announceMessage -.
//
//	@Tags		FCM V1
//	@Summary	Announce message to all devices
//	@security	JWT
//	@Accept		json
//	@Produce	json
//	@param		body	body	announceRequest{}	true	"Body"
//	@Success	200
//	@Failure	400	{object}	resp.Response{}
//	@Failure	500	{object}	resp.Response{}
//	@Router		/v1/fcm/announcement [post]
func (r *fcmRoutes) announceMessage(c *gin.Context) {
	var req announceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	if err := r.t.AnnounceMessage(req.Message); err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}

type pushRequest struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// pushMessage -.
//
//	@Tags		FCM V1
//	@Summary	Push message to devices which has push token
//	@security	JWT
//	@Accept		json
//	@Produce	json
//	@param		body	body	pushRequest{}	true	"Body"
//	@Success	200
//	@Failure	400	{object}	resp.Response{}
//	@Failure	500	{object}	resp.Response{}
//	@Router		/v1/fcm/push [post]
func (r *fcmRoutes) pushMessage(c *gin.Context) {
	var req pushRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	if err := r.t.PushNotification(req.Title, req.Message); err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}
