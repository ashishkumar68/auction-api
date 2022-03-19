package commands

import (
	"context"
	"github.com/ashishkumar68/auction-api/repositories"
)

type RegisterNewUserCommand struct {
	FirstName	string	`json:"firstName" binding:"required,min=3,max=200"`
	LastName	string	`json:"lastName" binding:"required,min=3,max=200"`
	Email		string	`json:"email" binding:"required,min=3,max=200,email"`
	Password	string	`json:"password" binding:"required,min=8,max=80"`
}

type RegisterNewUserHandler struct {
	userRepository *repositories.UserRepository
}

func (handler *RegisterNewUserHandler) Handle(ctx context.Context, command *RegisterNewUserCommand) error {

	return nil
}