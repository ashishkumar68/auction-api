package routes

import (
	"github.com/ashishkumar68/auction-api/actions"
	"github.com/ashishkumar68/auction-api/middleware"
	"github.com/gin-gonic/gin"
)

func MapIndexRoutes(engine *gin.Engine) {
	engine.GET("/", middleware.CorsRoute(), actions.IndexAction)
	engine.OPTIONS("/", middleware.CorsRoute())
}
