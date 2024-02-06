package token

type Token struct {
	PrivateKey []byte
	PublicKey  []byte
}

func NewToken(privateKey []byte, publicKey []byte) *Token {
	return &Token{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}
}
