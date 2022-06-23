package v1

import (
	"net/http"
	"sort"
	"time"

	"toc-machine-trading/internal/entity"
	"toc-machine-trading/internal/usecase"
	"toc-machine-trading/pkg/global"

	"github.com/gin-gonic/gin"
)

type analyzeRoutes struct {
	t usecase.Analyze
}

func newAnalyzeRoutes(handler *gin.RouterGroup, t usecase.Analyze) {
	r := &analyzeRoutes{t}

	h := handler.Group("/analyze")
	{
		h.GET("/below-quater", r.getQuaterTargets)
	}
}

// BelowQuaterMA BelowQuaterMA
type BelowQuaterMA struct {
	Date   string         `json:"date"`
	Stocks []entity.Stock `json:"stocks"`
}

// @Summary     getQuaterTargets
// @Description getQuaterTargets
// @ID          getQuaterTargets
// @Tags  	    analyze
// @Accept      json
// @Produce     json
// @Success     200 {object} []BelowQuaterMA{}
// @Router      /analyze/below-quater [get]
func (r *analyzeRoutes) getQuaterTargets(c *gin.Context) {
	mapData := r.t.GetBelowQuaterMap(c.Request.Context())
	result := []BelowQuaterMA{}
	dateArr := []time.Time{}
	for date := range mapData {
		dateArr = append(dateArr, date)
	}

	sort.Slice(dateArr, func(i, j int) bool {
		return dateArr[i].After(dateArr[j])
	})

	for _, date := range dateArr {
		result = append(result, BelowQuaterMA{
			Date:   date.Format(global.ShortTimeLayout),
			Stocks: mapData[date],
		})
	}
	c.JSON(http.StatusOK, result)
}
