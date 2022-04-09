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
)

func CreateItem(c *gin.Context) {
	var addItemForm forms.AddNewItemForm
	if err := c.ShouldBindJSON(&addItemForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	itemService := services.NewItemService(actions.GetDBConnectionByContext(c))
	item, err := itemService.AddNew(c, addItemForm)
	if err != nil {
		log.Println(fmt.Sprintf("Could not save item: %s", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": actions.InternalServerErrMsg})
		return
	}

	c.JSON(http.StatusCreated, item)
}

func ListItems(c *gin.Context) {
	pg := paginate.New()
	db := actions.GetDBConnectionByContext(c)
	itemRepo := repositories.NewItemRepository(db)

	c.JSON(http.StatusOK, pg.With(itemRepo.List()).Request(c.Request).Response(&[]models.Item{}))
	return
}
