//go:build prod

package router

import "github.com/gin-gonic/gin"

func init() {
	gin.SetMode(gin.ReleaseMode)
}
