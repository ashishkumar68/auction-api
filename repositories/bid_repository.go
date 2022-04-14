package repositories

import (
	"fmt"
	"github.com/ashishkumar68/auction-api/models"
	"log"
)

func (repo *Repository) FindBidById(id uint) *models.Bid {
	var bid models.Bid
	repo.connection.Model(&models.Bid{}).Find(&bid, id)
	if bid.IsZero() {
		return nil
	}

	return &bid
}

func (repo *Repository) SaveBid(bid *models.Bid) error {
	result := repo.connection.Model(&models.Bid{}).Create(bid)
	if result.Error != nil {
		log.Println(fmt.Sprintf("Could not insert new record for type: %T", bid))
		log.Println(fmt.Sprintf("Insert error: %s", result.Error))
		return result.Error
	}

	return nil
}

func (repo *Repository) UpdateBid(bid *models.Bid) error {
	result := repo.connection.Model(&models.Bid{}).Save(bid)
	if result.Error != nil {
		log.Println(fmt.Sprintf("Could not update record for type: %T", bid))
		log.Println(fmt.Sprintf("Update error: %s", result.Error))
		return result.Error
	}

	return nil
}

func (repo *Repository) FindBidByUuid(uuid string) *models.Bid {
	var bid models.Bid
	repo.connection.Model(&models.Bid{}).Preload("UserCreated").Where("uuid = ?", uuid).First(&bid)
	if bid.IsZero() {
		return nil
	}

	return &bid
}

func (repo *Repository) FindBidByItem(item *models.Item, user *models.User) *models.Bid {
	var bid models.Bid

	repo.connection.
		Model(&models.Bid{}).
		//Table("bids").
		Select("bids.*").
		Joins("JOIN items ON bids.item_id = items.id").
		Joins("JOIN users ON bids.created_by = users.id").
		Where("users.id = ? AND items.id = ?", user.ID, item.ID).
		Scan(&bid)
	if bid.IsZero() {
		return nil
	}

	return &bid
}
