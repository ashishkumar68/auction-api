package repositories

import (
	"github.com/ashishkumar68/auction-api/models"
	"gorm.io/gorm"
	"log"
)

func (repo *Repository) FindItemCommentById(id uint) *models.ItemComment {
	var itemComment models.ItemComment
	repo.connection.
		Joins("UserCreated").
		Find(&itemComment, id)

	if itemComment.IsZero() {
		return nil
	}

	return &itemComment
}

func (repo *Repository) FindItemCommentByUuid(uuid string) *models.ItemComment {
	var itemComment models.ItemComment
	repo.connection.
		Joins("UserCreated").
		Find(&itemComment, "uuid = ?", uuid)

	if itemComment.IsZero() {
		return nil
	}

	return &itemComment
}

func (repo *Repository) SaveItemComment(comment *models.ItemComment) error {
	result := repo.connection.Create(comment)
	if result.Error != nil {
		log.Println("could not save new item comment due to error:", result.Error)
		return result.Error
	}

	return nil
}

func (repo *Repository) UpdateItemComment(comment *models.ItemComment) error {
	result := repo.connection.Save(comment)
	if result.Error != nil {
		log.Println("could not update item comment due to error:", result.Error)
		return result.Error
	}

	return nil
}

func (repo *Repository) DeleteItemComment(comment *models.ItemComment) error {
	result := repo.connection.Delete(comment)
	if result.Error != nil {
		log.Println("could not delete item comment due to error:", result.Error)
		return result.Error
	}

	return nil
}

func (repo *Repository) CountCommentsByItem(item *models.Item) uint {
	var count uint
	repo.connection.
		Table("item_comments").
		Select("count(*)").
		Where("item_comments.item_id = ? AND deleted_at IS NULL", item.ID).
		Scan(&count)

	return count
}

func (repo *Repository) FindCommentsByItem(item *models.Item) *gorm.DB {
	return repo.connection.
		Model(&models.ItemComment{}).
		Preload("UserCreated").
		Joins("Item").
		Where("Item.id = ?", item.ID).
		Order("item_comments.id DESC")
}
