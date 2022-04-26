package routes

import (
	"github.com/ashishkumar68/auction-api/actions"
	"github.com/gin-gonic/gin"
)

func MapIndexRoutes(engine *gin.Engine) {
	engine.GET("/", actions.IndexAction)
}

var (
	AnonymousRoutes = gin.H{
		"/api/user/login":    []string{"POST"},
		"/api/user/register": []string{"POST"},
		"/api/items":         []string{"GET"},
	}
)
