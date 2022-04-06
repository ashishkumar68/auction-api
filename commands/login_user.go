package commands

import (
	"context"
	"github.com/ashishkumar68/auction-api/models"
	"github.com/ashishkumar68/auction-api/repositories"
	"github.com/ashishkumar68/auction-api/services"
	"github.com/gogolfing/cbus"
	"gorm.io/gorm"
	"log"
)

type LoginUserCommand struct {
	DB *gorm.DB

	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (cmd *LoginUserCommand) Type() string {
	return "LoginUser"
}

func LoginUserHandler(ctx context.Context, command cbus.Command) (interface{}, error) {
	var user models.LoggedInUser
	cmd := command.(*LoginUserCommand)
	existingUser := repositories.NewUserRepository(cmd.DB).FindByEmail(cmd.Email)
	if existingUser.IsZero() {
		return nil, services.UserEmailDoesntExist
	}
	if !services.CompareHashAndPass(existingUser.Password, cmd.Password) {
		return nil, services.PasswordsDontMatch
	}
	user = models.CreateLoggedInUserByUser(*existingUser)
	accessToken, err := services.GenerateNewJwtToken(*existingUser, services.TokenTypeAccess)
	if err != nil {
		log.Println("could not create new JWT access token")
		log.Println(err)
		return nil, err
	}
	refreshToken, err := services.GenerateNewJwtToken(*existingUser, services.TokenTypeRefresh)
	if err != nil {
		log.Println("could not create new JWT refresh token")
		log.Println(err)
		return nil, err
	}
	user.AccessToken = accessToken
	user.RefreshToken = refreshToken

	return user, nil
}
