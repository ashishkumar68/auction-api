package item

import (
	"fmt"
	"github.com/ashishkumar68/auction-api/actions"
	"github.com/ashishkumar68/auction-api/forms"
	"github.com/ashishkumar68/auction-api/services"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func PlaceBidOnItem(c *gin.Context) {
	var placeBidForm forms.PlaceNewItemBidForm
	itemId, err := strconv.Atoi(c.Param("itemId"))
	placeBidForm.ActionUser = actions.GetActionUserByContext(c)

	if err != nil {
		log.Println(fmt.Sprintf("Could not save bid: %s", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": actions.InvalidItemIdReceivedErr})
		return
	}

	placeBidForm.ItemId = uint(itemId)
	if err = c.ShouldBindJSON(&placeBidForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	placeBidForm.BidUserId = actions.GetActionUserByContext(c).GetId()
	itemService := services.NewItemService(actions.GetDBConnectionByContext(c))

	bid, err := itemService.PlaceItemBid(c, placeBidForm)
	if err != nil {
		log.Println(fmt.Sprintf("Could not save bid: %s", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": actions.InternalServerErrMsg})
		return
	}

	c.JSON(http.StatusCreated, bid)
}
