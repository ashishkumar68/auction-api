package repositories

import (
	"fmt"
	"github.com/ashishkumar68/auction-api/models"
	"gorm.io/gorm"
	"log"
)

type ItemRepository struct {
	BaseRepository
}

func InitItemRepository(conn *gorm.DB) *ItemRepository {
	return &ItemRepository{
		BaseRepository: BaseRepository{connection: conn},
	}
}

func (repo *ItemRepository) Find(id uint) *models.Item {
	var item models.Item
	repo.connection.Find(&item, id)

	return &item
}

func (repo *ItemRepository) Save(item *models.Item) error {
	result := repo.connection.Create(item)
	if result.Error != nil {
		log.Println(fmt.Sprintf("Could not insert new record for type: %T", item))
		log.Println(fmt.Sprintf("Insert error: %s", result.Error))
		return result.Error
	}

	return nil
}

func (repo *ItemRepository) Update(item *models.Item) error {
	result := repo.connection.Save(item)
	if result.Error != nil {
		log.Println(fmt.Sprintf("Could not update record for type: %T", item))
		log.Println(fmt.Sprintf("Update error: %s", result.Error))
		return result.Error
	}

	return nil
}

func (repo *ItemRepository) FindByUuid(uuid string) *models.Item {
	var item models.Item
	repo.connection.Preload("UserCreated").Where("uuid = ?", uuid).First(&item)

	return &item
}

func (repo *ItemRepository) FindByName(name string) []models.Item {
	var items []models.Item
	repo.connection.Where("name LIKE ?", "%"+name+"%").Find(&items)

	return items
}

func (repo *ItemRepository) List() *gorm.DB {
	return repo.connection.
		Joins("UserCreated").
		Joins("UserUpdated").
		Model(&models.Item{})
}
