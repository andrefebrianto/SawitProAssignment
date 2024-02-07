package token

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/SawitProRecruitment/UserService/model"
)

var (
	privateKey = []byte(`-----BEGIN PRIVATE KEY-----
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
/5oAYoH6cSImgOGTkBS6KdvC1o0++w==
-----END PRIVATE KEY-----`)

	publicKey = []byte(`-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAwH1JLQeMUuanO86pRdId
SCoKQwJCVjiLzy2uuWPpAHLaMUqfcLxSXqFN3ffIZ54beMORPdUpg8Afl+roOPBl
daq8B/RxnDZYs1L7op5xvA4ZVe3jnkNZ6tbUB1o7Os0WcOVBTV9hSKzff9mWl9hk
iiMmST3RcipFAiOKaU6Q0LG6tpISgxaIksSrSLxs37Mbdzs57NHh0URPnM/jVlqc
E3YhhhxW9h9gyEAmZ/mfwTZAm9rDdf/phUQ6K9c4rlgwGZSTZXvGYvSzeGOffe9O
kb4FKfYgDGWuP0AHBDag4RflsNeiZ+1/M3WHRGlivBAOYsqQ8RUTPoyS4nStrd4N
3bs4+S+dzAEOpUNcH5MbF4ztHyYy8xlPt64w9/ySpre3MR1fNjECFWqDqy3UMv05
teux2b+uMpk7TQ1NQn0nEbHAZoaUggX7AW1+sVCxUDZJk5HzyapxH3/Xyr6bH7LH
uurXCervVqeVLHtlOQa8jEf03/WHQXhil07Ss1A+1VwZBny6KJW8aPC+uQv6wP3t
0VPaUqv3thFq+5LeTTCdFxTvsXwMjUaEbCR/l/EpEuOj2bvWQUv6u/zlbV/qNxtj
3qKUODx9uIxm7uGDb8MxbdeXQKoA++aWwYiWXr/Hf5SwPMjpa/ztUAsxo41c8Iv3
l//aV86EY4yKk19JId8H8qECAwEAAQ==
-----END PUBLIC KEY-----`)
)

func TestToken_Generate(t *testing.T) {
	type fields struct {
		PrivateKey []byte
		PublicKey  []byte
	}
	type args struct {
		ttl  time.Duration
		data model.Token
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "success to generate jwt token",
			fields: fields{
				PrivateKey: privateKey,
				PublicKey:  publicKey,
			},
			args: args{
				ttl: 2 * time.Hour,
			},
			want:    "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9",
			wantErr: false,
		},
		{
			name: "failed to parse private key",
			fields: fields{
				PrivateKey: []byte{},
				PublicKey:  []byte{},
			},
			args: args{
				ttl: 2 * time.Hour,
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := NewToken(tt.fields.PrivateKey, tt.fields.PublicKey)
			got, err := tr.Generate(tt.args.ttl, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Token.Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !strings.Contains(got, tt.want) { // just compare the header because the slight difference in the timestamp will generate different claims and signature
				t.Errorf("Token.Generate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToken_Validate(t *testing.T) {
	tokenData := model.Token{
		UserID: 168,
		Name:   "John Doe",
	}

	tknHandler := &Token{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}

	tknString, _ := tknHandler.Generate(2*time.Hour, tokenData)

	type fields struct {
		PrivateKey []byte
		PublicKey  []byte
	}
	type args struct {
		token string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    model.Token
		wantErr bool
	}{
		{
			name: "success to validate access token",
			fields: fields{
				PrivateKey: privateKey,
				PublicKey:  publicKey,
			},
			args: args{
				token: tknString,
			},
			want:    tokenData,
			wantErr: false,
		},
		{
			name: "failed to parse public key",
			fields: fields{
				PrivateKey: []byte{},
				PublicKey:  []byte{},
			},
			args: args{
				token: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhIjp7IlVzZXJJRCI6MSwiTmFtZSI6IkpvaG4gRG9lIn0sImV4cCI6MTcwNzIxNDA0OSwibmJmIjoxNzA3MjA2ODQ5LCJpYXQiOjE3MDcyMDY4NDl9.gVh7vMXRxEdYP_DvzkTtUZq7fbF5E37aozRgLRaKoZsHV5b1KGSGm0JIrsfOBph3s-Iz7AiI_VvapJEOh5MdSWJhQEFpCUpokz1ceSTyvg8ZAXq4YQG5aM5SvXmwJZNZDQdM26d2XwrXW2eGlA8Tt5rwc6FhNwtqI-vI0Q7m-8PNM4CBhR9_Wgd43jbonJPLqX93LWooqyCtdC-w2rTu1vBCEocYjw6eKJzLkWZzdIGW-xTqLPYMpniJpO4oik5fdhGa5XuPWBGP2FfURuhF8RzI7QbvpralBd0SkUXzG-mn1gK7HIGfa0cOFI3FzMGuzdz0EDwXJm_7LWe8IAAIvgKJqZYjQ74xfQpxvRdfF7-6yGHV9KqExVBHQCX1U4Gum5trqJqA5KEXLZt9qtAqA2-QZCwWLHMYSq_RekRCIMF1X34hGZaDa28s8LRgl10yB-evVJtpleG9KRps7lUUMxXUplw40UyTCs0F8BkfuarWs-kDIuS7Oow4c798GeIU7y1bQlgisEeZEpcwyzC5795ChbiWdBEqCc4xDLz0ZVU8QFhfihIEq19gXe-T3B7xeZuvtP4gudE4MyIPUaqP1nG8Ty5Dp_9JeoO_EbLpQxC3blQsjg5Y3GSOaEIFBLslRuoH0b1Rt9jBDdxZS_uFSrR-qqGAQMbaCl1O9-MwObI",
			},
			want:    model.Token{},
			wantErr: true,
		},
		{
			name: "validate expired access token",
			fields: fields{
				PrivateKey: privateKey,
				PublicKey:  publicKey,
			},
			args: args{
				token: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhIjp7IlVzZXJJRCI6MSwiTmFtZSI6IkpvaG4gRG9lIn0sImV4cCI6MTcwNzExNDA0OSwibmJmIjoxNzA3MjA2ODQ5LCJpYXQiOjE3MDcyMDY4NDl9.MpD3waPsGaVZmi7MxNNRhqfKA2eUvRWUTZpuP2Nn5tFEni4287XrpM_VCcBPVT4J9t38FozB2GBPY5K_ixOxvvczQ90y5ohw8PAPg6mRNbbdXZgWIszaTSdix1Ct9EWuNvHu7qWNrgBvsLELEkXlYoV9xr559ElrLFAgLBHHm1hZgvMc9VDl-8r2vx2cCeVQdiqUpQn7JWqJaUNXgmpdOFAxZ2KGDPTrGeBsRUnBsBXXoD-iP5fLB_TDGKprYCrX-JFgHUnwYIr9fljI_tI2zqxgKUalcoB1S1tZ12oA4YxEB528m04pVHq1LNQzdFsFIAw6lRpSlV0Xc7CsnvG6-6LEKfjjjFPnKnNwMHbqP4mKjIqBAUY7JTXyUIf72_VUQJ9pwV9t-67-qZrPwDgLo1qT61lTmeq0J5lzpFNXB1ILIEcWN1to66Pq-bOV_1jVuXsoCSJiGmf0oui6KnhZIrTppUxwLbCgT3RIN7awHDJw-eafqOD7Yal8cIA9uFh58-HC6561gRGmdnTSswNvOkHOks0q1YPfcT2Mha2RWnp561bd5hC-V3r5Fu0qOFWb_IAhpW68xPQWJThU-6Z_9Aa0hUuBTLOXnalY_MM27nF5smGtjqnsG2AAPRwDk7MRfEGID-1J-jy-2RzCMnJ0IL53r9EuTddw7kUQZIt_eVA",
			},
			want:    model.Token{},
			wantErr: true,
		},
		{
			name: "validate invalid access token that using different key pair",
			fields: fields{
				PrivateKey: privateKey,
				PublicKey:  publicKey,
			},
			args: args{
				token: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhIjp7IlVzZXJJRCI6IjZiNmFmYjE0LWRkZGUtNGI3MC1iZDVmLTJiMWJiZDZkYmI3MiIsIk5hbWUiOiJKb2huIERvZSJ9LCJleHAiOjE3MDcxMTQwNDksIm5iZiI6MTcwNzIwNjg0OSwiaWF0IjoxNzA3MjA2ODQ5fQ.bZG3LWjlOczenc4r2_zPkqQi-qEfejJMIHWEzZ-OG6ilMdXz8jMwH65o8fYqTI1OijfRIpTFR5txSxBMgEs2EC2HT6t6LHmgGKqwALnwIk0WkMCODrXfHZRLCXpCKwXefFO8NbIP76IuSdjB5hOgOfNFjxrBtAkUPzedxpvQmnIoo72-xEFG3T6OdaGhpArFKowQmuQoQesqfRwoR6_X2DKwVQeG6u_p5XBHjfLfwZyaFDKkbPF23WGNNsO-KI1jR-qrreCcG8kIA1rFBkBvHJ5_fCOoTDBxOyAtlgWpkWkLzultIy1N9_tBGW1Wm8V_x1yOAi_8hJSCSHUeSUliZH6PNBbkSI6tjBIwr1fOWS6HQ1WRL4dj3WMiwHcM6cX_VWChXaEgfFR-KfrMYkPVohJXlLiSVaDWzRTvpEgRrQOKi0r6gla3IJiPd7OQjmpIPeO7Yd2GSq3STpRx320k9MNzNLElOulskF6araxd8Tfl5qor6_2XjjOEfFyVm_c_NcxPJuhPsopLqfHSXoKONFBP7AgqFZzW70u35RsFQ2OT6KSIBl7J-DTivo22m5jvuX1ZWNhs32iqGFvJM5jiz85b1MT28S42YzINaneIErsT-2Mu33wEVjc-4ItkfJ1-mUSSQZYsut1JadROcWi47cbj1zjWApMPxaVrQPb0U_w",
			},
			want:    model.Token{},
			wantErr: true,
		},
		{
			name: "validate invalid access token that using different algorithm",
			fields: fields{
				PrivateKey: privateKey,
				PublicKey:  publicKey,
			},
			args: args{
				token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhIjp7IlVzZXJJRCI6IjZiNmFmYjE0LWRkZGUtNGI3MC1iZDVmLTJiMWJiZDZkYmI3MiIsIk5hbWUiOiJKb2huIERvZSJ9LCJleHAiOjE3MDcxMTQwNDksIm5iZiI6MTcwNzIwNjg0OSwiaWF0IjoxNzA3MjA2ODQ5fQ.Dz-MYfKogpoF4s2o_ZECGQTgFhAqHjNAijxnqsY_M7w",
			},
			want:    model.Token{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := NewToken(tt.fields.PrivateKey, tt.fields.PublicKey)
			got, err := tr.Validate(tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("Token.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Token.Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
