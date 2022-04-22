//go:build wireinject
// +build wireinject

package services

import (
	"github.com/ashishkumar68/auction-api/repositories"
	"github.com/google/wire"
	"gorm.io/gorm"
)

func NewUserService(conn *gorm.DB) UserService {
	wire.Build(repositories.NewRepository, initUserService)

	return nil
}

func NewItemService(conn *gorm.DB) ItemService {
	wire.Build(repositories.NewRepository, initItemService)

	return nil
}

func NewReactionService(conn *gorm.DB) ReactionService {
	wire.Build(repositories.NewRepository, initReactionService)

	return nil
}

func NewItemCommentService(conn *gorm.DB) ItemCommentService {
	wire.Build(repositories.NewRepository, initItemCommentService)

	return nil
}