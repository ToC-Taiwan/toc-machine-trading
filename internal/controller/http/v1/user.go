package v1

import (
	"net/http"
	"net/mail"

	"tmt/internal/controller/http/auth"
	"tmt/internal/controller/http/resp"
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
	h.PUT("/user/auth", r.updateAuthTradeUser)
	h.PUT("/user/push-token", r.updateUserPushToken)
	h.GET("/user/verify/:user/:code", r.verifyEmailHandler)
}

// newUserHandler _.
//
//	@tags		User V1
//	@Summary	New user
//	@accept		json
//	@produce	json
//	@param		body	body	entity.User{}	true	"Body"
//	@success	200
//	@failure	400	{string}	string
//	@failure	500	{string}	string
//	@router		/v1/user [post]
func (u *userRoutes) newUserHandler(c *gin.Context) {
	user := entity.User{}
	if err := c.ShouldBindJSON(&user); err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if user.Username == "" || user.Password == "" || user.Email == "" {
		resp.ErrorResponse(c, http.StatusBadRequest, "username, password, email is required")
		return
	}

	_, err := mail.ParseAddress(user.Email)
	if err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, "email format error")
		return
	}

	if err := u.system.AddUser(c, &user); err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

// verifyEmailHandler _.
//
//	@tags		User V1
//	@Summary	Verify email
//	@accept		json
//	@produce	json
//	@param		user	path	string	true	"user"
//	@param		code	path	string	true	"code"
//	@success	200
//	@failure	400	{string}	string
//	@failure	500	{string}	string
//	@router		/v1/user/verify/{user}/{code} [get]
func (u *userRoutes) verifyEmailHandler(c *gin.Context) {
	result := "Success"
	user := c.Param("user")
	if user == "" || user == "undefined" || user == "{user}" {
		result = "User is required"
	}
	code := c.Param("code")
	if code == "" || code == "undefined" || code == "{code}" {
		result = "Code is required"
	}
	if err := u.system.VerifyEmail(c, user, code); err != nil {
		result = err.Error()
	}
	c.HTML(http.StatusOK, "mail_verification.tmpl", gin.H{"result": result})
}

// loginHandler _.
//
//	@tags		User V1
//	@Summary	Login
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
//	@tags		User V1
//	@Summary	Logout
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
//	@tags		User V1
//	@Summary	Refresh token
//	@security	JWT
//	@accept		json
//	@produce	json
//	@success	200	{object}	auth.RefreshResponseBody{}
//	@failure	401	{object}	auth.UnauthorizedResponseBody{}
//	@router		/v1/refresh [get]
func (u *userRoutes) refreshTokenHandler(c *gin.Context) {
	u.jwtHandler.RefreshHandler(c)
}

type userPushTokenRequest struct {
	PushToken string `json:"push_token"`
}

// updateUserPushToken _.
//
//	@tags		User V1
//	@Summary	Update user push token
//	@security	JWT
//	@accept		json
//	@produce	json
//	@param		body	body	userPushTokenRequest{}	true	"Body"
//	@success	200
//	@failure	400	{object}	resp.Response{}
//	@failure	401	{object}	auth.UnauthorizedResponseBody{}
//	@failure	500	{object}	resp.Response{}
//	@router		/v1/user/push-token [put]
func (u *userRoutes) updateUserPushToken(c *gin.Context) {
	p := userPushTokenRequest{}
	if err := c.ShouldBindJSON(&p); err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	username := auth.ExtractUsername(c)
	if username == "" {
		resp.ErrorResponse(c, http.StatusBadRequest, "username is required in token")
		return
	}

	if err := u.system.UpdateUserPushToken(c.Request.Context(), username, p.PushToken); err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

// updateAuthTradeUser _.
//
//	@tags		User V1
//	@Summary	Update auth trade user
//	@security	JWT
//	@accept		json
//	@produce	json
//	@success	200
//	@failure	401	{object}	auth.UnauthorizedResponseBody{}
//	@router		/v1/user/auth [put]
func (u *userRoutes) updateAuthTradeUser(c *gin.Context) {
	u.system.UpdateAuthTradeUser()
	c.JSON(http.StatusOK, nil)
}
