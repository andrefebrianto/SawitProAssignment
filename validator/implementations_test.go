package validator

import (
	"testing"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/go-playground/validator/v10"
)

func TestValidator_Validate(t *testing.T) {
	type fields struct {
		val *validator.Validate
	}
	type args struct {
		dataStruct interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "validate correct data",
			fields: fields{
				val: validator.New(),
			},
			args: args{
				dataStruct: generated.LoginUserJSONRequestBody{
					Phone:    "+628123456789",
					Password: "secure_password",
				},
			},
			wantErr: false,
		},
		{
			name: "validate incorrect data",
			fields: fields{
				val: validator.New(),
			},
			args: args{
				dataStruct: generated.LoginUserJSONRequestBody{
					Phone:    "08123456789",
					Password: "secure_password",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewValidator(tt.fields.val)
			if err := v.Validate(tt.args.dataStruct); (err != nil) != tt.wantErr {
				t.Errorf("Validator.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
