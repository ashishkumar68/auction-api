package services

import (
	"context"
	"fmt"
	"github.com/ashishkumar68/auction-api/forms"
	"github.com/ashishkumar68/auction-api/models"
	"github.com/ashishkumar68/auction-api/repositories"
	"github.com/ashishkumar68/auction-api/utils"
	"gorm.io/gorm"
	"log"
	"os"
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
	AddItemImages(ctx context.Context, form forms.AddItemImagesForm) ([]*models.ItemImage, error)
	RemoveItemImage(ctx context.Context, form forms.RemoveItemImageForm) error
	RemoveItemImages(ctx context.Context, form forms.RemoveItemImagesForm) error
	GetItemImage(ctx context.Context, itemImg *models.ItemImage) (fileName string, filePath string, err error)
	MarkItemImageThumbnail(ctx context.Context, form forms.MarkItemImageThumbnailForm) error
	RemoveItemImageThumbnail(ctx context.Context, form forms.RemoveItemThumbnailForm) error
}

type ItemServiceImplementor struct {
	repository *repositories.Repository
}

func initItemService(repository *repositories.Repository) ItemService {
	itemService := &ItemServiceImplementor{
		repository: repository,
	}

	return itemService
}

func (service *ItemServiceImplementor) DeleteFSItemImages(item *models.Item) error {
	itemDir := fmt.Sprintf("%s/items/%s", utils.GetGlobalUploadsDir(), item.Uuid)

	return os.RemoveAll(itemDir)
}

func (service *ItemServiceImplementor) DeleteFSItemImage(image *models.ItemImage) error {
	itemImagePath := fmt.Sprintf("%s/%s", utils.GetGlobalUploadsDir(), image.Path)
	return os.Remove(itemImagePath)
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
	err := service.repository.Save(newItem)
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
	if !form.Item.IsOwner(*form.ActionUser) {
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
	if !form.Item.IsOwner(*form.ActionUser) {
		return ItemNotOwnedByActionUser
	}
	item := form.Item
	err := item.MarkOffBid()
	if err != nil {
		log.Println("could not put item off bid due to error:", err)
		return err
	}
	err = service.repository.UpdateItem(item)
	if err != nil {
		log.Println("could not put item off bid due to error:", err)
		return err
	}

	return nil
}

func (service *ItemServiceImplementor) AddItemImages(
	_ context.Context,
	form forms.AddItemImagesForm) ([]*models.ItemImage, error) {
	if !form.Item.IsOwner(*form.ActionUser) {
		return nil, ItemNotOwnedByActionUser
	}
	item := form.Item
	itemImages := make([]*models.ItemImage, 0)
	for _, file := range form.ImageFiles {
		newItemImage, err := models.NewItemImageFromMultipartFile(item, file, form.ActionUser)
		if err != nil {
			log.Println(fmt.Sprintf("could not build a new item image due to error: %s", err.Error()))
			return nil, err
		}
		itemImages = append(itemImages, newItemImage)
	}

	err := service.repository.Transaction(func(trx *gorm.DB) error {
		var thumbnail *models.ItemImage
		// remove all existing item images if it's an override
		if form.RemoveExisting {
			deleteErr := service.repository.DeleteItemImages(item)
			if deleteErr != nil {
				return deleteErr
			}
		} else {
			currentImages := service.repository.FindItemImages(item)
			if nil != currentImages && (len(currentImages)+len(itemImages) > models.MaxImagesPerItem) {
				return models.MaxItemImagesReachedErr
			}
			for _, img := range currentImages {
				if img.IsThumbnail && thumbnail == nil {
					thumbnail = img
					break
				}
			}
		}
		for _, image := range itemImages {
			if thumbnail == nil {
				image.IsThumbnail = true
				thumbnail = image
			}
			saveErr := service.repository.Save(image)
			if saveErr != nil {
				log.Println(fmt.Sprintf("found error: %s while saving item image to database", saveErr))
				return saveErr
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	if form.RemoveExisting {
		err = service.DeleteFSItemImages(item)
		if err != nil {
			log.Println("could not delete existing item images due to error", err)
			return nil, err
		}
	}
	for _, itemImg := range itemImages {
		err = utils.SaveUploadedFile(itemImg.MultiPartImgFile, utils.GetFileSystemFilePath(itemImg.Path))
		if err != nil {
			log.Println(fmt.Sprintf("found error: %s while saving item image to file system", err.Error()))
			return nil, err
		}
	}

	return itemImages, nil
}

func (service *ItemServiceImplementor) RemoveItemImage(_ context.Context, form forms.RemoveItemImageForm) error {

	if !form.ItemImage.Item.IsOwner(*form.ActionUser) {
		return ItemNotOwnedByActionUser
	}
	err := service.repository.Delete(form.ItemImage)
	if err != nil {
		log.Println(fmt.Sprintf("could not delete item image from DB due to error: %s", err.Error()))
		return err
	}
	err = service.DeleteFSItemImage(form.ItemImage)
	if err != nil {
		log.Println(fmt.Sprintf("could not delete item image from File system due to error: %s", err.Error()))
		return err
	}

	return nil
}

func (service *ItemServiceImplementor) RemoveItemImages(_ context.Context, form forms.RemoveItemImagesForm) error {

	if !form.Item.IsOwner(*form.ActionUser) {
		return ItemNotOwnedByActionUser
	}
	if err := service.repository.DeleteItemImages(form.Item); err != nil {
		log.Println(fmt.Sprintf("could not delete item:%d images from DB due to error:%s", form.Item.ID, err.Error()))
		return err
	}
	if err := service.DeleteFSItemImages(form.Item); err != nil {
		log.Println(fmt.Sprintf("could not delete item:%d images from FS due to error:%s", form.Item.ID, err.Error()))
		return err
	}

	return nil
}

func (service *ItemServiceImplementor) GetItemImage(_ context.Context, itemImg *models.ItemImage) (string, string, error) {
	filePath := fmt.Sprintf("%s/%s", utils.GetGlobalUploadsDir(), itemImg.Path)
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return "", "", models.ItemImageNotFoundErr
	}

	return itemImg.Name, filePath, nil
}

func (service *ItemServiceImplementor) MarkItemImageThumbnail(
	_ context.Context,
	form forms.MarkItemImageThumbnailForm,
) error {
	if !form.ItemImg.Item.IsOwner(*form.ActionUser) {
		return ItemNotOwnedByActionUser
	}

	currentThumbnail := service.repository.FindItemThumbnail(form.ItemImg.Item)
	form.ItemImg.IsThumbnail = true
	err := service.repository.Transaction(func(trx *gorm.DB) error {
		if currentThumbnail != nil {
			updateErr := service.repository.UpdatesWithMap(currentThumbnail, map[string]interface{}{
				"is_thumbnail": false, "updated_by": form.ActionUser.ID,
			})
			if updateErr != nil {
				return updateErr
			}
		}
		updateErr := service.repository.Update(form.ItemImg)
		if updateErr != nil {
			return updateErr
		}
		return nil
	})
	if err != nil {
		log.Println("could not mark item image thumbnail due to error:", err.Error())
		return err
	}

	return nil
}

func (service *ItemServiceImplementor) RemoveItemImageThumbnail(
	_ context.Context,
	form forms.RemoveItemThumbnailForm,
) error {
	if !form.Item.IsOwner(*form.ActionUser) {
		return ItemNotOwnedByActionUser
	}
	currentThumbnail := service.repository.FindItemThumbnail(form.Item)
	if currentThumbnail == nil {
		return models.ItemThumbnailNotFoundErr
	}
	err := service.repository.UpdatesWithMap(currentThumbnail, map[string]interface{}{
		"is_thumbnail": false, "updated_by": form.ActionUser.ID,
	})
	if err != nil {
		log.Println("could not unmark/remove item image thumbnail due to error:", err.Error())
		return err
	}

	return nil
}
