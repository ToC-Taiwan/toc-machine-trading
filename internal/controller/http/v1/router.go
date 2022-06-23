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

// NewRouter -.
// Swagger spec:
// @title       TOC MACHINE TRADING API
// @description Auto Trade on sinopac
// @version     1.0.0
func NewRouter(handler *gin.Engine, b usecase.Basic, a usecase.Analyze) {
	docs.SwaggerInfo.BasePath = "/v1"
	docs.SwaggerInfo.Host = "127.0.0.1:8080"

	// Options
	// handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

	// Prometheus metrics
	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Routers
	h := handler.Group("/v1")
	{
		newBasicRoutes(h, b)
		newAnalyzeRoutes(h, a)
	}
}
