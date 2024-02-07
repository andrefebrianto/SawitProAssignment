package validator

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

type dummyStruct struct {
	Password string `validate:"password"`
}

type invalidStruct struct {
	Age int `validate:"password"`
}

func TestValidator_Password(t *testing.T) {
	type fields struct {
		val *validator.Validate
	}
	type args struct {
		dataStruct dummyStruct
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "password with length less than 6",
			fields: fields{
				val: validator.New(),
			},
			args: args{
				dataStruct: dummyStruct{
					Password: "p$1S",
				},
			},
			wantErr: true,
		},
		{
			name: "password with length more than 64",
			fields: fields{
				val: validator.New(),
			},
			args: args{
				dataStruct: dummyStruct{
					Password: "-,Zk2Xo4)Ko9~#cx~%^V=5Ef_sci)_,qFUQVGt>rqGK1+Tg]gay,kmy,qThFm~PKG%%@ui",
				},
			},
			wantErr: true,
		},
		{
			name: "password without symbol",
			fields: fields{
				val: validator.New(),
			},
			args: args{
				dataStruct: dummyStruct{
					Password: "Qmqpwz5xqWHE4rMa",
				},
			},
			wantErr: true,
		},
		{
			name: "password without number",
			fields: fields{
				val: validator.New(),
			},
			args: args{
				dataStruct: dummyStruct{
					Password: "TmCgRPWwyBU%E:TJ",
				},
			},
			wantErr: true,
		},
		{
			name: "password without capital",
			fields: fields{
				val: validator.New(),
			},
			args: args{
				dataStruct: dummyStruct{
					Password: ",fwy,b^t3gmtgfil",
				},
			},
			wantErr: true,
		},
		{
			name: "password with emojis",
			fields: fields{
				val: validator.New(),
			},
			args: args{
				dataStruct: dummyStruct{
					Password: "👍👍👍👍👍",
				},
			},
			wantErr: true,
		},
		{
			name: "password with an emoji, alphanumerics, and symbols",
			fields: fields{
				val: validator.New(),
			},
			args: args{
				dataStruct: dummyStruct{
					Password: "JsDG^2DYPz5Gdqys👍",
				},
			},
			wantErr: false,
		},
		{
			name: "correct password with length between 6 and 64",
			fields: fields{
				val: validator.New(),
			},
			args: args{
				dataStruct: dummyStruct{
					Password: "^1@Z@D)115cza6Ly",
				},
			},
			wantErr: false,
		},
		{
			name: "correct password with 6 characters",
			fields: fields{
				val: validator.New(),
			},
			args: args{
				dataStruct: dummyStruct{
					Password: "A_%cJ8",
				},
			},
			wantErr: false,
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

	t.Run("invalid data type", func(t *testing.T) {
		v := NewValidator(validator.New())

		input := invalidStruct{
			Age: 21,
		}

		if err := v.Validate(input); (err != nil) != true {
			t.Errorf("Validator.Validate() error = %v, wantErr %v", err, true)
		}
	})
}
