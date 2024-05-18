// Package router implements routing paths. Each services in own file.
package router

import (
	"fmt"

	"github.com/toc-taiwan/toc-machine-trading/docs"
	"github.com/toc-taiwan/toc-machine-trading/internal/controller/http/auth"
	v1 "github.com/toc-taiwan/toc-machine-trading/internal/controller/http/v1"
	"github.com/toc-taiwan/toc-machine-trading/internal/usecase"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	prefix = "/tmt"
)

// Router -.
type Router struct {
	rootHandler *gin.Engine
	v1Group     *gin.RouterGroup
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

	if swagHandler != nil {
		docs.SwaggerInfo.BasePath = prefix
		g.GET("/docs/*any", swagHandler)
		g.Use(swaggerMiddleware())
	}

	jwtHandler, err := auth.NewAuthMiddleware(system)
	if err != nil {
		panic(err)
	}

	v1Prefix := fmt.Sprintf("%s/v1", prefix)

	v1Public := g.Group(v1Prefix)
	v1Private := g.Group(v1Prefix)

	v1Private.Use(jwtHandler.MiddlewareFunc())
	v1.NewUserRoutes(v1Public, v1Private, jwtHandler, system)

	return &Router{
		rootHandler: g,
		v1Group:     v1Private,
		jwtHandler:  jwtHandler,
	}
}

func (r *Router) GetHandler() *gin.Engine {
	return r.rootHandler
}

func (r *Router) AddV1FCMRoutes(fcm usecase.FCM) *Router {
	v1.NewFCMRoutes(r.v1Group, fcm)
	return r
}

func (r *Router) AddV1TradeRoutes(trade usecase.Trade) *Router {
	v1.NewTradeRoutes(r.v1Group, trade)
	return r
}

func (r *Router) AddV1BasicRoutes(basic usecase.Basic) *Router {
	v1.NewBasicRoutes(r.v1Group, basic)
	return r
}

func (r *Router) AddV1AnalyzeRoutes(analyze usecase.Analyze) *Router {
	v1.NewAnalyzeRoutes(r.v1Group, analyze)
	return r
}

func (r *Router) AddV1TargetRoutes(target usecase.Target) *Router {
	v1.NewTargetRoutes(r.v1Group, target)
	return r
}

func (r *Router) AddV1OrderRoutes(trade usecase.Trade) *Router {
	v1.NewOrderRoutes(r.v1Group, trade)
	return r
}

func (r *Router) AddV1HistoryRoutes(history usecase.History) *Router {
	v1.NewHistoryRoutes(r.v1Group, history)
	return r
}

func (r *Router) AddV1RealTimeRoutes(basic usecase.Basic, realTime usecase.RealTime, history usecase.History) *Router {
	v1.NewRealTimeRoutes(r.v1Group, basic, realTime, history)
	return r
}

func swaggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		docs.SwaggerInfo.Host = c.Request.Host
		c.Next()
	}
}
