package user

import (
	"github.com/ashishkumar68/auction-api/actions"
	"github.com/gin-gonic/gin"
)

func MapUserRoutes(authGroup *gin.RouterGroup) {
	userTxRouteGroup := authGroup.Group("/users", actions.TransactionRoute(), actions.CorsRoute())

	userTxRouteGroup.POST("/register", RegisterUser)
	userTxRouteGroup.OPTIONS("/register")

	userTxRouteGroup.POST("/login", LoginUser)
	userTxRouteGroup.OPTIONS("/login")
	userTxRouteGroup.GET("/items", actions.AuthenticatedRoute(), ListUserItems)
	userTxRouteGroup.OPTIONS("/items")
}
