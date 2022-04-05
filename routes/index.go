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
		"/api/login":		[]string{"POST"},
		"/api/register":	[]string{"POST"},
		"/api/items":		[]string{"GET"},
	}
)