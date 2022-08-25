package actions

import (
	"github.com/ashishkumar68/auction-api/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"os"
)

var (
	InternalServerErrMsg          = "Something went wrong, sorry please try again later."
	AccountWithEmailExists        = "Sorry! a user with this email already exists"
	InvalidCredentials            = "Invalid credentials were found."
	InvalidItemIdReceivedErr      = "Invalid item id was received in request."
	InvalidCommentIdReceivedErr   = "Invalid comment id was received in request."
	CommentNotAuthoredByUser      = "Comment not authored by user."
	CouldNotReadFromMultipartForm = "Could not read from request form."
	ImagesNotFoundInRequest       = "Could not find images in request."
	CouldNotSaveImage             = "Could not save image"
	InvalidRemoveExistingVal      = "Invalid removeExisting value was found."
	InvalidImageIdFoundErr        = "Invalid image id value was found."
	ItemImageNotFoundError        = "Item image was not found."
)

func IndexAction(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

func FetchBuildVersionAction(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"version": os.Getenv("VERSION")})
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
