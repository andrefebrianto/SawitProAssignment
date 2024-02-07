package handler

import (
	"github.com/SawitProRecruitment/UserService/password"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/token"
	"github.com/SawitProRecruitment/UserService/validator"
)

type Server struct {
	Repository repository.RepositoryInterface
	Token      token.TokenInterface
	Password   password.PasswordInterface
	Validator  validator.ValidatorInterface
}

type NewServerOptions struct {
	Repository repository.RepositoryInterface
	Token      token.TokenInterface
	Password   password.PasswordInterface
	Validator  validator.ValidatorInterface
}

func NewServer(opts NewServerOptions) *Server {
	return &Server{
		Repository: opts.Repository,
		Token:      opts.Token,
		Password:   opts.Password,
		Validator:  opts.Validator,
	}
}
