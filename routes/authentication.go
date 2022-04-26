package routes

import (
	"github.com/ashishkumar68/auction-api/actions/user"
	"github.com/gin-gonic/gin"
)

func MapUserRoutes(authGroup *gin.RouterGroup) {
	authGroup.POST("/user/register", user.RegisterUser)
	authGroup.POST("/user/login", user.LoginUser)

	authGroup.GET("/user/items", user.ListUserItems)
}
