package password

type Password struct {
	hashCost int
}

func NewPassword(hashCost int) *Password {
	return &Password{
		hashCost: hashCost,
	}
}
