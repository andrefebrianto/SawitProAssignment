package validator

type ValidatorInterface interface {
	Validate(dataStruct interface{}) error
}
