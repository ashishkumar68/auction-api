package services

import (
	"context"
	"fmt"
	"github.com/ashishkumar68/auction-api/forms"
	"github.com/ashishkumar68/auction-api/models"
	"github.com/ashishkumar68/auction-api/repositories"
	"gorm.io/gorm"
	"log"
)

type UserService interface {
	NewRegister(ctx context.Context, form forms.RegisterNewUserForm) (*models.User, error)
	Login(ctx context.Context, form forms.LoginUserForm) (*models.LoggedInUser, error)
}

type UserServiceImplementor struct {
	repository *repositories.Repository
}

func initUserService(db *gorm.DB) UserService {
	return &UserServiceImplementor{
		repository: repositories.NewRepository(db),
	}
}

func (service *UserServiceImplementor) NewRegister(
	ctx context.Context,
	form forms.RegisterNewUserForm) (*models.User, error) {

	newUser := models.NewUserFromValues(form.FirstName, form.LastName, form.Email, form.Password)
	hashedPass, err := HashPassword(newUser.PlainPassword)
	if err != nil {
		log.Println(fmt.Sprintf("Could not hashed password while saving user information."))
		log.Println("err:", err)
		return nil, err
	}
	newUser.Password = hashedPass
	err = service.repository.SaveUser(newUser)
	if err != nil {
		log.Println(fmt.Sprintf("Could not save user information."))
		log.Println("err:", err)
		return nil, err
	}

	return newUser, nil
}

func (service *UserServiceImplementor) Login(
	ctx context.Context,
	form forms.LoginUserForm) (*models.LoggedInUser, error) {

	existingUser := service.repository.FindUserByEmail(form.Email)
	if nil == existingUser {
		return nil, UserEmailDoesntExist
	}
	if !CompareHashAndPass(existingUser.Password, form.Password) {
		return nil, PasswordsDontMatch
	}
	user := models.CreateLoggedInUserByUser(*existingUser)
	accessToken, err := GenerateNewJwtToken(*existingUser, TokenTypeAccess)
	if err != nil {
		log.Println("could not create new JWT access token")
		log.Println(err)
		return nil, err
	}
	refreshToken, err := GenerateNewJwtToken(*existingUser, TokenTypeRefresh)
	if err != nil {
		log.Println("could not create new JWT refresh token")
		log.Println(err)
		return nil, err
	}
	user.AccessToken = accessToken
	user.RefreshToken = refreshToken

	return user, nil
}
