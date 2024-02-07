package password

type PasswordInterface interface {
	GenerateHash(password string) (string, error)
	Validate(hashedPassword, password string) error
}
