package commands

import (
	"context"
	"fmt"
	"github.com/ashishkumar68/auction-api/models"
	"github.com/ashishkumar68/auction-api/repositories"
	"github.com/ashishkumar68/auction-api/services"
	"github.com/gogolfing/cbus"
	"gorm.io/gorm"
	"log"
)

type RegisterNewUserCommand struct {
	DB *gorm.DB

	FirstName	string	`json:"firstName" binding:"required,min=3,max=200"`
	LastName	string	`json:"lastName" binding:"required,min=3,max=200"`
	Email		string	`json:"email" binding:"required,min=3,max=200,email"`
	Password	string	`json:"password" binding:"required,min=8,max=80"`
}

func (cmd *RegisterNewUserCommand) Type() string {
	return "RegisterNewUser"
}

func RegisterNewUserHandler(ctx context.Context, command cbus.Command) (interface{}, error) {
	registerUserCmd := command.(*RegisterNewUserCommand)
	newUser := models.NewUserFromValues(
		registerUserCmd.FirstName,
		registerUserCmd.LastName,
		registerUserCmd.Email,
		registerUserCmd.Password)
	hashedPass, err := services.HashPassword(newUser.PlainPassword)
	if err != nil {
		log.Println(fmt.Sprintf("Could not hashed password while saving user information."))
		log.Println("err:", err)
		return nil, err
	}
	newUser.Password = hashedPass
	err = repositories.NewUserRepository(registerUserCmd.DB).Save(&newUser)
	if err != nil {
		log.Println(fmt.Sprintf("Could not save user information."))
		log.Println("err:", err)
		return nil, err
	}

	return newUser, nil
}