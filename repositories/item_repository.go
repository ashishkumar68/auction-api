package repositories

import (
	"fmt"
	"github.com/ashishkumar68/auction-api/models"
	"gorm.io/gorm"
	"log"
)

func (repo *Repository) FindItemById(id uint) *models.Item {
	var item models.Item
	repo.connection.
		Joins("LEFT JOIN item_images ON item_images.item_id = items.id").
		Preload("UserCreated").
		Preload("UserUpdated").
		Preload("ItemImages").
		Find(&item, id)
	if item.IsZero() {
		return nil
	}

	return &item
}

func (repo *Repository) SaveItem(item *models.Item) error {
	result := repo.connection.Create(item)
	if result.Error != nil {
		log.Println(fmt.Sprintf("Could not insert new record for type: %T", item))
		log.Println(fmt.Sprintf("Insert error: %s", result.Error))
		return result.Error
	}

	return nil
}

func (repo *Repository) UpdateItem(item *models.Item) error {
	result := repo.connection.Save(item)
	if result.Error != nil {
		log.Println(fmt.Sprintf("Could not update record for type: %T", item))
		log.Println(fmt.Sprintf("Update error: %s", result.Error))
		return result.Error
	}

	return nil
}

func (repo *Repository) FindItemByUuid(uuid string) *models.Item {
	var item models.Item
	repo.connection.Joins("UserCreated").Where("items.uuid = ?", uuid).First(&item)
	if item.IsZero() {
		return nil
	}

	return &item
}

func (repo *Repository) FindItemByName(name string) []models.Item {
	var items []models.Item
	repo.connection.Where("name LIKE ?", "%"+name+"%").Find(&items)

	return items
}

func (repo *Repository) ListItems() *gorm.DB {
	return repo.connection.
		Preload("ItemImages").
		Joins("UserCreated").
		Joins("UserUpdated").
		Joins("LEFT JOIN item_images ON item_images.item_id = items.id AND item_images.is_thumbnail = 1").
		Model(&models.Item{})
}

func (repo *Repository) ListUserItems(user *models.User) *gorm.DB {
	return repo.connection.
		Model(&models.Item{}).
		Preload("ItemImages").
		Joins("UserCreated").
		Joins("LEFT JOIN item_images ON item_images.item_id = items.id AND item_images.is_thumbnail = 1").
		Where("items.created_by = ?", user.ID).
		Order("items.id DESC")
}

func (repo *Repository) FindItemImages(item *models.Item) []*models.ItemImage {
	var itemImages []*models.ItemImage
	repo.connection.
		Model(&models.ItemImage{}).
		Where("item_id = ?", item.ID).Find(&itemImages)
	if len(itemImages) == 0 {
		return nil
	}

	return itemImages
}

func (repo *Repository) DeleteItemImages(item *models.Item) error {
	result := repo.connection.Where("item_id = ?", item.ID).Delete(&models.ItemImage{})
	if result.Error != nil {
		log.Println("could not delete item images due to error:", result.Error)
		return result.Error
	}

	return nil
}

func (repo *Repository) FindItemImage(imageId uint, itemId uint) *models.ItemImage {
	var image models.ItemImage
	repo.connection.
		Preload("Item").
		Preload("Item.UserCreated").
		Joins("JOIN items ON item_images.item_id = items.id").
		Joins("JOIN users ON items.created_by = users.id").
		Where("item_images.id = ? AND items.id = ?", imageId, itemId).
		Find(&image)

	if image.IsZero() {
		return nil
	}

	return &image
}
