package middleware

import (
	"github.com/ashishkumar68/auction-api/actions"
	"github.com/ashishkumar68/auction-api/repositories"
	"github.com/ashishkumar68/auction-api/routes"
	"github.com/ashishkumar68/auction-api/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strings"
)

func AuthenticatedRoute() gin.HandlerFunc {
	return func (c *gin.Context) {
		var dbConnection *gorm.DB
		if db, ok := c.Get("db"); ok {
			dbConnection = db.(*gorm.DB)
		}
		if _, ok := routes.AnonymousRoutes[c.Request.RequestURI]; ok {
			c.Next()
			return
		}
		tokenString := strings.Trim(c.Request.Header.Get("Authorization"), "")
		token, err := services.VerifyJwtToken(tokenString)
		if err != nil {
			log.Println("Could not verify auth token")
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": actions.InvalidCredentials})
			return
		}

		email := token.Header["username"].(string)
		loggedInUser := repositories.NewUserRepository(dbConnection).FindByEmail(email)
		if loggedInUser.IsZero() || !loggedInUser.IsActive {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": actions.InvalidCredentials})
			return
		}
		c.Set("actionUser", loggedInUser)
		dbConnection.Set("actionUser", loggedInUser)

		c.Next()
	}
}