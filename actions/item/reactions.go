package item

import (
	"fmt"
	"github.com/ashishkumar68/auction-api/actions"
	"github.com/ashishkumar68/auction-api/forms"
	"github.com/ashishkumar68/auction-api/repositories"
	"github.com/ashishkumar68/auction-api/services"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func AddReactionToItem(c *gin.Context) {
	itemId, err := strconv.Atoi(c.Param("itemId"))
	if err != nil {
		log.Println("could not parse item id:", c.Param("itemId"))
		c.JSON(http.StatusBadRequest, gin.H{"error": actions.InvalidItemIdReceivedErr})
		return
	}
	db := actions.GetDBConnectionByContext(c)
	item := repositories.NewRepository(db).FindItemById(uint(itemId))
	if item == nil {
		log.Println("could not find item with id:", c.Param("itemId"))
		c.JSON(http.StatusBadRequest, gin.H{"error": actions.InvalidItemIdReceivedErr})
		return
	}
	var addItemReaction forms.AddItemReactionForm
	addItemReaction.ActionUser = actions.GetActionUserByContext(c)
	addItemReaction.Item = item
	if err = c.ShouldBindJSON(&addItemReaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	reactionService := services.NewReactionService(db)
	reaction, err := reactionService.AddReactionToItem(c, addItemReaction)
	if err != nil {
		log.Println(fmt.Sprintf("Could not add reaction to item: %d due to error: %s", item.ID, err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": actions.InternalServerErrMsg})
		return
	}

	c.JSON(http.StatusCreated, reaction)
}

func RemoveItemReaction(c *gin.Context) {
	itemId, err := strconv.Atoi(c.Param("itemId"))
	if err != nil {
		log.Println("could not parse item id:", c.Param("itemId"))
		c.JSON(http.StatusBadRequest, gin.H{"error": actions.InvalidItemIdReceivedErr})
		return
	}
	db := actions.GetDBConnectionByContext(c)
	item := repositories.NewRepository(db).FindItemById(uint(itemId))
	if item == nil {
		log.Println("could not find item with id:", c.Param("itemId"))
		c.JSON(http.StatusBadRequest, gin.H{"error": actions.InvalidItemIdReceivedErr})
		return
	}
	deleteItemReaction := forms.NewRemoveItemReactionForm(actions.GetActionUserByContext(c), item)
	reactionService := services.NewReactionService(db)
	err = reactionService.RemoveReactionFromItem(c, *deleteItemReaction)
	if err != nil {
		log.Println(fmt.Sprintf("Could not remove reaction to item: %d due to error: %s", item.ID, err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": actions.InternalServerErrMsg})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
