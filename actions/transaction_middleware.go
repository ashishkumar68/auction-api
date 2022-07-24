package actions

import (
	"github.com/ashishkumar68/auction-api/database"
	"github.com/gin-gonic/gin"
	"log"
)

func TransactionRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		dbConn := database.GetDBHandle().WithContext(c).Begin()
		c.Set("db", dbConn)

		c.Next()

		if isSuccessStatus(c.Writer.Status()) {
			log.Println("committing transactions.")
			err := dbConn.Commit().Error
			if err != nil {
				log.Println("error while committing transactions", err)
			}
		} else {
			err := dbConn.Rollback().Error
			if err != nil {
				log.Println("error while doing rollback transactions", err)
			}
		}
	}
}
