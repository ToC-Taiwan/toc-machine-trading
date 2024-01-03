//go:build !prod

package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	gin.SetMode(gin.DebugMode)
	swagHandler = ginSwagger.WrapHandler(swaggerFiles.Handler)
}
