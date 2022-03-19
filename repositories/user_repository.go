package repositories

import "github.com/ashishkumar68/auction-api/models"

type UserRepository struct {
	BaseRepository
}

func (repo *UserRepository) Find(id int32) *models.User {
	var user *models.User
	repo.connection.Find(&user, id)

	return user
}