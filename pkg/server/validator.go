package server

import (
	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	validate *validator.Validate
}

func (v *CustomValidator) Validate(i interface{}) error {
	if err := v.validate.Struct(i); err != nil {
		return err
	}

	return nil
}

func newCustomValidator() *CustomValidator {
	return &CustomValidator{
		validate: validator.New(),
	}
}
