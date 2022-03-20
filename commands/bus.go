package commands

import (
	"github.com/gogolfing/cbus"
)

func NewCommandBus() *cbus.Bus {
	bus := &cbus.Bus{}
	bus.Handle(&RegisterNewUserCommand{}, cbus.HandlerFunc(RegisterNewUserHandler))

	return bus
}