// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"toc-machine-trading/docs"
	"toc-machine-trading/internal/usecase"
	"toc-machine-trading/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var log = logger.Get()

// RouterV1 -.
type RouterV1 struct {
	g *gin.RouterGroup
}

// NewRouter -.
// @title       TOC MACHINE TRADING
// @description Auto Trade
// @version     0.0.1
func NewRouter(handler *gin.Engine) *RouterV1 {
	apiVersion := "/v1"

	docs.SwaggerInfo.BasePath = apiVersion
	docs.SwaggerInfo.Host = "127.0.0.1:8080"

	// handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

	// Prometheus metrics
	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	return &RouterV1{handler.Group(apiVersion)}
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
func (r *RouterV1) AddTargetRoutes(handler *gin.Engine, target *usecase.TargetUseCase) {
	newTargetRoutes(r.g, target)
}

// AddOrderRoutes -.
func (r *RouterV1) AddOrderRoutes(handler *gin.Engine, order *usecase.OrderUseCase) {
	newOrderRoutes(r.g, order)
}

// AddHistoryRoutes -.
func (r *RouterV1) AddHistoryRoutes(handler *gin.Engine, history *usecase.HistoryUseCase) {
	newHistoryRoutes(r.g, history)
}

// AddStreamRoutes -.
func (r *RouterV1) AddStreamRoutes(handler *gin.Engine, stream *usecase.StreamUseCase) {
	newStreamRoutes(r.g, stream)
}
