package services

import (
	"context"
	"fmt"
	"github.com/ashishkumar68/auction-api/forms"
	"github.com/ashishkumar68/auction-api/models"
	"github.com/ashishkumar68/auction-api/repositories"
	"log"
)

var (
	ItemCommentEditFailed      = fmt.Errorf("failed to update item comment")
	DifferentAuthorItemComment = fmt.Errorf("found different item comment author")
	ItemCommentDeleteFailed    = fmt.Errorf("failed to delete item comment")
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
	_ context.Context,
	form forms.AddItemCommentForm) (*models.ItemComment, error) {

	newItemComment := models.NewItemComment(form.Comment, form.Item, form.ActionUser)
	err := service.repo.SaveItemComment(newItemComment)
	if err != nil {
		log.Println(fmt.Sprintf(
			"could not add comment: %s on item: %d due to error: %s", form.Comment, form.Item.ID, err,
		))
		return nil, err
	}

	return newItemComment, nil
}

func (service *ItemCommentServiceImplementor) UpdateItemComment(
	_ context.Context,
	form forms.EditItemCommentForm) error {

	comment := service.repo.FindItemCommentById(form.CommentId)
	if comment == nil {
		log.Println(fmt.Sprintf("comment: %d was not found while editing", form.CommentId))
		return models.ItemCommentNotFound
	}
	if !comment.IsAuthor(*form.ActionUser) {
		log.Println(fmt.Sprintf(
			"failed to update item comment: %d due to err: %s", comment.ID, DifferentAuthorItemComment,
		))
		return DifferentAuthorItemComment
	}
	if comment.Description != form.EditedComment {
		comment.Description = form.EditedComment
		err := service.repo.UpdateItemComment(comment)
		if err != nil {
			log.Println(fmt.Sprintf("failed to update item comment: %d due to err: %s", comment.ID, err))
			return ItemCommentEditFailed
		}
	}

	return nil
}

func (service *ItemCommentServiceImplementor) DeleteItemComment(
	_ context.Context,
	form forms.DeleteItemCommentForm) error {

	comment := service.repo.FindItemCommentById(form.CommentId)
	if comment == nil {
		log.Println(fmt.Sprintf("comment: %d was not found while editing", form.CommentId))
		return models.ItemCommentNotFound
	}
	if !comment.IsAuthor(*form.ActionUser) {
		log.Println(fmt.Sprintf(
			"failed to delete item comment: %d due to err: %s", comment.ID, DifferentAuthorItemComment,
		))
		return DifferentAuthorItemComment
	}
	err := service.repo.DeleteItemComment(comment)
	if err != nil {
		log.Println(fmt.Sprintf("could not delete item comment: %d due to error: %s", comment.ID, err))
		return ItemCommentDeleteFailed
	}

	return nil
}

func initItemCommentService(repo *repositories.Repository) ItemCommentService {
	return &ItemCommentServiceImplementor{repo: repo}
}
