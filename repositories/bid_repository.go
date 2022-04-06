package repositories

import (
	"fmt"
	"github.com/ashishkumar68/auction-api/models"
	"gorm.io/gorm"
	"log"
)

type BidRepository struct {
	BaseRepository
}

func InitBidRepository(conn *gorm.DB) *BidRepository {
	return &BidRepository{
		BaseRepository: BaseRepository{connection: conn},
	}
}

func (repo *BidRepository) Find(id uint) *models.Bid {
	var bid models.Bid
	repo.connection.Find(&bid, id)

	return &bid
}

func (repo *BidRepository) Save(bid *models.Bid) error {
	result := repo.connection.Create(bid)
	if result.Error != nil {
		log.Println(fmt.Sprintf("Could not insert new record for type: %T", bid))
		log.Println(fmt.Sprintf("Insert error: %s", result.Error))
		return result.Error
	}

	return nil
}

func (repo *BidRepository) Update(bid *models.Bid) error {
	result := repo.connection.Save(bid)
	if result.Error != nil {
		log.Println(fmt.Sprintf("Could not update record for type: %T", bid))
		log.Println(fmt.Sprintf("Update error: %s", result.Error))
		return result.Error
	}

	return nil
}

func (repo *BidRepository) FindByUuid(uuid string) *models.Bid {
	var bid models.Bid
	repo.connection.Preload("UserCreated").Where("uuid = ?", uuid).First(&bid)

	return &bid
}

func (repo *BidRepository) FindByItem(item *models.Item, user *models.User) *models.Bid {
	var bid models.Bid
	repo.connection.
		Joins("Item").
		Joins("UserCreated").
		Where("UserCreated.ID = ? AND Item.ID = ?", user.ID, item.ID).
		First(&bid)

	return &bid
}