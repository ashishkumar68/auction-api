package commands

import (
	"context"
	"fmt"
	"github.com/ashishkumar68/auction-api/models"
	"github.com/ashishkumar68/auction-api/repositories"
	"github.com/gogolfing/cbus"
	"gorm.io/gorm"
	"log"
)

type AddNewItemCommand struct {
	DB *gorm.DB

	Name        string              `json:"name" binding:"required"`
	Description string              `json:"description" binding:"required"`
	Category    models.ItemCategory `json:"category" binding:"required,oneof=0 1 2 3"`
	BrandName   string              `json:"brandName" binding:"required"`
	MarketValue models.Value        `json:"marketValue" binding:"required"`
}

func (cmd *AddNewItemCommand) Type() string {
	return "AddNewItem"
}

func AddNewItemHandler(ctx context.Context, command cbus.Command) (interface{}, error) {
	addNewItemCmd := command.(*AddNewItemCommand)
	newItem := models.NewItemFromValues(
		addNewItemCmd.Name,
		addNewItemCmd.Description,
		addNewItemCmd.Category,
		addNewItemCmd.BrandName,
		addNewItemCmd.MarketValue)
	err := repositories.NewItemRepository(addNewItemCmd.DB).Save(newItem)
	if err != nil {
		log.Println(fmt.Sprintf("could not create new item: %s", err))
		return nil, err
	}

	return newItem, nil
}
