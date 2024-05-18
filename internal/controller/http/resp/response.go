// Package resp package resp
package resp

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/toc-taiwan/toc-machine-trading/internal/usecase"
)

// Response -.
type Response struct {
	Code     int    `json:"code"`
	Response string `json:"response"`
}

// ErrorResponse -.
func ErrorResponse(c *gin.Context, code int, err any) {
	switch v := err.(type) {
	case *usecase.UseCaseError:
		c.AbortWithStatusJSON(code, Response{
			Code:     v.Code,
			Response: v.Message,
		})
	case error:
		c.AbortWithStatusJSON(code, Response{
			Code:     code,
			Response: v.Error(),
		})
	default:
		c.AbortWithStatusJSON(code, Response{
			Code:     code,
			Response: fmt.Sprintf("%+v", v),
		})
	}
}
