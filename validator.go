package main

import (
	"github.com/go-playground/validator"
)

//CustomValidator represents the custom input validator
type CustomValidator struct {
	Validator *validator.Validate
}

//Validate the input
func (cv *CustomValidator) Validate(input interface{}) error {
	//Validate agains the validator
	err := cv.Validator.Struct(input)

	if err != nil {
		return err
	}

	return nil
}
