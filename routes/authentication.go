package routes

import (
	"github.com/ashishkumar68/auction-api/actions/user"
	"github.com/ashishkumar68/auction-api/repositories"
	"github.com/ashishkumar68/auction-api/services"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

func MapAuthRoutes(authGroup *gin.RouterGroup) {
	authGroup.POST("/register", user.RegisterUser)
	authGroup.POST("/login", user.LoginUser)
}

func AuthenticatedRoute() gin.HandlerFunc {
	return func (c *gin.Context) {
		if _, ok := AnonymousRoutes[c.Request.RequestURI]; ok {
			c.Next()
			return
		}
		tokenString := strings.Trim(c.Request.Header.Get("Authorization"), "")
		token, err := services.VerifyJwtToken(tokenString)
		if err != nil {
			log.Println("Could not verify auth token")
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": user.InvalidCredentials})
			return
		}
		email := token.Header["username"].(string)
		loggedInUser := repositories.NewUserRepository().FindByEmail(email)
		if loggedInUser.IsZero() || !loggedInUser.IsActive {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": user.InvalidCredentials})
			return
		}
		c.Set("user", *loggedInUser)

		c.Next()
	}
}