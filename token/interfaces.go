package token

import (
	"time"

	"github.com/SawitProRecruitment/UserService/model"
)

type TokenInterface interface {
	Generate(ttl time.Duration, data model.Token) (string, error)
	Validate(token string) (model.Token, error)
}
