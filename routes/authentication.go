package routes

import (
	"github.com/ashishkumar68/auction-api/actions/user"
	"github.com/ashishkumar68/auction-api/middleware"
	"github.com/gin-gonic/gin"
)

func MapUserRoutes(authGroup *gin.RouterGroup) {
	userTxRouteGroup := authGroup.Group("/user", middleware.TransactionRoute())

	userTxRouteGroup.POST("/register", user.RegisterUser)
	userTxRouteGroup.POST("/login", user.LoginUser)
	userTxRouteGroup.GET("/items", middleware.AuthenticatedRoute(), user.ListUserItems)
}
