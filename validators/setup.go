package validators

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"log"
)

func SetupCustomValidators() {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		log.Println("could not tap into gin validator to register custom validations.")
		return
	}

	for tag, validatorFunc := range customValidators {
		err := v.RegisterValidation(tag, validatorFunc)
		if err != nil {
			log.Println(fmt.Sprintf("could not initialise validator: %s due to err: %s", tag, err.Error()))
		}
	}
}
