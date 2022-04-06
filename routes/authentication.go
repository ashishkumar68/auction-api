package routes

import (
	"github.com/ashishkumar68/auction-api/actions/user"
	"github.com/gin-gonic/gin"
)

func MapAuthRoutes(authGroup *gin.RouterGroup) {
	authGroup.POST("/register", user.RegisterUser)
	authGroup.POST("/login", user.LoginUser)
}
