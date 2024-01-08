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

// getAllRepoStock -.
//
//	@Tags		Basic V1
//	@Summary	Get all repo stock
//	@Accept		json
//	@Produce	json
//	@Success	200
//	@Failure	400	{object}	resp.Response{}
//	@Failure	500	{object}	resp.Response{}
//	@Router		/v1/fcm/announcement [get]
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
