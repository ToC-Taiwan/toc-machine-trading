// Package resp package resp
package resp

import (
	"github.com/gin-gonic/gin"
)

// Response -.
type Response struct {
	Response string `json:"response"`
}

// ErrorResponse -.
func ErrorResponse(c *gin.Context, code int, msg string) {
	c.AbortWithStatusJSON(code, Response{msg})
}
