package main

import (
	"log"
	"os"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/handler"
	"github.com/SawitProRecruitment/UserService/password"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/token"
	validatorwrapper "github.com/SawitProRecruitment/UserService/validator"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4/middleware"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	tokenHandler := initTokenSigningKey()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(ValidateJWT(tokenHandler))

	var server generated.ServerInterface = initServer(tokenHandler)

	generated.RegisterHandlers(e, server)
	e.Logger.Fatal(e.Start(":1323"))
}

func initTokenSigningKey() *token.Token {
	prvKey, err := os.ReadFile("cert/id_rsa")
	if err != nil {
		log.Fatalln(err)
	}
	pubKey, err := os.ReadFile("cert/id_rsa.pub")
	if err != nil {
		log.Fatalln(err)
	}

	return token.NewToken(prvKey, pubKey)
}

func initServer(tokenHandler *token.Token) *handler.Server {
	dbDsn := os.Getenv("DATABASE_URL")
	var repo repository.RepositoryInterface = repository.NewRepository(repository.NewRepositoryOptions{
		Dsn: dbDsn,
	})

	passwordHandler := password.NewPassword(12)
	inputValidator := validatorwrapper.NewValidator(validator.New())

	opts := handler.NewServerOptions{
		Repository: repo,
		Token:      tokenHandler,
		Password:   passwordHandler,
		Validator:  inputValidator,
	}
	return handler.NewServer(opts)
}
