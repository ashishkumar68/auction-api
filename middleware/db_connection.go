package middleware

import (
	"github.com/ashishkumar68/auction-api/database"
	"github.com/gin-gonic/gin"
)

func AssignNewDatabaseConnection() gin.HandlerFunc {
	return func(c *gin.Context) {
		dbConnection := database.NewConnectionWithContext(c)

		c.Set("db", dbConnection)

		c.Next()
	}
}
