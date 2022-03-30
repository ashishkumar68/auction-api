package item

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateItem(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"item": 123,
	})
}
