package server

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/yurikilian/bills/pkg/exception"
	"io"
)

type Binder struct {
	validator *CustomValidator
}

func NewBinder() *Binder {
	return &Binder{
		validator: newCustomValidator(),
	}
}

type ValidationError struct {
	Path    string
	Message string
}

func (b *Binder) ReadBody(c *HttpContext, result interface{}) error {

	if err := readBody(c, result); err != nil {
		return exception.NewMalformedRequestProblem()
	}

	if vErr := b.validator.Validate(result); vErr != nil {
		vErrors := vErr.(validator.ValidationErrors)
		customErrors := make([]*exception.ValidationProblemDetail, 0)

		for _, vErr := range vErrors {
			customErrors = append(customErrors, exception.NewValidationProblemDetail(vErr.Tag(), vErr.Field(), vErr.Param()))
		}

		return exception.NewValidationProblem(customErrors)
	}

	return nil
}

func readBody(c *HttpContext, toBind interface{}) error {

	if c.Request().ContentLength == 0 {
		return nil
	}

	read, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(read, toBind)
	if err != nil {
		return err
	}

	return nil
}
