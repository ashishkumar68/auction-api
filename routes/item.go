package routes

import (
	"github.com/ashishkumar68/auction-api/actions/item"
	"github.com/ashishkumar68/auction-api/middleware"
	"github.com/gin-gonic/gin"
)

func MapItemRoutes(itemsGroup *gin.RouterGroup) {
	itemsTxRouteGroup := itemsGroup.Group("/items", middleware.TransactionRoute())

	itemsTxRouteGroup.POST("", middleware.AuthenticatedRoute(), item.CreateItem)
	itemsTxRouteGroup.GET("", item.ListItems)
	itemsTxRouteGroup.PATCH("/:itemId", middleware.AuthenticatedRoute(), item.EditItem)
	itemsTxRouteGroup.POST("/:itemId/images", middleware.AuthenticatedRoute(), item.AddItemImages)
	itemsTxRouteGroup.DELETE("/:itemId/images/:imageId", middleware.AuthenticatedRoute(), item.DeleteItemImage)
	itemsTxRouteGroup.DELETE("/:itemId/images", middleware.AuthenticatedRoute(), item.DeleteItemImages)

	itemsTxRouteGroup.PUT("/:itemId/mark-off-bid", middleware.AuthenticatedRoute(), item.MarkItemOffBid)
	itemsTxRouteGroup.POST("/:itemId/bid", middleware.AuthenticatedRoute(), item.PlaceBidOnItem)

	itemsTxRouteGroup.POST("/:itemId/reaction", middleware.AuthenticatedRoute(), item.AddReactionToItem)
	itemsTxRouteGroup.DELETE("/:itemId/reaction", middleware.AuthenticatedRoute(), item.RemoveItemReaction)

	itemsTxRouteGroup.POST("/:itemId/comment", middleware.AuthenticatedRoute(), item.AddItemComment)
	itemsTxRouteGroup.PATCH("/:itemId/comment/:commentId", middleware.AuthenticatedRoute(), item.UpdateItemComment)
	itemsTxRouteGroup.DELETE("/:itemId/comment/:commentId", middleware.AuthenticatedRoute(), item.DeleteItemComment)
}
