package repositories

import (
	"gorm.io/gorm"
)

type BaseRepository struct {
	connection *gorm.DB
}

