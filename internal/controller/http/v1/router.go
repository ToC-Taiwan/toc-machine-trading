// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"fmt"
	"net/http"

	"tmt/docs"
	"tmt/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RouterV1 -.
type RouterV1 struct {
	g *gin.RouterGroup
}

// NewRouter -.
// @title       TOC MACHINE TRADING
// @description API docs for Auto Trade
// @version     0.0.1
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
func NewRouter(handler *gin.Engine) *RouterV1 {
	prefix := "/tmt/v1"

	docs.SwaggerInfo.BasePath = prefix
	docs.SwaggerInfo.Host = "127.0.0.1:26670"

	// handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

	// Prometheus metrics
	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Prometheus metrics
	handler.GET(fmt.Sprintf("%s/-/health", prefix), healthCheck)

	return &RouterV1{
		g: handler.Group(prefix),
	}
}

// AddBasicRoutes -.
func (r *RouterV1) AddBasicRoutes(handler *gin.Engine, basic usecase.Basic) {
	newBasicRoutes(r.g, basic)
}

// AddAnalyzeRoutes -.
func (r *RouterV1) AddAnalyzeRoutes(handler *gin.Engine, analyze usecase.Analyze) {
	newAnalyzeRoutes(r.g, analyze)
}

// AddTargetRoutes -.
func (r *RouterV1) AddTargetRoutes(handler *gin.Engine, target usecase.Target) {
	newTargetRoutes(r.g, target)
}

// AddTradeRoutes -.
func (r *RouterV1) AddTradeRoutes(handler *gin.Engine, trade usecase.Trade) {
	newTradeRoutes(r.g, trade)
}

// AddHistoryRoutes -.
func (r *RouterV1) AddHistoryRoutes(handler *gin.Engine, history usecase.History) {
	newHistoryRoutes(r.g, history)
}

// AddRealTimeRoutes -.
func (r *RouterV1) AddRealTimeRoutes(handler *gin.Engine, realTime usecase.RealTime, trade usecase.Trade, history usecase.History) {
	newRealTimeRoutes(r.g, realTime, trade, history)
}

// @Summary     healthCheck
// @Description healthCheck
// @ID          healthCheck
// @Tags  	    healthCheck
// @Accept      json
// @Produce     json
// @Success     200 {string} string
// @Router      /-/health [get]
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, "OK")
}
