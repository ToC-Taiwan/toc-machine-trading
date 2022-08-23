package v1

import (
	"net/http"
	"sort"
	"time"
	"tmt/internal/entity"
	"tmt/internal/usecase"
	"tmt/pkg/global"

	"github.com/gin-gonic/gin"
)

type analyzeRoutes struct {
	t usecase.Analyze
}

func newAnalyzeRoutes(handler *gin.RouterGroup, t usecase.Analyze) {
	r := &analyzeRoutes{t}

	h := handler.Group("/analyze")
	{
		h.GET("/reborn", r.getRebornTargets)
		h.GET("/simulate-historytick", r.startSimulateHistoryTick)
	}
}

type reborn struct {
	Date   string         `json:"date"`
	Stocks []entity.Stock `json:"stocks"`
}

// @Summary     getRebornTargets
// @Description getRebornTargets
// @ID          getRebornTargets
// @Tags  	    analyze
// @Accept      json
// @Produce     json
// @Success     200 {object} []reborn{}
// @Router      /analyze/reborn [get]
func (r *analyzeRoutes) getRebornTargets(c *gin.Context) {
	mapData := r.t.GetRebornMap(c.Request.Context())
	result := []reborn{}
	dateArr := []time.Time{}
	for date := range mapData {
		dateArr = append(dateArr, date)
	}

	sort.Slice(dateArr, func(i, j int) bool {
		return dateArr[i].After(dateArr[j])
	})

	for _, date := range dateArr {
		result = append(result, reborn{
			Date:   date.Format(global.ShortTimeLayout),
			Stocks: mapData[date],
		})
	}
	c.JSON(http.StatusOK, result)
}

// @Summary     startSimulateHistoryTick
// @Description startSimulateHistoryTick
// @ID          startSimulateHistoryTick
// @Tags  	    analyze
// @Accept      json
// @param use_default header bool true "use_default"
// @Produce     json
// @Success     200
// @Router      /analyze/simulate-historytick [get]
func (r *analyzeRoutes) startSimulateHistoryTick(c *gin.Context) {
	useDefaultString := c.Request.Header.Get("use_default")
	go r.t.SimulateOnHistoryTick(c.Request.Context(), useDefaultString == "true")
	c.JSON(http.StatusOK, gin.H{
		"message": "start simulate history tick, check result in logs",
	})
}
