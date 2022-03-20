//+build wireinject

package repositories

import (
	"github.com/ashishkumar68/auction-api/database"
	"github.com/google/wire"
)

func NewUserRepository() *UserRepository {
	wire.Build(database.NewConnection, InitUserRepository)

	return nil
}