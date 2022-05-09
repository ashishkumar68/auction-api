package repositories

import (
	"database/sql"
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

// FindReactionsCountByItems returns item reactions mapped by item id.
func (repo *Repository) FindReactionsCountByItems(items []*models.Item) map[uint][]models.ItemReactionTypeCount {
	itemIds := make([]uint, 0)
	for _, item := range items {
		itemIds = append(itemIds, item.ID)
	}

	var itemReactionsCount []models.ItemReactionTypeCount
	repo.connection.Raw(`
SELECT i.id as item_id, r.type as reaction_type, count(r.id) as reaction_count
FROM items i
JOIN reactions r ON r.item_id = i.id
WHERE i.id IN @ids
GROUP BY i.id, r.type
ORDER BY i.id ASC
;
`, sql.Named("ids", itemIds)).Scan(&itemReactionsCount)

	itemReactionsMap := make(map[uint][]models.ItemReactionTypeCount)
	for _, itemReactionCount := range itemReactionsCount {
		itemReactionCount.ReactionTypeText = models.FindReactionTypeString(itemReactionCount.ReactionType)
		if itemReactions, ok := itemReactionsMap[itemReactionCount.ItemId]; ok {
			itemReactionsMap[itemReactionCount.ItemId] = append(itemReactions, itemReactionCount)
		} else {
			itemReactionsMap[itemReactionCount.ItemId] = []models.ItemReactionTypeCount{itemReactionCount}
		}
	}

	return itemReactionsMap
}
