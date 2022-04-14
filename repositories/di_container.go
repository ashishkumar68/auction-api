//go:build wireinject
// +build wireinject

package repositories

import (
	"github.com/google/wire"
	"gorm.io/gorm"
)

func NewRepository(conn *gorm.DB) *Repository {
	wire.Build(initRepository)

	return nil
}
