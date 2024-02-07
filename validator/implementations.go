package validator

func (v *Validator) Validate(dataStruct interface{}) error {
	return v.val.Struct(dataStruct)
}
