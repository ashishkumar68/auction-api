package routes

import (
	"github.com/ashishkumar68/auction-api/actions/item"
	"github.com/gin-gonic/gin"
)

func MapItemRoutes(itemsGroup *gin.RouterGroup) {
	itemsGroup.POST("/items", item.CreateItem)
	itemsGroup.GET("/items", item.ListItems)
	itemsGroup.POST("/items/:itemId/bid", item.PlaceBidOnItem)
}
