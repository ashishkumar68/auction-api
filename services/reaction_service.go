package services

import (
	"context"
	"fmt"
	"github.com/ashishkumar68/auction-api/forms"
	"github.com/ashishkumar68/auction-api/models"
	"github.com/ashishkumar68/auction-api/repositories"
	"gorm.io/gorm"
	"log"
)

type ReactionService interface {
	AddReactionToItem(ctx context.Context, form forms.AddItemReactionForm) (*models.Reaction, error)
	RemoveReactionFromItem(ctx context.Context, form forms.RemoveItemReactionForm) error
}

type ReactionServiceImplementor struct {
	repository *repositories.Repository
}

func (service *ReactionServiceImplementor) AddReactionToItem(
	_ context.Context,
	form forms.AddItemReactionForm) (*models.Reaction, error) {

	existingReaction := service.repository.FindReactionByItemAndUser(form.Item, form.ActionUser)
	var newReaction *models.Reaction
	if existingReaction == nil {
		newReaction = models.NewReactionFromValues(form.Item, form.ReactionType, form.ActionUser)
		if err := service.repository.SaveReaction(newReaction); err != nil {
			log.Println("could not add reaction to item due to err:", err)
			return nil, err
		}
	} else {
		if existingReaction.Type == form.ReactionType {
			return existingReaction, nil
		}
		existingReaction.Type = form.ReactionType
		newReaction = existingReaction
		if err := service.repository.UpdateReaction(existingReaction); err != nil {
			log.Println("could not add reaction to item due to err:", err)
			return nil, err
		}
	}

	return newReaction, nil
}

func (service *ReactionServiceImplementor) RemoveReactionFromItem(
	_ context.Context,
	form forms.RemoveItemReactionForm) error {

	existingReaction := service.repository.FindReactionByItemAndUser(form.Item, form.ActionUser)
	if existingReaction != nil {
		if err := service.repository.DeleteReaction(existingReaction); err != nil {
			log.Println(fmt.Sprintf(
				"could not delete user:%s reaction for item: %d due to err", form.ActionUser.Email, form.Item.ID,
			), err)
			return err
		}
	}

	return nil
}

func InitReactionService(db *gorm.DB) ReactionService {
	return &ReactionServiceImplementor{repository: repositories.NewRepository(db)}
}
