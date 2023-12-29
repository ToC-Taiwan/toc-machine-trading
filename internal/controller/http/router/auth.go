package router

import (
	"net/http"
	"time"

	"tmt/internal/usecase"

	jwt "github.com/appleboy/gin-jwt/v2"
	v4jwt "github.com/golang-jwt/jwt/v4"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	identityKey = "tmt_identity"
	timeOut     = time.Hour
)

func newAuthMiddleware(system usecase.System) (*jwt.GinJWTMiddleware, error) {
	m := jwt.GinJWTMiddleware{
		TokenLookup:      "header:Authorization",
		SigningAlgorithm: "HS256",
		Timeout:          timeOut,
		TimeFunc:         time.Now,
		TokenHeadName:    "Bearer",
		Authorizator: func(interface{}, *gin.Context) bool {
			return true
		},
		Unauthorized:          unauthorized,
		LoginResponse:         loginResponse,
		LogoutResponse:        logoutResponse,
		RefreshResponse:       refreshResponse,
		IdentityKey:           identityKey,
		IdentityHandler:       identityHandler,
		HTTPStatusMessageFunc: hTTPStatusMessageFunc,
		Realm:                 "tmt_jwt",
		CookieMaxAge:          timeOut,
		CookieName:            "tmt",

		Key:           []byte(uuid.New().String()),
		MaxRefresh:    time.Hour,
		Authenticator: authenticator(system),
		PayloadFunc:   payloadFunc,

		// PrivKeyFile:          "",
		// PrivKeyBytes:         []byte{},
		// PubKeyFile:           "",
		// PrivateKeyPassphrase: "",
		// PubKeyBytes:          []byte{},
		// CookieDomain:      "",
		// SendCookie:        false,
		// SecureCookie:      false,
		// CookieHTTPOnly:    false,
		// SendAuthorization: false,
		// DisabledAbort:     false,
		// CookieSameSite:    1,

		ParseOptions: []v4jwt.ParserOption{},
	}
	return jwt.New(&m)
}

func unauthorized(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}

type loginResponseBody struct {
	Token  string `json:"token"`
	Expire string `json:"expire"`
	Code   int    `json:"code"`
}

func loginResponse(c *gin.Context, code int, token string, expire time.Time) {
	c.JSON(http.StatusOK, loginResponseBody{
		Token:  token,
		Expire: expire.Format(time.RFC3339),
		Code:   http.StatusOK,
	})
}

func logoutResponse(c *gin.Context, code int) {
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
	})
}

func refreshResponse(c *gin.Context, code int, token string, expire time.Time) {
	c.JSON(http.StatusOK, gin.H{
		"code":   http.StatusOK,
		"token":  token,
		"expire": expire.Format(time.RFC3339),
	})
}

func identityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	return claims[identityKey]
}

func hTTPStatusMessageFunc(e error, c *gin.Context) string {
	return e.Error()
}

type loginBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func authenticator(system usecase.System) func(c *gin.Context) (interface{}, error) {
	return func(c *gin.Context) (interface{}, error) {
		var loginVals loginBody
		if err := c.ShouldBind(&loginVals); err != nil {
			return "", jwt.ErrMissingLoginValues
		}
		user := loginVals.Username
		password := loginVals.Password
		ok, err := system.Login(c.Request.Context(), user, password)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, jwt.ErrFailedAuthentication
		}
		return user, nil
	}
}

func payloadFunc(data interface{}) jwt.MapClaims {
	return nil
}
