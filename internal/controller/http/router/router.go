// Package router implements routing paths. Each services in own file.
package router

import (
	"fmt"
	"net/http"
	"net/mail"
	"os"

	"tmt/docs"
	v1 "tmt/internal/controller/http/v1"
	"tmt/internal/entity"
	"tmt/internal/usecase"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const (
	prefix = "/tmt"
)

// Router -.
type Router struct {
	v1Public    *gin.RouterGroup
	v1Private   *gin.RouterGroup
	rootHandler *gin.Engine
}

// NewRouter -.
//
//	@title			TOC MACHINE TRADING
//	@description	API docs for Auto Trade
//	@version		1.0.0
//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html
func NewRouter(system usecase.System) *Router {
	g := gin.New()
	g.Use(gin.Recovery())
	g.GET("/metrics", gin.WrapH(promhttp.Handler()))
	g.GET("/-/health", healthCheck)

	if os.Getenv("DISABLE_SWAGGER_HTTP_HANDLER") == "" {
		docs.SwaggerInfo.BasePath = prefix
		g.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		g.Use(swaggerMiddleware())
	}

	jwtHandler, err := newAuthMiddleware(system)
	if err != nil {
		panic(err)
	}

	v1Public := g.Group(fmt.Sprintf("%s/v1", prefix))
	v1Public.POST("/login", loginHandler(jwtHandler))
	v1Public.POST("/user", newUserHandler(system))
	v1Public.GET("/user/verify/:user/:code", verifyEmailHandler(system))
	v1Public.POST("/user/activate", activateUserHandler(system))

	v1Private := g.Group(fmt.Sprintf("%s/v1", prefix))
	v1Private.Use(jwtHandler.MiddlewareFunc())

	return &Router{
		v1Public:    v1Public,
		v1Private:   v1Private,
		rootHandler: g,
	}
}

// loginHandler loginHandler
//
//	@Description	Every api request will extend token expired time, websocket will not extend.
//	@tags			user
//	@accept			json
//	@produce		json
//	@param			body	body		loginBody{}	true	"Body"
//	@success		200		{object}	loginResponseBody{}
//	@router			/v1/login [post]
func loginHandler(jwtHandler *jwt.GinJWTMiddleware) gin.HandlerFunc {
	return jwtHandler.LoginHandler
}

// newUserHandler newUserHandler
//
//	@Description	add new user
//	@tags			user
//	@accept			json
//	@produce		json
//	@param			body	body	entity.User{}	true	"Body"
//	@success		200
//	@failure		400	{string}	string	"Bad Request"
//	@failure		500	{string}	string	"Internal Server Error"
//	@router			/v1/user [post]
func newUserHandler(system usecase.System) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := entity.User{}
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		if user.Username == "" || user.Password == "" || user.Email == "" {
			c.JSON(http.StatusBadRequest, "username, password, email is required")
			return
		}

		_, err := mail.ParseAddress(user.Email)
		if err != nil {
			c.JSON(http.StatusBadRequest, "email format error")
			return
		}

		if err := system.AddUser(c, &user); err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, nil)
	}
}

// verifyEmailHandler verifyEmailHandler
//
//	@Description	email verification
//	@tags			user
//	@accept			json
//	@produce		json
//	@param			user	path	string	true	"user"
//	@param			code	path	string	true	"code"
//	@success		200
//	@failure		400	{string}	string	"Bad Request"
//	@failure		500	{string}	string	"Internal Server Error"
//	@router			/v1/user/verify/{user}/{code} [get]
func verifyEmailHandler(system usecase.System) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.Param("user")
		if user == "" || user == "undefined" || user == "{user}" {
			c.JSON(http.StatusBadRequest, "user is required")
			return
		}
		code := c.Param("code")
		if code == "" || code == "undefined" || code == "{code}" {
			c.JSON(http.StatusBadRequest, "code is required")
			return
		}
		if err := system.VerifyEmail(c, user, code); err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, "success")
	}
}

// activateUserHandler activateUserHandler
//
//	@Description	active user
//	@tags			user
//	@accept			json
//	@produce		json
//	@param			user	header	string	true	"user"
//	@success		200
//	@failure		400	{string}	string	"Bad Request"
//	@failure		500	{string}	string	"Internal Server Error"
//	@router			/v1/user/activate [post]
func activateUserHandler(system usecase.System) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.GetHeader("user")
		if user == "" {
			c.JSON(http.StatusBadRequest, "user is required")
			return
		}
		if err := system.ActivateUser(c, user); err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, nil)
	}
}

func (r *Router) GetHandler() *gin.Engine {
	return r.rootHandler
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
