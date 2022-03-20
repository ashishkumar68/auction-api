package routes

import (
	"github.com/ashishkumar68/auction-api/actions/user"
	"github.com/gin-gonic/gin"
)

func MapAuthRoutes(engine *gin.Engine) {
	engine.POST("/register", user.RegisterUser)
}