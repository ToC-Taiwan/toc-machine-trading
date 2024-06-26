package v1

import (
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/toc-taiwan/toc-machine-trading/internal/controller/http/auth"
	"github.com/toc-taiwan/toc-machine-trading/internal/controller/http/resp"
	"github.com/toc-taiwan/toc-machine-trading/internal/entity"
	"github.com/toc-taiwan/toc-machine-trading/internal/usecase"
)

type userRoutes struct {
	system     usecase.System
	jwtHandler *jwt.GinJWTMiddleware
}

func NewUserRoutes(public *gin.RouterGroup, private *gin.RouterGroup, jwtHandler *jwt.GinJWTMiddleware, system usecase.System) {
	r := &userRoutes{
		system:     system,
		jwtHandler: jwtHandler,
	}

	public.POST("/login", r.loginHandler)
	public.GET("/refresh", r.refreshTokenHandler)

	public.POST("/user", r.newUserHandler)
	public.POST("/user/verify/:user/:code", r.verifyEmailHandler)

	private.GET("/user/info", r.getUserInfo)
	private.PUT("/user/auth", r.updateAuthTradeUser)
	private.GET("/user/push-token", r.getUserPushTokenStatus)
	private.PUT("/user/push-token", r.updateUserPushToken)
	private.DELETE("/user/push-token", r.clearAllPushToken)

	private.GET("/logout", r.logutHandler)
}

// newUserHandler _.
//
//	@tags		User V1
//	@Summary	New user
//	@accept		json
//	@produce	json
//	@param		body	body	entity.NewUser{}	true	"Body"
//	@success	200
//	@failure	400	{object}	resp.Response{}
//	@failure	500	{object}	resp.Response{}
//	@router		/v1/user [post]
func (u *userRoutes) newUserHandler(c *gin.Context) {
	user := entity.NewUser{}
	if err := c.ShouldBindJSON(&user); err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	if err := u.system.AddUser(c, &user); err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err)
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
//	@router		/v1/user/verify/{user}/{code} [post]
func (u *userRoutes) verifyEmailHandler(c *gin.Context) {
	user := c.Param("user")
	if user == "" || user == "undefined" || user == "{user}" {
		resp.ErrorResponse(c, http.StatusBadRequest, "User is required")
		return
	}
	code := c.Param("code")
	if code == "" || code == "undefined" || code == "{code}" {
		resp.ErrorResponse(c, http.StatusBadRequest, "Code is required")
		return
	}
	if err := u.system.VerifyEmail(c, user, code); err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, nil)
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
//	@success	200
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
//	@success	200	{object}	auth.LoginResponseBody{}
//	@failure	401	{object}	resp.Response{}
//	@router		/v1/refresh [get]
func (u *userRoutes) refreshTokenHandler(c *gin.Context) {
	u.jwtHandler.RefreshHandler(c)
}

type userPushTokenRequest struct {
	PushToken string `json:"push_token"`
	Enabled   bool   `json:"enabled"`
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
//	@failure	401	{object}	resp.Response{}
//	@failure	500	{object}	resp.Response{}
//	@router		/v1/user/push-token [put]
func (u *userRoutes) updateUserPushToken(c *gin.Context) {
	p := userPushTokenRequest{}
	if err := c.ShouldBindJSON(&p); err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	if p.PushToken == "" {
		resp.ErrorResponse(c, http.StatusBadRequest, "push_token is required")
		return
	}

	username := auth.ExtractUsername(c)
	if username == "" {
		resp.ErrorResponse(c, http.StatusBadRequest, "username is required in token")
		return
	}

	if err := u.system.InsertPushToken(c.Request.Context(), p.PushToken, username, p.Enabled); err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}

type pushTokenStatusResponse struct {
	Enabled bool `json:"enabled"`
}

// getUserPushTokenStatus _.
//
//	@tags		User V1
//	@Summary	Get user push token status
//	@security	JWT
//	@accept		json
//	@produce	json
//	@param		token	header		string	true	"token"
//	@success	200		{object}	pushTokenStatusResponse{}
//	@failure	400		{object}	resp.Response{}
//	@failure	401		{object}	resp.Response{}
//	@failure	500		{object}	resp.Response{}
//	@router		/v1/user/push-token [get]
func (u *userRoutes) getUserPushTokenStatus(c *gin.Context) {
	token := c.GetHeader("token")
	if token == "" {
		resp.ErrorResponse(c, http.StatusBadRequest, "token is required")
		return
	}
	enabled, err := u.system.IsPushTokenEnabled(c.Request.Context(), token)
	if err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, pushTokenStatusResponse{Enabled: enabled})
}

// clearAllPushToken _.
//
//	@tags		User V1
//	@Summary	Clear all push token
//	@security	JWT
//	@accept		json
//	@produce	json
//	@success	200
//	@failure	500	{object}	resp.Response{}
//	@router		/v1/user/push-token [delete]
func (u *userRoutes) clearAllPushToken(c *gin.Context) {
	if err := u.system.DeleteAllPushTokens(c); err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err)
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
//	@failure	401	{object}	resp.Response{}
//	@router		/v1/user/auth [put]
func (u *userRoutes) updateAuthTradeUser(c *gin.Context) {
	u.system.UpdateAuthTradeUser()
	c.JSON(http.StatusOK, nil)
}

// getUserInfo _.
//
//	@tags		User V1
//	@Summary	Get user info
//	@security	JWT
//	@accept		json
//	@produce	json
//	@success	200
//	@failure	401	{object}	resp.Response{}
//	@router		/v1/user/info [get]
func (u *userRoutes) getUserInfo(c *gin.Context) {
	user := auth.ExtractUsername(c)
	if user == "" {
		resp.ErrorResponse(c, http.StatusBadRequest, "username is required in token")
		return
	}
	info, err := u.system.GetUserInfo(c.Request.Context(), user)
	if err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, info)
}
