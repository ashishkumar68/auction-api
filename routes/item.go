package routes

import (
	"github.com/ashishkumar68/auction-api/actions/item"
	"github.com/gin-gonic/gin"
)

func MapItemRoutes(itemsGroup *gin.RouterGroup) {
	itemsGroup.POST("/items", item.CreateItem)
	itemsGroup.GET("/items", item.ListItems)
	itemsGroup.PATCH("/items/:itemId", item.EditItem)
	itemsGroup.POST("/items/:itemId/images", item.AddItemImages)

	itemsGroup.PUT("/items/:itemId/mark-off-bid", item.MarkItemOffBid)
	itemsGroup.POST("/items/:itemId/bid", item.PlaceBidOnItem)

	itemsGroup.POST("/items/:itemId/reaction", item.AddReactionToItem)
	itemsGroup.DELETE("/items/:itemId/reaction", item.RemoveItemReaction)

	itemsGroup.POST("/items/:itemId/comment", item.AddItemComment)
	itemsGroup.PATCH("/items/:itemId/comment/:commentId", item.UpdateItemComment)
	itemsGroup.DELETE("/items/:itemId/comment/:commentId", item.DeleteItemComment)
}
