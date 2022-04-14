package repositories

import (
	"fmt"
	"github.com/ashishkumar68/auction-api/models"
	"log"
	"strings"
)

func (repo *Repository) FindUserById(id uint) *models.User {
	var user models.User
	repo.connection.Find(&user, id)
	if user.IsZero() {
		return nil
	}

	return &user
}

func (repo *Repository) SaveUser(user *models.User) error {

	result := repo.connection.Create(user)
	if result.Error != nil {
		log.Println(fmt.Sprintf("Could not insert new record for type: %T", user))
		log.Println(fmt.Sprintf("Insert error: %s", result.Error))
		if !strings.Contains(result.Error.Error(), "Deadlock") {
			return result.Error
		}
		// try 4 more times when facing deadlock.
		for retry := 1; retry <= 4; retry++ {
			result = repo.connection.Create(user)
			if result.Error != nil && strings.Contains(result.Error.Error(), "Deadlock") {
				log.Println("retry attempt saving user", retry)
				continue
			} else {
				break
			}
		}
	}

	return nil
}

func (repo *Repository) UpdateUser(user *models.User) error {
	result := repo.connection.Save(user)
	if result.Error != nil {
		log.Println(fmt.Sprintf("Could not update record for type: %T", user))
		log.Println(fmt.Sprintf("Update error: %s", result.Error))
		return result.Error
	}

	return nil
}

func (repo *Repository) FindUserByEmail(email string) *models.User {
	var user models.User
	repo.connection.Where("email = ? AND deleted_at IS NULL", email).First(&user)
	if user.IsZero() {
		return nil
	}

	return &user
}
