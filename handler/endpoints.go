package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/SawitProRecruitment/UserService/constant"
	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/model"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/labstack/echo/v4"
)

const tokenTTL = 24 * time.Hour // access token will be expired in 24 hours

func (s *Server) RegisterUser(ctx echo.Context) error {
	reqCtx := ctx.Request().Context()
	request := new(generated.RegisterUserJSONRequestBody)

	if err := ctx.Bind(request); err != nil {
		return ctx.JSON(http.StatusBadRequest, BadRequestErr)
	}

	err := s.Validator.Validate(request)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, BadRequestErr)

	}

	user, err := s.Repository.GetUserByPhone(reqCtx, repository.GetUserByPhoneInput{
		Phone: request.Phone,
	})
	if err != nil && !errors.Is(err, constant.NotFoundErr) {
		return ctx.JSON(http.StatusInternalServerError, InternalErr)
	}
	if user.ID != 0 {
		return ctx.JSON(http.StatusBadRequest, BadRequestErr)
	}

	hashedPassword, err := s.Password.GenerateHash(request.Password)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, BadRequestErr)
	}

	userID, err := s.Repository.SaveUser(reqCtx, repository.SaveUserInput{
		Name:     request.Name,
		Phone:    request.Phone,
		Password: hashedPassword,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, InternalErr)
	}

	return ctx.JSON(http.StatusOK, generated.UserRegistrationResponse{
		Id: userID,
	})
}

func (s *Server) LoginUser(ctx echo.Context) error {
	reqCtx := ctx.Request().Context()
	request := new(generated.LoginUserJSONRequestBody)

	if err := ctx.Bind(request); err != nil {
		return ctx.JSON(http.StatusBadRequest, BadRequestErr)
	}

	err := s.Validator.Validate(request)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, BadRequestErr)

	}

	user, err := s.Repository.GetUserByPhone(reqCtx, repository.GetUserByPhoneInput{
		Phone: request.Phone,
	})
	if err != nil {
		if errors.Is(err, constant.NotFoundErr) {
			return ctx.JSON(http.StatusBadRequest, BadRequestErr)
		}
		return ctx.JSON(http.StatusInternalServerError, InternalErr)
	}

	err = s.Password.Validate(user.Password, request.Password)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, BadRequestErr)
	}

	err = s.Repository.IncreaseUserLoginCounterByID(reqCtx, repository.IncreaseUserLoginCounterByIDInput{
		ID: user.ID,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, InternalErr)
	}

	token, err := s.Token.Generate(tokenTTL, model.Token{
		UserID: user.ID,
		Name:   user.Name,
	})
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, BadRequestErr)
	}

	return ctx.JSON(http.StatusOK, &generated.UserLoginResponse{
		Id:          user.ID,
		AccessToken: token,
	})
}

func (s *Server) GetUserProfile(ctx echo.Context) error {
	reqCtx := ctx.Request().Context()
	session, ok := ctx.Get(constant.AuthContextKey).(model.Token)
	if !ok {
		return ctx.JSON(http.StatusForbidden, ForbiddenErr)
	}

	user, err := s.Repository.GetUserByID(reqCtx, repository.GetUserByIDInput{
		ID: session.UserID,
	})
	if err != nil {
		if errors.Is(err, constant.NotFoundErr) {
			return ctx.JSON(http.StatusForbidden, ForbiddenErr)
		}
		return ctx.JSON(http.StatusInternalServerError, InternalErr)
	}

	return ctx.JSON(http.StatusOK, &generated.UserResponse{
		Name:  user.Name,
		Phone: user.Phone,
	})
}

func (s *Server) UpdateUserProfile(ctx echo.Context) error {
	reqCtx := ctx.Request().Context()
	request := new(generated.UpdateUserProfileJSONRequestBody)
	session, ok := ctx.Get(constant.AuthContextKey).(model.Token)
	if !ok {
		return ctx.JSON(http.StatusForbidden, ForbiddenErr)
	}

	if err := ctx.Bind(request); err != nil {
		return ctx.JSON(http.StatusBadRequest, BadRequestErr)
	}

	err := s.Validator.Validate(request)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, BadRequestErr)

	}

	currentUser, err := s.Repository.GetUserByID(reqCtx, repository.GetUserByIDInput{
		ID: session.UserID,
	})
	if err != nil {
		if errors.Is(err, constant.NotFoundErr) {
			return ctx.JSON(http.StatusForbidden, ForbiddenErr)
		}
		return ctx.JSON(http.StatusInternalServerError, InternalErr)
	}

	if currentUser.Phone != request.Phone {
		existingUser, err := s.Repository.GetUserByPhone(reqCtx, repository.GetUserByPhoneInput{
			Phone: request.Phone,
		})
		if err == nil && existingUser.ID != 0 {
			return ctx.JSON(http.StatusConflict, ConflictErr)
		}
	}

	err = s.Repository.UpdateUserByID(reqCtx, repository.UpdateUserByIDInput{
		ID:    session.UserID,
		Name:  request.Name,
		Phone: request.Phone,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, InternalErr)
	}

	return ctx.JSON(http.StatusOK, &generated.UserResponse{
		Name:  request.Name,
		Phone: request.Phone,
	})
}
