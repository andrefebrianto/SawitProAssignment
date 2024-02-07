package password

import (
	"golang.org/x/crypto/bcrypt"
)

func (p *Password) GenerateHash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), p.hashCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func (*Password) Validate(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return err
	}

	return nil
}
