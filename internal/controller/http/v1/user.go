package v1

import (
	"net/http"
	"net/mail"

	"tmt/internal/entity"
	"tmt/internal/usecase"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type userRoutes struct {
	system     usecase.System
	jwtHandler *jwt.GinJWTMiddleware
}

func NewUserRoutes(h *gin.RouterGroup, jwtHandler *jwt.GinJWTMiddleware, system usecase.System) {
	r := &userRoutes{
		system:     system,
		jwtHandler: jwtHandler,
	}
	h.POST("/login", r.loginHandler)
	h.GET("/logout", r.logutHandler)
	h.GET("/refresh", r.refreshTokenHandler)

	h.POST("/user", r.newUserHandler)
	h.GET("/user/verify/:user/:code", r.verifyEmailHandler)
	h.POST("/user/activate", r.activateUserHandler)
}

// loginHandler _.
//
//	@Summary	loginHandler
//	@tags		user
//	@accept		json
//	@produce	json
//	@param		body	body		auth.LoginBody{}	true	"Body"
//	@success	200		{object}	auth.LoginResponseBody{}
//	@router		/v1/login [post]
func (u *userRoutes) loginHandler(c *gin.Context) {
	u.jwtHandler.LoginHandler(c)
}

// logutHandler _.
//
//	@Summary	logutHandler
//	@tags		user
//	@security	JWT
//	@accept		json
//	@produce	json
//	@success	200	{object}	auth.LogoutResponseBody{}
//	@router		/v1/logout [get]
func (u *userRoutes) logutHandler(c *gin.Context) {
	u.jwtHandler.LogoutHandler(c)
}

// refreshTokenHandler _.
//
//	@Summary	refreshTokenHandler
//	@tags		user
//	@security	JWT
//	@accept		json
//	@produce	json
//	@success	200	{object}	auth.RefreshResponseBody{}
//	@router		/v1/refresh [get]
func (u *userRoutes) refreshTokenHandler(c *gin.Context) {
	u.jwtHandler.RefreshHandler(c)
}

// newUserHandler _.
//
//	@Summary	newUserHandler
//	@tags		user
//	@accept		json
//	@produce	json
//	@param		body	body	entity.User{}	true	"Body"
//	@success	200
//	@failure	400	{string}	string	"Bad Request"
//	@failure	500	{string}	string	"Internal Server Error"
//	@router		/v1/user [post]
func (u *userRoutes) newUserHandler(c *gin.Context) {
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

	if err := u.system.AddUser(c, &user); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

// verifyEmailHandler _.
//
//	@Summary	verifyEmailHandler
//	@tags		user
//	@accept		json
//	@produce	json
//	@param		user	path	string	true	"user"
//	@param		code	path	string	true	"code"
//	@success	200
//	@failure	400	{string}	string	"Bad Request"
//	@failure	500	{string}	string	"Internal Server Error"
//	@router		/v1/user/verify/{user}/{code} [get]
func (u *userRoutes) verifyEmailHandler(c *gin.Context) {
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
	if err := u.system.VerifyEmail(c, user, code); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, "success")
}

// activateUserHandler _.
//
//	@Summary	activateUserHandler
//	@tags		user
//	@accept		json
//	@produce	json
//	@param		user	header	string	true	"user"
//	@success	200
//	@failure	400	{string}	string	"Bad Request"
//	@failure	500	{string}	string	"Internal Server Error"
//	@router		/v1/user/activate [post]
func (u *userRoutes) activateUserHandler(c *gin.Context) {
	user := c.GetHeader("user")
	if user == "" {
		c.JSON(http.StatusBadRequest, "user is required")
		return
	}
	if err := u.system.ActivateUser(c, user); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, nil)
}
