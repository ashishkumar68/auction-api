//+build wireinject

package repositories

import (
	"github.com/ashishkumar68/auction-api/database"
	"github.com/google/wire"
	"gorm.io/gorm"
)

func NewUserRepository(config *gorm.Config) *UserRepository {
	wire.Build(database.NewConnection, InitUserRepository)

	return nil
}