package routes

import (
	"github.com/ashishkumar68/auction-api/actions/exp"
	"github.com/ashishkumar68/auction-api/actions/user"
	"github.com/ashishkumar68/auction-api/middleware"
	"github.com/gin-gonic/gin"
)

func MapUserRoutes(authGroup *gin.RouterGroup) {
	userTxRouteGroup := authGroup.Group("/users", middleware.TransactionRoute(), middleware.CorsRoute())

	userTxRouteGroup.POST("/register", user.RegisterUser)
	userTxRouteGroup.OPTIONS("/register")

	userTxRouteGroup.POST("/login", user.LoginUser)
	userTxRouteGroup.OPTIONS("/login")
	userTxRouteGroup.GET("/items", middleware.AuthenticatedRoute(), user.ListUserItems)
	userTxRouteGroup.OPTIONS("/items")

	userTxRouteGroup.POST("/sample-uploads", exp.SampleUploadFiles)
}
