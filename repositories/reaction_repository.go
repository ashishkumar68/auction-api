package repositories

import (
	"github.com/ashishkumar68/auction-api/models"
	"log"
)

func (repo *Repository) FindReactionByItemAndUser(item *models.Item, user *models.User) *models.Reaction {
	var reaction models.Reaction
	repo.connection.
		Joins("JOIN items ON items.id = reactions.item_id").
		Joins("JOIN users ON users.id = reactions.created_by").
		Where("items.id = ? AND users.id = ?", item.ID, user.ID).
		Find(&reaction)

	if reaction.IsZero() {
		return nil
	}

	return &reaction
}

func (repo *Repository) SaveReaction(reaction *models.Reaction) error {
	result := repo.connection.Create(reaction)
	if result.Error != nil {
		log.Println("could not save reaction due to error:", result.Error)
		return result.Error
	}

	return nil
}

func (repo *Repository) UpdateReaction(reaction *models.Reaction) error {
	result := repo.connection.Updates(reaction)
	if result.Error != nil {
		log.Println("could not update reaction due to error:", result.Error)
		return result.Error
	}

	return nil
}

func (repo *Repository) DeleteReaction(reaction *models.Reaction) error {
	result := repo.connection.Delete(reaction)
	if result.Error != nil {
		log.Println("could not delete reaction due to error:", result.Error)
		return result.Error
	}

	return nil
}
