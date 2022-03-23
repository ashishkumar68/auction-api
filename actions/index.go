package actions

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func IndexAction(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}
