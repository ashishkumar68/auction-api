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

func AddItemComment(c *gin.Context) {
	itemId, err := strconv.Atoi(c.Param("itemId"))
	if err != nil {
		log.Println(fmt.Sprintf("Could not save item comment: %s", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": actions.InvalidItemIdReceivedErr})
		return
	}
	db := actions.GetDBConnectionByContext(c)
	repository := repositories.NewRepository(db)
	item := repository.FindItemById(uint(itemId))
	if item == nil {
		log.Println(fmt.Sprintf("Could not save item comment: %s", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": actions.InvalidItemIdReceivedErr})
		return
	}
	var addItemComment forms.AddItemCommentForm
	addItemComment.Item = item
	addItemComment.ActionUser = actions.GetActionUserByContext(c)
	if err = c.ShouldBindJSON(&addItemComment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	itemCommentService := services.NewItemCommentService(db)
	newComment, err := itemCommentService.AddItemComment(c, addItemComment)
	if err != nil {
		log.Println(fmt.Sprintf("Could not save item comment: %s", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": actions.InternalServerErrMsg})
		return
	}

	c.JSON(http.StatusCreated, newComment)
}

func UpdateItemComment(c *gin.Context) {
	var updateItemComment forms.EditItemCommentForm
	updateItemComment.ActionUser = actions.GetActionUserByContext(c)
	itemId, err := strconv.Atoi(c.Param("itemId"))
	if err != nil {
		log.Println(fmt.Sprintf("Could not update item comment: %s", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": actions.InvalidItemIdReceivedErr})
		return
	}
	db := actions.GetDBConnectionByContext(c)
	repository := repositories.NewRepository(db)
	item := repository.FindItemById(uint(itemId))
	if item == nil {
		log.Println(fmt.Sprintf("Could not update item comment: %s", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": actions.InvalidItemIdReceivedErr})
		return
	}
	updateItemComment.Item = item
	commentId, err := strconv.Atoi(c.Param("commentId"))
	if err != nil {
		log.Println(fmt.Sprintf("Could not update item comment: %s", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": actions.InvalidCommentIdReceivedErr})
		return
	}
	updateItemComment.CommentId = uint(commentId)
	comment := repository.FindItemCommentById(uint(commentId))
	if comment == nil {
		log.Println(fmt.Sprintf("Could not update item comment: %s", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": actions.InvalidCommentIdReceivedErr})
		return
	}
	if !comment.UserCreated.IsSameAs(updateItemComment.ActionUser.BaseModel) {
		log.Println(fmt.Sprintf("Could not update item comment: %s", err))
		c.JSON(http.StatusForbidden, gin.H{"error": actions.CommentNotAuthoredByUser})
		return
	}
	if err = c.ShouldBindJSON(&updateItemComment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	itemCommentService := services.NewItemCommentService(db)
	err = itemCommentService.UpdateItemComment(c, updateItemComment)
	if err != nil {
		log.Println(fmt.Sprintf("Could not update item comment: %s", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": actions.InternalServerErrMsg})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func DeleteItemComment(c *gin.Context) {
	var deleteItemComment forms.DeleteItemCommentForm
	deleteItemComment.ActionUser = actions.GetActionUserByContext(c)
	itemId, err := strconv.Atoi(c.Param("itemId"))
	if err != nil {
		log.Println(fmt.Sprintf("Could not delete item comment: %s", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": actions.InvalidItemIdReceivedErr})
		return
	}
	db := actions.GetDBConnectionByContext(c)
	repository := repositories.NewRepository(db)
	item := repository.FindItemById(uint(itemId))
	if item == nil {
		log.Println(fmt.Sprintf("Could not delete item comment: %s", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": actions.InvalidItemIdReceivedErr})
		return
	}
	deleteItemComment.Item = item
	commentId, err := strconv.Atoi(c.Param("commentId"))
	if err != nil {
		log.Println(fmt.Sprintf("Could not delete item comment: %s", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": actions.InvalidCommentIdReceivedErr})
		return
	}
	deleteItemComment.CommentId = uint(commentId)
	comment := repository.FindItemCommentById(uint(commentId))
	if comment == nil {
		log.Println(fmt.Sprintf("Could not delete item comment: %s", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": actions.InvalidCommentIdReceivedErr})
		return
	}
	if !comment.UserCreated.IsSameAs(deleteItemComment.ActionUser.BaseModel) {
		log.Println(fmt.Sprintf("Could not delete item comment: %s", err))
		c.JSON(http.StatusForbidden, gin.H{"error": actions.CommentNotAuthoredByUser})
		return
	}

	itemCommentService := services.NewItemCommentService(db)
	err = itemCommentService.DeleteItemComment(c, deleteItemComment)
	if err != nil {
		log.Println(fmt.Sprintf("Could not delete item comment: %s", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": actions.InternalServerErrMsg})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func ListItemComments(c *gin.Context) {
	itemId, err := strconv.Atoi(c.Param("itemId"))
	if err != nil {
		log.Println(fmt.Sprintf("Could not fetch item comments due to error: %s", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": actions.InvalidItemIdReceivedErr})
		return
	}
	db := actions.GetDBConnectionByContext(c)
	repository := repositories.NewRepository(db)
	item := repository.FindItemById(uint(itemId))
	if item == nil {
		log.Println(fmt.Sprintf("Could not fetch item comments due to error: %s", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": actions.InvalidItemIdReceivedErr})
		return
	}
	pg := paginate.New()

	c.JSON(http.StatusOK, pg.With(repository.FindCommentsByItem(item)).Request(c.Request).Response(&[]models.ItemComment{}))
}
