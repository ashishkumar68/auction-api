package item

import (
	"fmt"
	"github.com/ashishkumar68/auction-api/actions"
	"github.com/ashishkumar68/auction-api/models"
	"github.com/gin-gonic/gin"
	"os"
)

const (
	BaseItemsRoute = "/items"
)

func MapItemRoutes(itemsGroup *gin.RouterGroup) {
	itemsTxRouteGroup := itemsGroup.Group(BaseItemsRoute, actions.TransactionRoute(), actions.CorsRoute())

	itemsTxRouteGroup.POST("", actions.AuthenticatedRoute(), CreateItem)
	itemsTxRouteGroup.GET("", ListItems)
	itemsTxRouteGroup.OPTIONS("")

	itemsTxRouteGroup.PATCH("/:itemId", actions.AuthenticatedRoute(), EditItem)
	itemsTxRouteGroup.OPTIONS("/:itemId")

	itemsTxRouteGroup.POST("/:itemId/images", actions.AuthenticatedRoute(), AddItemImages)
	itemsTxRouteGroup.OPTIONS("/:itemId/images")

	itemsTxRouteGroup.DELETE("/:itemId/images/:imageId", actions.AuthenticatedRoute(), DeleteItemImage)
	itemsTxRouteGroup.OPTIONS("/:itemId/images/:imageId")

	itemsTxRouteGroup.PATCH("/:itemId/images/:imageId/make-thumbnail", actions.AuthenticatedRoute(), MakeItemImageThumbnail)
	itemsTxRouteGroup.OPTIONS("/:itemId/images/:imageId/make-thumbnail")

	itemsTxRouteGroup.DELETE("/:itemId/images/remove-thumbnail", actions.AuthenticatedRoute(), RemoveItemImageThumbnail)
	itemsTxRouteGroup.OPTIONS("/:itemId/images/remove-thumbnail")

	itemsTxRouteGroup.DELETE("/:itemId/images", actions.AuthenticatedRoute(), DeleteItemImages)

	itemsTxRouteGroup.GET("/:itemId/images/:imageId", GetItemImage)

	itemsTxRouteGroup.PUT("/:itemId/mark-off-bid", actions.AuthenticatedRoute(), MarkItemOffBid)
	itemsTxRouteGroup.OPTIONS("/:itemId/mark-off-bid")

	itemsTxRouteGroup.POST("/:itemId/bid", actions.AuthenticatedRoute(), PlaceBidOnItem)
	itemsTxRouteGroup.OPTIONS("/:itemId/bid")

	itemsTxRouteGroup.GET("/:itemId/bids", ListItemBids)
	itemsTxRouteGroup.OPTIONS("/:itemId/bids")

	itemsTxRouteGroup.POST("/:itemId/reaction", actions.AuthenticatedRoute(), AddReactionToItem)
	itemsTxRouteGroup.OPTIONS("/:itemId/reaction")

	itemsTxRouteGroup.DELETE("/:itemId/reaction", actions.AuthenticatedRoute(), RemoveItemReaction)

	itemsTxRouteGroup.POST("/:itemId/comment", actions.AuthenticatedRoute(), AddItemComment)
	itemsTxRouteGroup.OPTIONS("/:itemId/comment")

	itemsTxRouteGroup.PATCH("/:itemId/comment/:commentId", actions.AuthenticatedRoute(), UpdateItemComment)
	itemsTxRouteGroup.OPTIONS("/:itemId/comment/:commentId")

	itemsTxRouteGroup.DELETE("/:itemId/comment/:commentId", actions.AuthenticatedRoute(), DeleteItemComment)
	itemsTxRouteGroup.GET("/:itemId/comments", ListItemComments)
	itemsTxRouteGroup.OPTIONS("/:itemId/comments")
}

func BuildItemImageRoute(img models.ItemImage) string {
	return fmt.Sprintf("%s://%s%s%s/%d/images/%d", "http", os.Getenv("HOST"), actions.BaseApiRoute, BaseItemsRoute, img.ItemId, img.ID)
}
