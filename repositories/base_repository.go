package repositories

import (
	"fmt"
	"github.com/ashishkumar68/auction-api/models"
	"gorm.io/gorm"
	"log"
)

type BaseRepository struct {
	connection *gorm.DB
}

func (repo *BaseRepository) Save(identity *models.Identity) error {
	result := repo.connection.Create(identity)
	if result.Error != nil {
		log.Println(fmt.Sprintf("Could not insert new record for type: %T", identity))
		log.Println(fmt.Sprintf("Insert error: %s", result.Error))
		return result.Error
	}

	return nil
}

func (repo *BaseRepository) Update(updatedIdentity *models.Identity) error {
	result := repo.connection.Save(updatedIdentity)
	if result.Error != nil {
		log.Println(fmt.Sprintf("Could not update record for type: %T", updatedIdentity))
		log.Println(fmt.Sprintf("Update error: %s", result.Error))
		return result.Error
	}

	return nil
}
