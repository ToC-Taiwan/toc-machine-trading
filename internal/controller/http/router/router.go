// Package router implements routing paths. Each services in own file.
package router

import (
	"fmt"
	"net/http"

	"tmt/docs"
	"tmt/internal/controller/http/auth"
	v1 "tmt/internal/controller/http/v1"
	"tmt/internal/usecase"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	prefix = "/tmt"
)

// Router -.
type Router struct {
	v1Public    *gin.RouterGroup
	v1Private   *gin.RouterGroup
	rootHandler *gin.Engine
	jwtHandler  *jwt.GinJWTMiddleware
}

var swagHandler gin.HandlerFunc

// NewRouter -.
//
//	@title						TMT OpenAPI
//	@description				Toc Machine Trading's API docs
//	@version					2.5.0
//	@securityDefinitions.apikey	JWT
//	@in							header
//	@name						Authorization
//	@license.name				GPLv3
//	@license.url				https://www.gnu.org/licenses/gpl-3.0.html#license-text
func NewRouter(system usecase.System) *Router {
	g := gin.New()
	g.Use(gin.Recovery())
	g.GET("/metrics", gin.WrapH(promhttp.Handler()))
	g.GET("/-/health", healthCheck)
	g.LoadHTMLGlob("templates/*")

	if swagHandler != nil {
		docs.SwaggerInfo.BasePath = prefix
		g.GET("/docs/*any", swagHandler)
		g.Use(swaggerMiddleware())
	}

	jwtHandler, err := auth.NewAuthMiddleware(system)
	if err != nil {
		panic(err)
	}

	v1Public := g.Group(fmt.Sprintf("%s/v1", prefix))
	v1Private := g.Group(fmt.Sprintf("%s/v1", prefix))
	v1Private.Use(jwtHandler.MiddlewareFunc())
	return &Router{
		v1Public:    v1Public,
		v1Private:   v1Private,
		rootHandler: g,
		jwtHandler:  jwtHandler,
	}
}

func (r *Router) GetHandler() *gin.Engine {
	return r.rootHandler
}

// AddV1UserRoutes -.
func (r *Router) AddV1UserRoutes(system usecase.System) *Router {
	v1.NewUserRoutes(r.v1Public, r.jwtHandler, system)
	return r
}

// AddV1BasicRoutes -.
func (r *Router) AddV1BasicRoutes(basic usecase.Basic) *Router {
	v1.NewBasicRoutes(r.v1Public, basic)
	return r
}

// AddV1AnalyzeRoutes -.
func (r *Router) AddV1AnalyzeRoutes(analyze usecase.Analyze) *Router {
	v1.NewAnalyzeRoutes(r.v1Public, analyze)
	return r
}

// AddV1TargetRoutes -.
func (r *Router) AddV1TargetRoutes(target usecase.Target) *Router {
	v1.NewTargetRoutes(r.v1Public, target)
	return r
}

// AddV1OrderRoutes -.
func (r *Router) AddV1OrderRoutes(trade usecase.Trade) *Router {
	v1.NewOrderRoutes(r.v1Public, trade)
	return r
}

// AddV1TradeRoutes -.
func (r *Router) AddV1TradeRoutes(trade usecase.Trade) *Router {
	v1.NewTradeRoutes(r.v1Private, trade)
	return r
}

// AddV1HistoryRoutes -.
func (r *Router) AddV1HistoryRoutes(history usecase.History) *Router {
	v1.NewHistoryRoutes(r.v1Public, history)
	return r
}

// AddV1RealTimeRoutes -.
func (r *Router) AddV1RealTimeRoutes(realTime usecase.RealTime, trade usecase.Trade, history usecase.History) *Router {
	v1.NewRealTimeRoutes(r.v1Public, realTime, trade, history)
	return r
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, "OK")
}

func swaggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		docs.SwaggerInfo.Host = c.Request.Host
		c.Next()
	}
}
