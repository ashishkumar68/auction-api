package item

import (
	"fmt"
	"github.com/ashishkumar68/auction-api/actions"
	"github.com/ashishkumar68/auction-api/forms"
	"github.com/ashishkumar68/auction-api/models"
	"github.com/ashishkumar68/auction-api/repositories"
	"github.com/ashishkumar68/auction-api/services"
	"github.com/gin-gonic/gin"
	"github.com/morkid/paginate"
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
		if err == services.ItemNotBidEligible {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else if err == services.BidsNotAllowedByOwner {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": actions.InternalServerErrMsg})
		}
		return
	}

	c.JSON(http.StatusCreated, bid)
}

func ListItemBids(c *gin.Context) {
	itemId, err := strconv.Atoi(c.Param("itemId"))
	if err != nil {
		log.Println(fmt.Sprintf("Could not fetch item bids due to error: %s", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": actions.InvalidItemIdReceivedErr})
		return
	}
	db := actions.GetDBConnectionByContext(c)
	repository := repositories.NewRepository(db)
	item := repository.FindItemById(uint(itemId))
	if item == nil {
		log.Println(fmt.Sprintf("Could not fetch item bids due to error: %s", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": actions.InvalidItemIdReceivedErr})
		return
	}
	pg := paginate.New()

	c.JSON(http.StatusOK, pg.With(repository.FindBidsByItem(item)).Request(c.Request).Response(&[]models.Bid{}))
}
