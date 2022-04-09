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
)

type ItemService interface {
	AddNew(ctx context.Context, form forms.AddNewItemForm) (*models.Item, error)
	PlaceItemBid(ctx context.Context, form forms.PlaceNewItemBidForm) (*models.Bid, error)
}

type ItemServiceImplementor struct {
	itemRepository *repositories.ItemRepository
	bidRepository *repositories.BidRepository
}

func (service *ItemServiceImplementor) AddNew(_ context.Context, form forms.AddNewItemForm) (*models.Item, error) {
	newItem := models.NewItemFromValues(
		form.Name,
		form.Description,
		form.Category,
		form.BrandName,
		form.MarketValue)
	err := service.itemRepository.Save(newItem)
	if err != nil {
		log.Println(fmt.Sprintf("could not create new item: %s", err))
		return nil, err
	}

	return newItem, nil
}

func (service *ItemServiceImplementor) PlaceItemBid(
	_ context.Context,
	form forms.PlaceNewItemBidForm) (*models.Bid, error) {

	item := service.itemRepository.Find(form.ItemId)
	if nil == item || item.IsZero() || item.IsDeleted() {
		return nil, BidItemNotFoundError
	}
	var placedBid *models.Bid
	existingBid := service.bidRepository.FindByItem(item, form.BidUser)

	if nil == existingBid || existingBid.IsZero() || existingBid.IsDeleted() {
		newBid := models.NewBidFromValues(item, form.BidValue, form.BidUser)
		err := service.bidRepository.Save(newBid)
		if err != nil {
			log.Println("could not save bid on item due to err:", err)
			return nil, err
		}
		placedBid = newBid
	} else {
		existingBid.Value = form.BidValue
		err := service.bidRepository.Update(existingBid)
		if err != nil {
			log.Println("could not save bid on item due to err:", err)
			return nil, err
		}
		placedBid = existingBid
	}

	return placedBid, nil
}

func initItemService(conn *gorm.DB) ItemService {
	return &ItemServiceImplementor{
		itemRepository: repositories.NewItemRepository(conn),
		bidRepository: repositories.NewBidRepository(conn),
	}
}