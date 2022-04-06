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

var (
	BidItemNotFoundError = fmt.Errorf("bid item was not found")
)

type PlaceNewItemBidCommand struct {
	DB *gorm.DB

	ItemId		uint			`json:"itemId" binding:"required,min=1,max=12"`
	BidValue	models.Value	`json:"bidValue" binding:"required,min=1"`
	BidUser		*models.User
}

func (cmd *PlaceNewItemBidCommand) Type() string {
	return "PlaceNewItemBid"
}

func PlaceNewItemBidHandler(ctx context.Context, command cbus.Command) (interface{}, error) {
	cmd := command.(*PlaceNewItemBidCommand)
	itemRepo := repositories.NewItemRepository(cmd.DB)
	item := itemRepo.Find(cmd.ItemId)
	if nil == item || item.IsZero() || item.IsDeleted() {
		return nil, BidItemNotFoundError
	}
	var placedBid *models.Bid
	bidRepo := repositories.NewBidRepository(cmd.DB)
	existingBid := bidRepo.FindByItem(item, cmd.BidUser)
	if nil == existingBid || existingBid.IsZero() || existingBid.IsDeleted() {
		newBid := models.NewBidFromValues(item, cmd.BidValue)
		err := bidRepo.Save(newBid)
		if err != nil {
			log.Println("could not save bid on item due to err:", err)
			return nil, err
		}
		placedBid = newBid
	} else {
		existingBid.Value = cmd.BidValue
		err := bidRepo.Update(existingBid)
		if err != nil {
			log.Println("could not save bid on item due to err:", err)
			return nil, err
		}
		placedBid = existingBid
	}

	return placedBid, nil
}