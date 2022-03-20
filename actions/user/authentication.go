package user

import (
	"context"
	"github.com/ashishkumar68/auction-api/commands"
	"github.com/ashishkumar68/auction-api/repositories"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var (
	InternalServerErrMsg	= "Something went wrong, sorry please try again later."
	AccountWithEmailExists	= "Sorry! a user with this email already exists"
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
	//user, err := bus.Execute(&registerUserCmd)
	if err != nil {
		log.Println("Could not save user")
		log.Println("err:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": InternalServerErrMsg})
		return
	}

	c.JSON(http.StatusCreated, user)
}
