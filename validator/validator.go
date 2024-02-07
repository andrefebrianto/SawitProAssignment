package validator

import "github.com/go-playground/validator/v10"

type Validator struct {
	val *validator.Validate
}

func NewValidator(val *validator.Validate) *Validator {
	// register custom password validation
	_ = val.RegisterValidation("password", password)

	return &Validator{
		val: val,
	}
}
