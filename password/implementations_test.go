package password

import "testing"

func TestPassword_GenerateFromPassword(t *testing.T) {
	type fields struct {
		hashCost int
	}
	type args struct {
		password string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "hash password fewer than 72 bytes",
			fields: fields{
				hashCost: 12,
			},
			args: args{
				password: "(F2r+CYG~'AZb'sPwcV#;PHkaS(tSZC4ukEV^rd7KNGN-M2l8[]&.Z8Imm%C404Pa9_Ld8Twhrex1K;HY{85}J(HS.,6+oA1a+m7",
			},
			wantErr: true,
		},
		{
			name: "hash password fewer than 72 bytes",
			fields: fields{
				hashCost: 12,
			},
			args: args{
				password: "veryverysecurepassword",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Password{
				hashCost: tt.fields.hashCost,
			}
			_, err := p.GenerateFromPassword(tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Password.GenerateFromPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestPassword_VerifyPassword(t *testing.T) {
	type fields struct {
		hashCost int
	}
	type args struct {
		hashedPassword string
		password       string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "verify correct password",
			fields: fields{
				hashCost: 12,
			},
			args: args{
				hashedPassword: "$2a$12$hWASkUwEkcS1CbsyRRwoBew5r7qwmXwH4YJyP.S149hghOg77UEQW",
				password:       "veryverysecurepassword",
			},
			wantErr: false,
		},
		{
			name: "verify incorrect password",
			fields: fields{
				hashCost: 12,
			},
			args: args{
				hashedPassword: "$2a$12$hWASkUwEkcS1CbsyRRwoBew5r7qwmXwH4YJyP.S149hghOg77UEQW",
				password:       "verysafepassword",
			},
			wantErr: true,
		},
		{
			name: "verify incorrect hashed password",
			fields: fields{
				hashCost: 12,
			},
			args: args{
				hashedPassword: "$2a$12$",
				password:       "veryverysecurepassword",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Password{
				hashCost: tt.fields.hashCost,
			}
			if err := p.VerifyPassword(tt.args.hashedPassword, tt.args.password); (err != nil) != tt.wantErr {
				t.Errorf("Password.VerifyPassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
