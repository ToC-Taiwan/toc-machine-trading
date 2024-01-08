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
	}
}

type announcementRequest struct {
	Message string `json:"message"`
}

// announceMessage -.
//
//	@Tags		FCM V1
//	@Summary	Announce message to all devices
//	@security	JWT
//	@Accept		json
//	@Produce	json
//	@param		body	body	announcementRequest{}	true	"Body"
//	@Success	200
//	@Failure	400	{object}	resp.Response{}
//	@Failure	500	{object}	resp.Response{}
//	@Router		/v1/fcm/announcement [post]
func (r *fcmRoutes) announceMessage(c *gin.Context) {
	var req announcementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := r.t.AnnounceMessage(req.Message); err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}
