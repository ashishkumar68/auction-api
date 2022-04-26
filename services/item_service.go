package services

import (
	"context"
	"fmt"
	"github.com/ashishkumar68/auction-api/forms"
	"github.com/ashishkumar68/auction-api/models"
	"github.com/ashishkumar68/auction-api/repositories"
	"log"
)

var (
	BidItemNotFoundError     = fmt.Errorf("bid item was not found")
	BidUserNotFoundError     = fmt.Errorf("bid user details was not found")
	ItemNotBidEligible       = fmt.Errorf("this item is not eligible for bidding")
	ItemNotOwnedByActionUser = fmt.Errorf("item is not owned by action user")
	BidsNotAllowedByOwner    = fmt.Errorf("bids can not be created by item owners")
)

type ItemService interface {
	AddNew(ctx context.Context, form forms.AddNewItemForm) (*models.Item, error)
	EditItem(ctx context.Context, form forms.EditItemForm) error
	MarkItemOffBid(ctx context.Context, form forms.MarkItemOffBidForm) error
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
	if item.IsOwner(*bidUser) {
		return nil, BidsNotAllowedByOwner
	}
	if !item.IsBidEligible() {
		return nil, ItemNotBidEligible
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

func (service *ItemServiceImplementor) EditItem(_ context.Context, form forms.EditItemForm) error {
	if !form.Item.UserCreated.IsSameAs(form.ActionUser.BaseModel) {
		return ItemNotOwnedByActionUser
	}
	item := form.Item
	err := item.UpdateFromValues(
		form.Name,
		form.Description,
		form.Category,
		form.BrandName,
		form.MarketValue,
		form.LastBidDate,
		form.ActionUser,
	)
	if err != nil {
		log.Println("could not edit item due to error:", err)
		return err
	}
	err = service.repository.UpdateItem(item)
	if err != nil {
		log.Println("could not edit item due to error:", err)
		return err
	}

	return nil
}

func (service *ItemServiceImplementor) MarkItemOffBid(_ context.Context, form forms.MarkItemOffBidForm) error {
	if !form.Item.UserCreated.IsSameAs(form.ActionUser.BaseModel) {
		return ItemNotOwnedByActionUser
	}
	item := form.Item
	err := item.MarkOffBid()
	if err != nil {
		log.Println("could not put item off bid due to error:", err)
		return err
	}

	return nil
}

func initItemService(repository *repositories.Repository) ItemService {
	itemService := &ItemServiceImplementor{
		repository: repository,
	}

	return itemService
}
