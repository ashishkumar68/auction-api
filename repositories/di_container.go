//+build wireinject

package repositories

import (
	"github.com/google/wire"
	"gorm.io/gorm"
)

func NewUserRepository(conn *gorm.DB) *UserRepository {
	wire.Build(InitUserRepository)

	return nil
}

func NewItemRepository(conn *gorm.DB) *ItemRepository {
	wire.Build(InitItemRepository)

	return nil
}