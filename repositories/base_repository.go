package repositories

import (
	"gorm.io/gorm"
)

type Repository struct {
	connection *gorm.DB
}

func initRepository(conn *gorm.DB) *Repository {
	return &Repository{
		connection: conn,
	}
}
