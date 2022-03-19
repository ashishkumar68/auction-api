package routes

import (
	"github.com/ashishkumar68/auction-api/actions"
	"github.com/gin-gonic/gin"
)

func MapIndexRoutes(engine *gin.Engine) {
	engine.GET("/", actions.IndexAction)
}