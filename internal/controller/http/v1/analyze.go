// Package v1 package v1
package v1

import (
	"net/http"
	"sort"
	"time"

	"github.com/toc-taiwan/toc-machine-trading/internal/entity"
	"github.com/toc-taiwan/toc-machine-trading/internal/usecase"

	"github.com/gin-gonic/gin"
)

type analyzeRoutes struct {
	t usecase.Analyze
}

func NewAnalyzeRoutes(handler *gin.RouterGroup, t usecase.Analyze) {
	r := &analyzeRoutes{t}

	h := handler.Group("/analyze")
	{
		h.GET("/reborn", r.getRebornTargets)
	}
}

type reborn struct {
	Date   string         `json:"date"`
	Stocks []entity.Stock `json:"stocks"`
}

// getRebornTargets -.
//
//	@Tags		Analyze V1
//	@Summary	Get reborn targets
//	@security	JWT
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	[]reborn{}
//	@Router		/v1/analyze/reborn [get]
func (r *analyzeRoutes) getRebornTargets(c *gin.Context) {
	mapData := r.t.GetRebornMap(c.Request.Context())
	result := []reborn{}
	dateArr := []time.Time{}
	for date := range mapData {
		dateArr = append(dateArr, date)
	}

	sort.SliceStable(dateArr, func(i, j int) bool {
		return dateArr[i].After(dateArr[j])
	})

	for _, date := range dateArr {
		result = append(result, reborn{
			Date:   date.Format(entity.ShortTimeLayout),
			Stocks: mapData[date],
		})
	}
	c.JSON(http.StatusOK, result)
}
