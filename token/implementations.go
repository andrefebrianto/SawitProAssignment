package token

import (
	"fmt"
	"time"

	"github.com/SawitProRecruitment/UserService/model"
	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	Data model.Token `json:"data"`
	jwt.RegisteredClaims
}

func (t *Token) Generate(ttl time.Duration, data model.Token) (string, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM(t.PrivateKey)
	if err != nil {
		return "", fmt.Errorf("create: parse key: %w", err)
	}

	now := time.Now().UTC()

	claims := CustomClaims{
		data,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(2 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		return "", fmt.Errorf("create: sign token: %w", err)
	}

	return token, nil
}

func (t *Token) Validate(token string) (model.Token, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM(t.PublicKey)
	if err != nil {
		return model.Token{}, fmt.Errorf("validate: parse key: %w", err)
	}

	tok, err := jwt.ParseWithClaims(token, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
		}

		return key, nil
	})
	if err != nil {
		return model.Token{}, fmt.Errorf("validate: %w", err)
	}

	claims, ok := tok.Claims.(*CustomClaims)
	if !ok || !tok.Valid {
		return model.Token{}, fmt.Errorf("validate: invalid")
	}

	return claims.Data, nil
}
