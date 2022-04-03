package item

import (
	"fmt"
	"github.com/ashishkumar68/auction-api/actions"
	"github.com/ashishkumar68/auction-api/commands"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func CreateItem(c *gin.Context) {
	var addItemCommand commands.AddNewItemCommand
	if err := c.ShouldBindJSON(&addItemCommand); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	addItemCommand.DB = actions.GetDBConnectionByContext(c)

	bus := commands.NewCommandBus()
	item, err := bus.ExecuteContext(c, &addItemCommand)
	if err != nil {
		log.Println(fmt.Sprintf("Could not save item: %s", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": actions.InternalServerErrMsg})
		return
	}

	c.JSON(http.StatusCreated, item)
}