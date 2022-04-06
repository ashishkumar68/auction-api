package middleware

import (
	"fmt"
	"github.com/ashishkumar68/auction-api/actions"
	"github.com/ashishkumar68/auction-api/database"
	"github.com/ashishkumar68/auction-api/repositories"
	"github.com/ashishkumar68/auction-api/routes"
	"github.com/ashishkumar68/auction-api/services"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

func AuthenticatedRoute() gin.HandlerFunc {
	isAnonymousRoute := func(c *gin.Context) bool {
		isAnonymous := false
		if allowedMethods, ok := routes.AnonymousRoutes[c.Request.URL.Path]; ok {
			for _, method := range allowedMethods.([]string) {
				if method == c.Request.Method {
					isAnonymous = true
					break
				}
			}
		}

		return isAnonymous
	}

	return func(c *gin.Context) {
		dbConnection := database.NewConnectionWithContext(c)
		if isAnonymousRoute(c) {
			c.Set("db", dbConnection)
			c.Next()
			return
		}
		tokenString := strings.Trim(c.Request.Header.Get("Authorization"), "")
		token, err := services.VerifyJwtToken(tokenString)
		if err != nil {
			log.Println(fmt.Sprintf("Could not verify auth token, err: %s", err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": actions.InvalidCredentials})
			return
		}

		email := token.Header["username"].(string)
		loggedInUser := repositories.NewUserRepository(dbConnection).FindByEmail(email)
		if loggedInUser.IsZero() || !loggedInUser.IsActive {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": actions.InvalidCredentials})
			return
		}
		dbConnection = dbConnection.Set("actionUser", *loggedInUser)
		// set logged-in user and db connection to context to be used in further handlers.
		c.Set("actionUser", loggedInUser)
		c.Set("db", dbConnection)

		c.Next()
	}
}
