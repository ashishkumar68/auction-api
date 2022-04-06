package user

import (
	"fmt"
	"github.com/ashishkumar68/auction-api/actions"
	"github.com/ashishkumar68/auction-api/commands"
	"github.com/ashishkumar68/auction-api/repositories"
	"github.com/ashishkumar68/auction-api/services"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func RegisterUser(c *gin.Context) {
	var registerUserCmd commands.RegisterNewUserCommand
	// Validate
	if err := c.ShouldBindJSON(&registerUserCmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Tap into DB connection.
	dbConn := actions.GetDBConnectionByContext(c)
	if !repositories.NewUserRepository(dbConn).FindByEmail(registerUserCmd.Email).IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": actions.AccountWithEmailExists,
			"email": registerUserCmd.Email,
		})
		return
	}
	// Add new user.
	bus := commands.NewCommandBus()
	registerUserCmd.DB = dbConn
	user, err := bus.ExecuteContext(c, &registerUserCmd)
	if err != nil {
		log.Println(fmt.Sprintf("Could not save user: %s", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": actions.InternalServerErrMsg})
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
	dbConn := actions.GetDBConnectionByContext(c)
	loginUserCmd.DB = dbConn
	bus := commands.NewCommandBus()
	loggedInUser, err := bus.ExecuteContext(c, &loginUserCmd)
	if err != nil {
		log.Println(fmt.Sprintf("Could not login user: %s", err))
		if err == services.PasswordsDontMatch {
			c.JSON(http.StatusUnauthorized, gin.H{"error": actions.InvalidCredentials})
		} else if err == services.UserEmailDoesntExist {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": actions.InternalServerErrMsg})
		}
		return
	}

	c.JSON(http.StatusOK, loggedInUser)
}
