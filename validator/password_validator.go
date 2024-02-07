package validator

import (
	reflect "reflect"

	"github.com/dlclark/regexp2"
	"github.com/go-playground/validator/v10"
)

var (
	validPasswordRegex = regexp2.MustCompile(`^(?=.*[A-Z])(?=.*[~`+"`"+`!@#$%^&*()_\-+={\[}\]|\:;\"'<,>.?\/])(?=.*[0-9]).{6,64}$`, 0)
)

func password(fl validator.FieldLevel) bool {
	field := fl.Field()

	switch field.Kind() {
	case reflect.String:
		text := field.String()
		isValid, _ := validPasswordRegex.MatchString(text)
		return isValid
	default:
		return false
	}
}
