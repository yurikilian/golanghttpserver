package server

import (
	"encoding/json"
	"github.com/yurikilian/bills/pkg/exception"
	"io"
)

type Binder struct {
	validator *CustomValidator
}

func NewBinder() *Binder {
	return &Binder{
		validator: Validator,
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
		customErrors := b.validator.MapValidationProblems(vErr)

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
