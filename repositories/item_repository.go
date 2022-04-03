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

func (repo *ItemRepository) Find(id int32) *models.Item {
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
	repo.connection.Where("uuid = ?", uuid).First(&item)

	return &item
}