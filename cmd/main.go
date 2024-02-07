package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/handler"
	"github.com/SawitProRecruitment/UserService/middleware"
	"github.com/SawitProRecruitment/UserService/password"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/token"
	validatorwrapper "github.com/SawitProRecruitment/UserService/validator"
	"github.com/go-playground/validator/v10"
	echomiddleware "github.com/labstack/echo/v4/middleware"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	tokenHandler := initTokenSigningKey()

	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORS())
	e.Use(middleware.ValidateJWT(tokenHandler))

	var server generated.ServerInterface = initServer(tokenHandler)

	generated.RegisterHandlers(e, server)

	errChan := make(chan error)
	stopChan := make(chan os.Signal, 1)

	// bind OS events to the signal channel
	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		errChan <- e.Start(":1323")
	}()

	select {
	case err := <-errChan:
		e.Logger.Error(err)
	case <-stopChan:
		e.Logger.Info("stopping server...")
	}
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
