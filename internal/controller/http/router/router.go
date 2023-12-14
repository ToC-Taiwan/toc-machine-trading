// Package router implements routing paths. Each services in own file.
package router

import (
	"fmt"
	"net/http"
	"os"

	"tmt/docs"
	v1 "tmt/internal/controller/http/v1"
	"tmt/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const (
	prefixV1 = "/tmt/v1"
)

// Router -.
type Router struct {
	public  *gin.RouterGroup
	handler *gin.Engine
}

// NewRouter -.
//
//	@title			TOC MACHINE TRADING
//	@description	API docs for Auto Trade
//	@version		1.0.0
//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html
func NewRouter() *Router {
	gin.SetMode(os.Getenv("GIN_MODE"))
	docs.SwaggerInfo.BasePath = prefixV1

	handler := gin.New()
	handler.Use(gin.Recovery())

	// Swagger
	if os.Getenv("DISABLE_SWAGGER_HTTP_HANDLER") != "" {
		handler.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	} else {
		handler.Use(swaggerMiddleware())
	}

	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))
	handler.GET(fmt.Sprintf("%s/-/health", prefixV1), healthCheck)

	return &Router{
		public:  handler.Group(prefixV1),
		handler: handler,
	}
}

func (r *Router) GetHandler() *gin.Engine {
	return r.handler
}

// AddV1BasicRoutes -.
func (r *Router) AddV1BasicRoutes(basic usecase.Basic) *Router {
	v1.NewBasicRoutes(r.public, basic)
	return r
}

// AddV1AnalyzeRoutes -.
func (r *Router) AddV1AnalyzeRoutes(analyze usecase.Analyze) *Router {
	v1.NewAnalyzeRoutes(r.public, analyze)
	return r
}

// AddV1TargetRoutes -.
func (r *Router) AddV1TargetRoutes(target usecase.Target) *Router {
	v1.NewTargetRoutes(r.public, target)
	return r
}

// AddV1TradeRoutes -.
func (r *Router) AddV1TradeRoutes(trade usecase.Trade) *Router {
	v1.NewTradeRoutes(r.public, trade)
	return r
}

// AddV1HistoryRoutes -.
func (r *Router) AddV1HistoryRoutes(history usecase.History) *Router {
	v1.NewHistoryRoutes(r.public, history)
	return r
}

// AddV1RealTimeRoutes -.
func (r *Router) AddV1RealTimeRoutes(realTime usecase.RealTime, trade usecase.Trade, history usecase.History) *Router {
	v1.NewRealTimeRoutes(r.public, realTime, trade, history)
	return r
}

// healthCheck -.
//
//	@Summary		healthCheck
//	@Description	healthCheck
//	@ID				healthCheck
//	@Tags			healthCheck
//	@Accept			json
//	@Produce		json
//	@Success		200	{string}	string
//	@Router			/-/health [get]
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, "OK")
}

func swaggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		docs.SwaggerInfo.Host = c.Request.Host
		c.Next()
	}
}
