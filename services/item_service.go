package services

import (
	"context"
	"fmt"
	"github.com/ashishkumar68/auction-api/forms"
	"github.com/ashishkumar68/auction-api/models"
	"github.com/ashishkumar68/auction-api/repositories"
	"gorm.io/gorm"
	"log"
)

var (
	BidItemNotFoundError = fmt.Errorf("bid item was not found")
	BidUserNotFoundError = fmt.Errorf("bid user details was not found")
)

type ItemService interface {
	AddNew(ctx context.Context, form forms.AddNewItemForm) (*models.Item, error)
	PlaceItemBid(ctx context.Context, form forms.PlaceNewItemBidForm) (*models.Bid, error)
}

type ItemServiceImplementor struct {
	repository *repositories.Repository
}

func (service *ItemServiceImplementor) AddNew(_ context.Context, form forms.AddNewItemForm) (*models.Item, error) {
	newItem := models.NewItemFromValues(
		form.Name,
		form.Description,
		form.Category,
		form.BrandName,
		form.MarketValue,
		form.LastBidDate,
	)
	newItem.UserCreatedBy = form.ActionUser.ID
	newItem.UserUpdatedBy = form.ActionUser.ID
	err := service.repository.SaveItem(newItem)
	if err != nil {
		log.Println(fmt.Sprintf("could not create new item: %s", err))
		return nil, err
	}

	return newItem, nil
}

func (service *ItemServiceImplementor) PlaceItemBid(
	ctx context.Context,
	form forms.PlaceNewItemBidForm) (*models.Bid, error) {

	bidUser := service.repository.FindUserById(form.BidUserId)
	if nil == bidUser {
		return nil, BidUserNotFoundError
	}
	item := service.repository.FindItemById(form.ItemId)
	if nil == item {
		return nil, BidItemNotFoundError
	}

	var placedBid *models.Bid
	existingBid := service.repository.FindBidByItem(item, bidUser)
	if nil == existingBid {
		newBid := models.NewBidFromValues(item, form.BidValue, form.ActionUser)
		err := service.repository.SaveBid(newBid)
		if err != nil {
			log.Println("could not save bid on item due to err:", err)
			return nil, err
		}
		placedBid = newBid
	} else {
		existingBid.Value = form.BidValue
		existingBid.UserUpdatedBy = form.ActionUser.ID
		err := service.repository.UpdateBid(existingBid)
		if err != nil {
			log.Println("could not save bid on item due to err:", err)
			return nil, err
		}
		placedBid = existingBid
	}

	return placedBid, nil
}

func initItemService(conn *gorm.DB) ItemService {
	itemService := &ItemServiceImplementor{
		repository: repositories.NewRepository(conn),
	}

	return itemService
}
