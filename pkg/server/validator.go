package server

import (
	"github.com/go-playground/validator/v10"
	"github.com/yurikilian/bills/pkg/exception"
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

func (v *CustomValidator) MapValidationProblems(vErr error) []*exception.ValidationProblemDetail {
	vErrors := vErr.(validator.ValidationErrors)
	customErrors := make([]*exception.ValidationProblemDetail, 0)

	for _, vErr := range vErrors {
		customErrors = append(customErrors, exception.NewValidationProblemDetail(vErr.Tag(), vErr.Field(), vErr.Param()))
	}
	return customErrors
}

func newCustomValidator() *CustomValidator {
	return &CustomValidator{
		validate: validator.New(),
	}
}

var Validator = newCustomValidator()
