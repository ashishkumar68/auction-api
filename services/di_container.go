//go:build wireinject
// +build wireinject

package services

import (
	"github.com/google/wire"
	"gorm.io/gorm"
)

func NewUserService(conn *gorm.DB) UserService {
	wire.Build(initUserService)

	return nil
}

func NewItemService(conn *gorm.DB) ItemService {
	wire.Build(initItemService)

	return nil
}