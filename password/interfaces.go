package password

type PasswordInterface interface {
	GenerateFromPassword(password string) (string, error)
	VerifyPassword(hashedPassword, password string) error
}
