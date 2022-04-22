package services

import (
	"context"
	"github.com/ashishkumar68/auction-api/forms"
	"github.com/ashishkumar68/auction-api/models"
	"github.com/ashishkumar68/auction-api/repositories"
)

type ItemCommentService interface {
	AddItemComment(ctx context.Context, form forms.AddItemCommentForm) (*models.ItemComment, error)
	UpdateItemComment(ctx context.Context, form forms.EditItemCommentForm) error
	DeleteItemComment(ctx context.Context, form forms.DeleteItemCommentForm) error
}

type ItemCommentServiceImplementor struct {
	repo *repositories.Repository
}

func (service *ItemCommentServiceImplementor) AddItemComment(
	ctx context.Context,
	form forms.AddItemCommentForm) (*models.ItemComment, error) {

	return nil, nil
}

func (service *ItemCommentServiceImplementor) UpdateItemComment(
	ctx context.Context,
	form forms.EditItemCommentForm) error {

	return nil
}

func (service *ItemCommentServiceImplementor) DeleteItemComment(
	ctx context.Context,
	form forms.DeleteItemCommentForm) error {

	return nil
}

func initItemCommentService(repo *repositories.Repository) ItemCommentService {
	return &ItemCommentServiceImplementor{repo: repo}
}