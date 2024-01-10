// Package auth package auth
package auth

import (
	"fmt"
	"net/http"
	"time"

	"tmt/internal/usecase"

	jwt "github.com/appleboy/gin-jwt/v2"
	v4jwt "github.com/golang-jwt/jwt/v4"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	tokenHeaderName = "Bearer"
	identityKey     = "tmt_identity"
	timeOut         = time.Hour
)

func NewAuthMiddleware(system usecase.System) (*jwt.GinJWTMiddleware, error) {
	m := jwt.GinJWTMiddleware{
		TokenLookup:      "header:Authorization",
		SigningAlgorithm: "HS256",
		Timeout:          timeOut,
		TimeFunc:         time.Now,
		TokenHeadName:    tokenHeaderName,
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
		MaxRefresh:    timeOut,
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

func loginResponse(c *gin.Context, code int, token string, expire time.Time) {
	c.JSON(http.StatusOK, LoginResponseBody{
		Token:  fmt.Sprintf("%s %s", tokenHeaderName, token),
		Expire: expire.Format(time.RFC3339),
		Code:   http.StatusOK,
	})
}

func logoutResponse(c *gin.Context, code int) {
	c.JSON(http.StatusOK, LogoutResponseBody{
		Code: code,
	})
}

func refreshResponse(c *gin.Context, code int, token string, expire time.Time) {
	c.JSON(http.StatusOK, RefreshResponseBody{
		Token:  fmt.Sprintf("%s %s", tokenHeaderName, token),
		Expire: expire.Format(time.RFC3339),
		Code:   code,
	})
}

func identityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	return claims[identityKey]
}

func hTTPStatusMessageFunc(e error, c *gin.Context) string {
	return e.Error()
}

func authenticator(system usecase.System) func(c *gin.Context) (interface{}, error) {
	return func(c *gin.Context) (interface{}, error) {
		var loginVals LoginBody
		if err := c.ShouldBind(&loginVals); err != nil {
			return "", jwt.ErrMissingLoginValues
		}
		err := system.Login(c.Request.Context(), loginVals.Username, loginVals.Password)
		if err != nil {
			return nil, err
		}
		return loginVals, nil
	}
}

func payloadFunc(data interface{}) jwt.MapClaims {
	if v, ok := data.(LoginBody); ok {
		return jwt.MapClaims{
			"username": v.Username,
		}
	}
	return nil
}

func ExtractUsername(c *gin.Context) string {
	claims := jwt.ExtractClaims(c)
	if v, ok := claims["username"]; ok {
		return v.(string)
	}
	return ""
}
