package middleware

import (
	"fmt"
	"github.com/ashishkumar68/auction-api/actions"
	"github.com/ashishkumar68/auction-api/repositories"
	"github.com/ashishkumar68/auction-api/services"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

func AuthenticatedRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := strings.Trim(c.Request.Header.Get("Authorization"), "")
		token, err := services.VerifyJwtToken(tokenString)
		if err != nil {
			log.Println(fmt.Sprintf("Could not verify auth token, err: %s", err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": actions.InvalidCredentials})
			return
		}

		email := token.Header["username"].(string)
		dbConnection := actions.GetDBConnectionByContext(c)
		loggedInUser := repositories.NewRepository(dbConnection).FindUserByEmail(email)
		if nil == loggedInUser || loggedInUser.IsZero() || !loggedInUser.IsActive {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": actions.InvalidCredentials})
			return
		}
		// set logged-in user to context to be used in further handlers.
		c.Set("actionUser", loggedInUser)
		c.Next()
	}
}

func isSuccessStatus(statusCode int) bool {
	for _, code := range []int{http.StatusOK, http.StatusCreated, http.StatusAccepted, http.StatusNoContent} {
		if code == statusCode {
			return true
		}
	}

	return false
}
