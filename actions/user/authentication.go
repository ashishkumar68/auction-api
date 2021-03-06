package user

import (
	"fmt"
	"github.com/ashishkumar68/auction-api/actions"
	"github.com/ashishkumar68/auction-api/forms"
	"github.com/ashishkumar68/auction-api/repositories"
	"github.com/ashishkumar68/auction-api/services"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func RegisterUser(c *gin.Context) {
	var registerUserForm forms.RegisterNewUserForm
	// Validate
	if err := c.ShouldBindJSON(&registerUserForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Tap into DB connection.
	dbConn := actions.GetDBConnectionByContext(c)
	if nil != repositories.NewRepository(dbConn).FindUserByEmail(registerUserForm.Email) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": actions.AccountWithEmailExists,
			"email": registerUserForm.Email,
		})
		return
	}
	// Add new user.
	user, err := services.NewUserService(dbConn).NewRegister(c, registerUserForm)
	if err != nil {
		log.Println(fmt.Sprintf("Could not save user: %s", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": actions.InternalServerErrMsg})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func LoginUser(c *gin.Context) {
	var loginUserForm forms.LoginUserForm
	if err := c.ShouldBindJSON(&loginUserForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	dbConn := actions.GetDBConnectionByContext(c)
	loggedInUser, err := services.NewUserService(dbConn).Login(c, loginUserForm)
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
