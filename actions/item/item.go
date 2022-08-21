package item

import (
	"fmt"
	"github.com/ashishkumar68/auction-api/actions"
	"github.com/ashishkumar68/auction-api/forms"
	"github.com/ashishkumar68/auction-api/models"
	"github.com/ashishkumar68/auction-api/repositories"
	"github.com/ashishkumar68/auction-api/services"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/morkid/paginate"
	"log"
	"net/http"
	"strconv"
)

func CreateItem(c *gin.Context) {
	var addItemForm forms.AddNewItemForm

	addItemForm.ActionUser = actions.GetActionUserByContext(c)
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
	repository := repositories.NewRepository(db)

	var items []*models.Item
	page := pg.With(repository.ListItems()).Request(c.Request).Response(&items)
	itemReactionMap := repository.FindReactionsCountByItems(items)
	for _, item := range items {
		if itemReactions, ok := itemReactionMap[item.ID]; ok {
			item.Reactions = itemReactions
		}
		item.CategoryText = models.GetItemCategoryString(item.Category)
		for _, itemImg := range item.ItemImages {
			itemImg.Url = BuildItemImageRoute(*itemImg)
		}
	}

	c.JSON(http.StatusOK, page)
}

func EditItem(c *gin.Context) {
	var editItemForm forms.EditItemForm
	editItemForm.ActionUser = actions.GetActionUserByContext(c)
	itemId, err := strconv.Atoi(c.Param("itemId"))
	if err != nil {
		log.Println(fmt.Sprintf("Could not update item: %s", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": actions.InvalidItemIdReceivedErr})
		return
	}
	db := actions.GetDBConnectionByContext(c)
	repository := repositories.NewRepository(db)
	item := repository.FindItemById(uint(itemId))
	if item == nil {
		log.Println(fmt.Sprintf("Could not update item: %s", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": actions.InvalidItemIdReceivedErr})
		return
	}
	editItemForm.Item = item
	if err = c.ShouldBindJSON(&editItemForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	itemService := services.NewItemService(db)
	err = itemService.EditItem(c, editItemForm)
	if err != nil {
		log.Println(fmt.Sprintf("Could not edit item: %s", err))
		if err == services.ItemNotOwnedByActionUser {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": actions.InternalServerErrMsg})
		}
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func MarkItemOffBid(c *gin.Context) {
	var form forms.MarkItemOffBidForm
	form.ActionUser = actions.GetActionUserByContext(c)
	itemId, err := strconv.Atoi(c.Param("itemId"))
	if err != nil {
		log.Println(fmt.Sprintf("Could not put item off bid due to err: %s", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": actions.InvalidItemIdReceivedErr})
		return
	}
	db := actions.GetDBConnectionByContext(c)
	repository := repositories.NewRepository(db)
	item := repository.FindItemById(uint(itemId))
	if item == nil {
		log.Println(fmt.Sprintf("Could not put item off bid due to err: %s", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": actions.InvalidItemIdReceivedErr})
		return
	}
	form.Item = item
	itemService := services.NewItemService(db)
	err = itemService.MarkItemOffBid(c, form)
	if err != nil {
		log.Println(fmt.Sprintf("Could not put item off bid due to error: %s", err))
		if err == services.ItemNotOwnedByActionUser {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": actions.InternalServerErrMsg})
		}
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func AddItemImages(c *gin.Context) {
	var form forms.AddItemImagesForm
	form.ActionUser = actions.GetActionUserByContext(c)
	itemId, err := strconv.Atoi(c.Param("itemId"))
	if err != nil {
		log.Println(fmt.Sprintf("Could not add item images due to err: %s", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": actions.InvalidItemIdReceivedErr})
		return
	}
	if c.Query("removeExisting") != "" {
		removeExisting, err := strconv.ParseBool(c.Query("removeExisting"))
		if err != nil {
			log.Println(fmt.Sprintf("Could not add item images due to err: %s", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": actions.InvalidRemoveExistingVal})
			return
		}
		form.RemoveExisting = removeExisting
	}
	db := actions.GetDBConnectionByContext(c)
	repository := repositories.NewRepository(db)
	item := repository.FindItemById(uint(itemId))
	if item == nil {
		log.Println(fmt.Sprintf("Could not add item images due to err: %s", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": actions.InvalidItemIdReceivedErr})
		return
	}
	form.Item = item

	if err = c.MustBindWith(&form, binding.FormMultipart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	itemService := services.NewItemService(db)
	itemImages, err := itemService.AddItemImages(c, form)
	if err != nil {
		log.Println(fmt.Sprintf("Could not add item images due to error: %s", err))
		if err == models.MaxItemImagesReachedErr {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else if err == services.ItemNotOwnedByActionUser {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": actions.InternalServerErrMsg})
		}
		return
	}
	for _, itemImg := range itemImages {
		itemImg.Url = BuildItemImageRoute(*itemImg)
	}

	c.JSON(http.StatusCreated, itemImages)
}

func DeleteItemImage(c *gin.Context) {
	var form forms.RemoveItemImageForm
	form.ActionUser = actions.GetActionUserByContext(c)

	itemId, err := strconv.Atoi(c.Param("itemId"))
	if err != nil {
		log.Println(fmt.Sprintf("Could not delete item image due to err: %s", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": actions.InvalidItemIdReceivedErr})
		return
	}
	imageId, err := strconv.Atoi(c.Param("imageId"))
	if err != nil {
		log.Println(fmt.Sprintf("Could not delete item image due to err: %s", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": actions.InvalidImageIdFoundErr})
		return
	}
	db := actions.GetDBConnectionByContext(c)
	repository := repositories.NewRepository(db)
	itemImage := repository.FindItemImage(uint(imageId), uint(itemId))
	if itemImage == nil {
		log.Println(fmt.Sprintf("Could not delete item image due to err: %s", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": actions.InvalidImageIdFoundErr})
		return
	}

	form.ItemImage = itemImage
	itemService := services.NewItemService(db)
	err = itemService.RemoveItemImage(c, form)
	if err != nil {
		log.Println(fmt.Sprintf("Could not delete item image due to err: %s", err))
		if err == services.ItemNotOwnedByActionUser {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": actions.InternalServerErrMsg})
		}
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func DeleteItemImages(c *gin.Context) {
	var form forms.RemoveItemImagesForm
	form.ActionUser = actions.GetActionUserByContext(c)

	itemId, err := strconv.Atoi(c.Param("itemId"))
	if err != nil {
		log.Println(fmt.Sprintf("Could not delete item images due to err: %s", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": actions.InvalidItemIdReceivedErr})
		return
	}
	db := actions.GetDBConnectionByContext(c)
	repository := repositories.NewRepository(db)
	item := repository.FindItemById(uint(itemId))
	if item == nil {
		log.Println(fmt.Sprintf("Could not delete item images due to err: %s", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": actions.InvalidItemIdReceivedErr})
		return
	}

	form.Item = item
	itemService := services.NewItemService(db)
	err = itemService.RemoveItemImages(c, form)
	if err != nil {
		log.Println(fmt.Sprintf("Could not delete item images due to err: %s", err))
		if err == services.ItemNotOwnedByActionUser {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": actions.InternalServerErrMsg})
		}
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func GetItemImage(c *gin.Context) {
	itemId, err := strconv.Atoi(c.Param("itemId"))
	if err != nil {
		log.Println(fmt.Sprintf("Could not find item image due to err: %s", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": actions.InvalidItemIdReceivedErr})
		return
	}
	db := actions.GetDBConnectionByContext(c)
	repository := repositories.NewRepository(db)
	imageId, err := strconv.Atoi(c.Param("imageId"))
	if err != nil {
		log.Println(fmt.Sprintf("Could not delete item image due to err: %s", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": actions.InvalidImageIdFoundErr})
		return
	}
	itemImg := repository.FindItemImage(uint(imageId), uint(itemId))
	if itemImg == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": actions.ItemImageNotFoundError})
		return
	}
	itemService := services.NewItemService(db)
	_, filePath, err := itemService.GetItemImage(c, itemImg)
	if err != nil {
		log.Println(fmt.Sprintf("Could not find item image due to err: %s", err.Error()))
		c.JSON(http.StatusNotFound, gin.H{"error": models.ItemImageNotFoundErr.Error()})
		return
	}

	c.Header("Content-Type", "application/octet-stream")
	// don't want to force download.
	//c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache")
	c.File(filePath)
}

func MakeItemImageThumbnail(c *gin.Context) {
	itemId, err := strconv.Atoi(c.Param("itemId"))
	if err != nil {
		log.Println(fmt.Sprintf("Could not mark item image thubmnail due to err: %s", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": actions.InvalidItemIdReceivedErr})
		return
	}
	db := actions.GetDBConnectionByContext(c)
	repository := repositories.NewRepository(db)
	imageId, err := strconv.Atoi(c.Param("imageId"))
	if err != nil {
		log.Println(fmt.Sprintf("Could not mark item image thumbnail due to err: %s", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": actions.InvalidImageIdFoundErr})
		return
	}
	itemImg := repository.FindItemImage(uint(imageId), uint(itemId))
	if itemImg == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": actions.ItemImageNotFoundError})
		return
	}
	var form forms.MarkItemImageThumbnailForm
	form.ActionUser = actions.GetActionUserByContext(c)
	form.ItemImg = itemImg

	itemService := services.NewItemService(db)
	err = itemService.MarkItemImageThumbnail(c, form)
	if err != nil {
		log.Println(fmt.Sprintf("Could not mark item image thumbnail due to err: %s", err))
		if err == services.ItemNotOwnedByActionUser {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": actions.InternalServerErrMsg})
		}
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func RemoveItemImageThumbnail(c *gin.Context) {
	itemId, err := strconv.Atoi(c.Param("itemId"))
	if err != nil {
		log.Println(fmt.Sprintf("Could not remove item image thubmnail due to err: %s", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": actions.InvalidItemIdReceivedErr})
		return
	}
	db := actions.GetDBConnectionByContext(c)
	repository := repositories.NewRepository(db)
	item := repository.FindItemById(uint(itemId))
	if item == nil {
		log.Println(fmt.Sprintf("Could not remove item image thumbnail due to err: %s", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": actions.InvalidItemIdReceivedErr})
		return
	}
	var form forms.RemoveItemThumbnailForm
	form.ActionUser = actions.GetActionUserByContext(c)
	form.Item = item

	itemService := services.NewItemService(db)
	err = itemService.RemoveItemImageThumbnail(c, form)
	if err != nil {
		log.Println(fmt.Sprintf("Could not remove item images due to err: %s", err))
		if err == services.ItemNotOwnedByActionUser {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": actions.InternalServerErrMsg})
		}
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
