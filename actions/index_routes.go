package actions

import (
	"github.com/gin-gonic/gin"
)

const (
	BaseApiRoute = "/api"
)

func MapIndexRoutes(engine *gin.Engine) {
	engine.GET("/", CorsRoute(), IndexAction)
	engine.OPTIONS("/", CorsRoute())
}
