package actions

import (
	"github.com/ashishkumar68/auction-api/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

var (
	InternalServerErrMsg     = "Something went wrong, sorry please try again later."
	AccountWithEmailExists   = "Sorry! a user with this email already exists"
	InvalidCredentials       = "Invalid credentials were found."
	InvalidItemIdReceivedErr = "Invalid item id was received in request."
)

func IndexAction(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

func GetDBConnectionByContext(c *gin.Context) *gorm.DB {
	var dbConn *gorm.DB
	if db, ok := c.Get("db"); ok {
		dbConn = db.(*gorm.DB)
	}

	return dbConn
}

func GetActionUserByContext(c *gin.Context) *models.User {
	var actionUser *models.User
	if user, ok := c.Get("actionUser"); ok {
		actionUser = user.(*models.User)
	}

	return actionUser
}
