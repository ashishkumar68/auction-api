package user

import (
	"context"
	"github.com/ashishkumar68/auction-api/commands"
	"github.com/ashishkumar68/auction-api/repositories"
	"github.com/ashishkumar68/auction-api/services"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var (
	InternalServerErrMsg	= "Something went wrong, sorry please try again later."
	AccountWithEmailExists	= "Sorry! a user with this email already exists"
	InvalidCredentials		= "Invalid credentials were found."
)

func RegisterUser(c *gin.Context) {
	var registerUserCmd commands.RegisterNewUserCommand
	if err := c.ShouldBindJSON(&registerUserCmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !repositories.NewUserRepository().FindByEmail(registerUserCmd.Email).IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": AccountWithEmailExists,
			"email": registerUserCmd.Email,
		})
		return
	}

	bus := commands.NewCommandBus()
	user, err := bus.ExecuteContext(context.Background(), &registerUserCmd)
	if err != nil {
		log.Println("Could not save user")
		log.Println("err:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": InternalServerErrMsg})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func LoginUser(c *gin.Context) {
	var loginUserCmd commands.LoginUserCommand
	if err := c.ShouldBindJSON(&loginUserCmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	bus := commands.NewCommandBus()
	loggedInUser, err := bus.ExecuteContext(context.Background(), &loginUserCmd)
	if err != nil {
		log.Println("Could not login user")
		log.Println("err:", err)
		if err == services.PasswordsDontMatch {
			c.JSON(http.StatusUnauthorized, gin.H{"error": InvalidCredentials})
		} else if err == services.UserEmailDoesntExist {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": InternalServerErrMsg})
		}
		return
	}

	c.JSON(http.StatusOK, loggedInUser)
}