package repositories

import (
	"fmt"
	"github.com/ashishkumar68/auction-api/models"
	"gorm.io/gorm"
	"log"
)

type UserRepository struct {
	BaseRepository
}

func InitUserRepository(conn *gorm.DB) *UserRepository {
	return &UserRepository{
		BaseRepository: BaseRepository{connection: conn},
	}
}

func (repo *UserRepository) Find(id int32) *models.User {
	var user *models.User
	repo.connection.Find(&user, id)

	return user
}

func (repo *UserRepository) Save(user *models.User) error {
	result := repo.connection.Create(user)
	if result.Error != nil {
		log.Println(fmt.Sprintf("Could not insert new record for type: %T", user))
		log.Println(fmt.Sprintf("Insert error: %s", result.Error))
		return result.Error
	}

	return nil
}

func (repo *UserRepository) Update(user *models.User) error {
	result := repo.connection.Save(user)
	if result.Error != nil {
		log.Println(fmt.Sprintf("Could not update record for type: %T", user))
		log.Println(fmt.Sprintf("Update error: %s", result.Error))
		return result.Error
	}

	return nil
}

func (repo *UserRepository) FindByEmail(email string) *models.User {
	var user models.User
	repo.connection.Where("email = ? AND deleted_at IS NULL", email).First(&user)

	return &user
}