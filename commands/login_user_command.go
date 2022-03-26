package commands

import (
	"context"
	"github.com/ashishkumar68/auction-api/models"
	"github.com/gogolfing/cbus"
)

type LoginUserCommand struct {
	Email		string	`json:"email" binding:"required"`
	Password	string	`json:"password" binding:"required"`
}

func (cmd *LoginUserCommand) Type() string {
	return "LoginUser"
}

func LoginUserHandler(ctx context.Context, command cbus.Command) (interface{}, error) {
	var user models.LoggedInUser
	_ = command.(*LoginUserCommand)


	return user, nil
}