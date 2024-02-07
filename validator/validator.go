package validator

import "github.com/go-playground/validator/v10"

type Validator struct {
	val *validator.Validate
}

func NewValidator(val *validator.Validate) *Validator {
	return &Validator{
		val: val,
	}
}
