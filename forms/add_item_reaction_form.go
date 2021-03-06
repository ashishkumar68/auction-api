package forms

import "github.com/ashishkumar68/auction-api/models"

type AddItemReactionForm struct {
	AuditableForm

	Item         *models.Item `binding:"required"`
	ReactionType uint8        `json:"reactionType" binding:"oneof=0 1"`
}
