//go:build wireinject
// +build wireinject

package repositories

import (
	"github.com/google/wire"
	"gorm.io/gorm"
)

func NewUserRepository(conn *gorm.DB) *UserRepository {
	wire.Build(initUserRepository)

	return nil
}

func NewItemRepository(conn *gorm.DB) *ItemRepository {
	wire.Build(initItemRepository)

	return nil
}

func NewBidRepository(conn *gorm.DB) *BidRepository {
	wire.Build(initBidRepository)

	return nil
}