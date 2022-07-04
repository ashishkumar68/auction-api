package routes

import (
	"github.com/ashishkumar68/auction-api/actions/item"
	"github.com/ashishkumar68/auction-api/middleware"
	"github.com/gin-gonic/gin"
)

func MapItemRoutes(itemsGroup *gin.RouterGroup) {
	itemsTxRouteGroup := itemsGroup.Group("/items", middleware.TransactionRoute(), middleware.CorsRoute())

	itemsTxRouteGroup.POST("", middleware.AuthenticatedRoute(), item.CreateItem)
	itemsTxRouteGroup.GET("", item.ListItems)
	itemsTxRouteGroup.OPTIONS("")

	itemsTxRouteGroup.PATCH("/:itemId", middleware.AuthenticatedRoute(), item.EditItem)
	itemsTxRouteGroup.OPTIONS("/:itemId")

	itemsTxRouteGroup.POST("/:itemId/images", middleware.AuthenticatedRoute(), item.AddItemImages)
	itemsTxRouteGroup.OPTIONS("/:itemId/images")

	itemsTxRouteGroup.DELETE("/:itemId/images/:imageId", middleware.AuthenticatedRoute(), item.DeleteItemImage)
	itemsTxRouteGroup.OPTIONS("/:itemId/images/:imageId")

	itemsTxRouteGroup.PATCH("/:itemId/images/:imageId/make-thumbnail", middleware.AuthenticatedRoute(), item.MakeItemImageThumbnail)
	itemsTxRouteGroup.OPTIONS("/:itemId/images/:imageId/make-thumbnail")

	itemsTxRouteGroup.DELETE("/:itemId/images/remove-thumbnail", middleware.AuthenticatedRoute(), item.RemoveItemImageThumbnail)
	itemsTxRouteGroup.OPTIONS("/:itemId/images/remove-thumbnail")

	itemsTxRouteGroup.DELETE("/:itemId/images", middleware.AuthenticatedRoute(), item.DeleteItemImages)

	itemsTxRouteGroup.GET("/:itemId/images/:imageId", item.GetItemImage)

	itemsTxRouteGroup.PUT("/:itemId/mark-off-bid", middleware.AuthenticatedRoute(), item.MarkItemOffBid)
	itemsTxRouteGroup.OPTIONS("/:itemId/mark-off-bid")

	itemsTxRouteGroup.POST("/:itemId/bid", middleware.AuthenticatedRoute(), item.PlaceBidOnItem)
	itemsTxRouteGroup.OPTIONS("/:itemId/bid")

	itemsTxRouteGroup.GET("/:itemId/bids", item.ListItemBids)
	itemsTxRouteGroup.OPTIONS("/:itemId/bids")

	itemsTxRouteGroup.POST("/:itemId/reaction", middleware.AuthenticatedRoute(), item.AddReactionToItem)
	itemsTxRouteGroup.OPTIONS("/:itemId/reaction")

	itemsTxRouteGroup.DELETE("/:itemId/reaction", middleware.AuthenticatedRoute(), item.RemoveItemReaction)

	itemsTxRouteGroup.POST("/:itemId/comment", middleware.AuthenticatedRoute(), item.AddItemComment)
	itemsTxRouteGroup.OPTIONS("/:itemId/comment")

	itemsTxRouteGroup.PATCH("/:itemId/comment/:commentId", middleware.AuthenticatedRoute(), item.UpdateItemComment)
	itemsTxRouteGroup.OPTIONS("/:itemId/comment/:commentId")

	itemsTxRouteGroup.DELETE("/:itemId/comment/:commentId", middleware.AuthenticatedRoute(), item.DeleteItemComment)
	itemsTxRouteGroup.GET("/:itemId/comments", item.ListItemComments)
	itemsTxRouteGroup.OPTIONS("/:itemId/comments")
}
