package handler_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/SawitProRecruitment/UserService/constant"
	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/handler"
	"github.com/SawitProRecruitment/UserService/model"
	"github.com/SawitProRecruitment/UserService/password"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/token"
	"github.com/SawitProRecruitment/UserService/validator"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func generateResponseString(resp interface{}) string {
	respJSON, _ := json.Marshal(resp)
	return string(respJSON) + "\n" // echo append this escape sequence at the end of the body
}

var _ = Describe("API Handler", func() {
	var (
		mockController *gomock.Controller

		repoMock      *repository.MockRepositoryInterface
		tokenMock     *token.MockTokenInterface
		passwordMock  *password.MockPasswordInterface
		validatorMock *validator.MockValidatorInterface
		server        *handler.Server

		echoCtx         echo.Context
		responseRec     *httptest.ResponseRecorder
		plainPassword   = "very_secure_password"
		hashedPassword  = "$2a$12$GknIDMsf25xSzcevHPople2zvckIn57bKzpr89..1pLjz53Fw8cte"
		tokenTTL        = 24 * time.Hour
		databaseConnErr = errors.New("database connection error")
	)

	Describe("Register new user", func() {
		var (
			requestBody generated.RegisterUserJSONRequestBody
		)

		BeforeEach(func() {
			gofakeit.Seed(time.Now().UnixNano())
			mockController = gomock.NewController(GinkgoT())

			repoMock = repository.NewMockRepositoryInterface(mockController)
			tokenMock = token.NewMockTokenInterface(mockController)
			passwordMock = password.NewMockPasswordInterface(mockController)
			validatorMock = validator.NewMockValidatorInterface(mockController)
			opts := handler.NewServerOptions{
				Repository: repoMock,
				Token:      tokenMock,
				Password:   passwordMock,
				Validator:  validatorMock,
			}
			server = handler.NewServer(opts)

			gofakeit.Struct(&requestBody)

			requestBodyJSON, _ := json.Marshal(requestBody)

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/users/registrations", strings.NewReader(string(requestBodyJSON)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			responseRec = httptest.NewRecorder()
			echoCtx = e.NewContext(req, responseRec)
		})

		AfterEach(func() {
			mockController.Finish()
		})

		When("the user sends input that does not match the requirements", func() {
			Context("input has different data type", func() {
				It("tells the user that input is invalid", func() {
					ec := echo.New()
					req := httptest.NewRequest(http.MethodPost, "/users/registrations", strings.NewReader(`{"name":"John Doe","password":"secretpassword","phone":123}`))
					req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
					responseRec = httptest.NewRecorder()
					ecCtx := ec.NewContext(req, responseRec)

					err := server.RegisterUser(ecCtx)

					resBody := generateResponseString(handler.BadRequestErr)
					Expect(err).Should(BeNil())
					Expect(responseRec.Code).Should(Equal(http.StatusBadRequest))
					Expect(responseRec.Body.String()).Should(Equal(resBody))
				})
			})

			It("tells the user that input is invalid", func() {
				validatorMock.EXPECT().Validate(&requestBody).Return(errors.New("Key: 'RegisterUserJSONRequestBody.Name' Error:Field validation for 'Name' failed on the 'min' tag"))

				err := server.RegisterUser(echoCtx)

				resBody := generateResponseString(handler.BadRequestErr)
				Expect(err).Should(BeNil())
				Expect(responseRec.Code).Should(Equal(http.StatusBadRequest))
				Expect(responseRec.Body.String()).Should(Equal(resBody))
			})
		})

		When("the user sends phone number that exist in database", func() {
			Context("there is a problem with the database connection when retrieving existing user", func() {
				It("tells the user that internal error occur", func() {
					ctx := echoCtx.Request().Context()

					validatorMock.EXPECT().Validate(&requestBody).Return(nil)
					repoMock.EXPECT().GetUserByPhone(ctx, repository.GetUserByPhoneInput{Phone: requestBody.Phone}).Return(repository.GetUserByPhoneOutput{}, databaseConnErr)

					err := server.RegisterUser(echoCtx)

					resBody := generateResponseString(handler.InternalErr)
					Expect(err).Should(BeNil())
					Expect(responseRec.Code).Should(Equal(http.StatusInternalServerError))
					Expect(responseRec.Body.String()).Should(Equal(resBody))
				})
			})

			It("tells the user that input is invalid", func() {
				ctx := echoCtx.Request().Context()

				var userByPhoneOutput repository.GetUserByPhoneOutput
				gofakeit.Struct(&userByPhoneOutput)
				userByPhoneOutput.Phone = requestBody.Phone

				validatorMock.EXPECT().Validate(&requestBody).Return(nil)
				repoMock.EXPECT().GetUserByPhone(ctx, repository.GetUserByPhoneInput{Phone: requestBody.Phone}).Return(userByPhoneOutput, nil)

				err := server.RegisterUser(echoCtx)

				resBody := generateResponseString(handler.BadRequestErr)
				Expect(err).Should(BeNil())
				Expect(responseRec.Code).Should(Equal(http.StatusBadRequest))
				Expect(responseRec.Body.String()).Should(Equal(resBody))
			})

		})

		When("the user sends input that match the requirements", func() {
			Context("there is a problem when generating password hash", func() {
				It("tells the user that input is invalid", func() {
					ctx := echoCtx.Request().Context()

					validatorMock.EXPECT().Validate(&requestBody).Return(nil)
					repoMock.EXPECT().GetUserByPhone(ctx, repository.GetUserByPhoneInput{Phone: requestBody.Phone}).Return(repository.GetUserByPhoneOutput{}, nil)
					passwordMock.EXPECT().GenerateHash(requestBody.Password).Return("", errors.New("hash password error"))

					err := server.RegisterUser(echoCtx)

					resBody := generateResponseString(handler.BadRequestErr)
					Expect(err).Should(BeNil())
					Expect(responseRec.Code).Should(Equal(http.StatusBadRequest))
					Expect(responseRec.Body.String()).Should(Equal(resBody))
				})
			})

			Context("there is a problem with the database connection when saving new user", func() {
				It("tells the user that internal error occur", func() {
					ctx := echoCtx.Request().Context()

					validatorMock.EXPECT().Validate(&requestBody).Return(nil)
					repoMock.EXPECT().GetUserByPhone(ctx, repository.GetUserByPhoneInput{Phone: requestBody.Phone}).Return(repository.GetUserByPhoneOutput{}, nil)
					passwordMock.EXPECT().GenerateHash(requestBody.Password).Return(hashedPassword, nil)
					repoMock.EXPECT().SaveUser(ctx, repository.SaveUserInput{
						Name:     requestBody.Name,
						Phone:    requestBody.Phone,
						Password: hashedPassword,
					}).Return(int64(0), databaseConnErr)

					err := server.RegisterUser(echoCtx)

					resBody := generateResponseString(handler.InternalErr)
					Expect(err).Should(BeNil())
					Expect(responseRec.Code).Should(Equal(http.StatusInternalServerError))
					Expect(responseRec.Body.String()).Should(Equal(resBody))
				})
			})

			It("sends user id", func() {
				ctx := echoCtx.Request().Context()
				userID := int64(1)

				validatorMock.EXPECT().Validate(&requestBody).Return(nil)
				repoMock.EXPECT().GetUserByPhone(ctx, repository.GetUserByPhoneInput{Phone: requestBody.Phone}).Return(repository.GetUserByPhoneOutput{}, nil)
				passwordMock.EXPECT().GenerateHash(requestBody.Password).Return(hashedPassword, nil)
				repoMock.EXPECT().SaveUser(ctx, repository.SaveUserInput{
					Name:     requestBody.Name,
					Phone:    requestBody.Phone,
					Password: hashedPassword,
				}).Return(userID, nil)

				err := server.RegisterUser(echoCtx)

				resBody := generateResponseString(generated.UserRegistrationResponse{Id: userID})
				Expect(err).Should(BeNil())
				Expect(responseRec.Code).Should(Equal(http.StatusOK))
				Expect(responseRec.Body.String()).Should(Equal(resBody))
			})
		})
	})

	Describe("Login with user credential", func() {
		var (
			requestBody generated.LoginUserJSONRequestBody
		)

		BeforeEach(func() {
			gofakeit.Seed(time.Now().UnixNano())
			mockController = gomock.NewController(GinkgoT())

			repoMock = repository.NewMockRepositoryInterface(mockController)
			tokenMock = token.NewMockTokenInterface(mockController)
			passwordMock = password.NewMockPasswordInterface(mockController)
			validatorMock = validator.NewMockValidatorInterface(mockController)
			opts := handler.NewServerOptions{
				Repository: repoMock,
				Token:      tokenMock,
				Password:   passwordMock,
				Validator:  validatorMock,
			}
			server = handler.NewServer(opts)

			gofakeit.Struct(&requestBody)
			requestBody.Password = plainPassword

			requestBodyJSON, _ := json.Marshal(requestBody)

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/users/login", strings.NewReader(string(requestBodyJSON)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			responseRec = httptest.NewRecorder()
			echoCtx = e.NewContext(req, responseRec)
		})

		AfterEach(func() {
			mockController.Finish()
		})

		When("the user sends input that does not match the requirements", func() {
			Context("input has different data type", func() {
				It("tells the user that input is invalid", func() {
					ec := echo.New()
					req := httptest.NewRequest(http.MethodPost, "/users/login", strings.NewReader(`{"password":"secretpassword","phone":123}`))
					req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
					responseRec = httptest.NewRecorder()
					ecCtx := ec.NewContext(req, responseRec)

					err := server.LoginUser(ecCtx)

					resBody := generateResponseString(handler.BadRequestErr)
					Expect(err).Should(BeNil())
					Expect(responseRec.Code).Should(Equal(http.StatusBadRequest))
					Expect(responseRec.Body.String()).Should(Equal(resBody))
				})
			})

			It("tells the user that input is invalid", func() {
				validatorMock.EXPECT().Validate(&requestBody).Return(errors.New("Key: 'LoginUserJSONRequestBody.Name' Error:Field validation for 'Phone' failed on the 'min' tag"))

				err := server.LoginUser(echoCtx)

				resBody := generateResponseString(handler.BadRequestErr)
				Expect(err).Should(BeNil())
				Expect(responseRec.Code).Should(Equal(http.StatusBadRequest))
				Expect(responseRec.Body.String()).Should(Equal(resBody))
			})
		})

		When("the user sends phone number that not exist in database", func() {
			Context("there is a problem with the database connection when retrieving existing user", func() {
				It("tells the user that internal error occur", func() {
					ctx := echoCtx.Request().Context()

					validatorMock.EXPECT().Validate(&requestBody).Return(nil)
					repoMock.EXPECT().GetUserByPhone(ctx, repository.GetUserByPhoneInput{Phone: requestBody.Phone}).Return(repository.GetUserByPhoneOutput{}, databaseConnErr)

					err := server.LoginUser(echoCtx)

					resBody := generateResponseString(handler.InternalErr)
					Expect(err).Should(BeNil())
					Expect(responseRec.Code).Should(Equal(http.StatusInternalServerError))
					Expect(responseRec.Body.String()).Should(Equal(resBody))
				})
			})

			It("tells the user that input is invalid", func() {
				ctx := echoCtx.Request().Context()

				validatorMock.EXPECT().Validate(&requestBody).Return(nil)
				repoMock.EXPECT().GetUserByPhone(ctx, repository.GetUserByPhoneInput{Phone: requestBody.Phone}).Return(repository.GetUserByPhoneOutput{}, constant.NotFoundErr)

				err := server.LoginUser(echoCtx)

				resBody := generateResponseString(handler.BadRequestErr)
				Expect(err).Should(BeNil())
				Expect(responseRec.Code).Should(Equal(http.StatusBadRequest))
				Expect(responseRec.Body.String()).Should(Equal(resBody))
			})
		})

		When("the user sends input that match the requirements but combination of phone and password are incorect", func() {
			It("tells the user that input is invalid", func() {
				ctx := echoCtx.Request().Context()

				var userByPhoneOutput repository.GetUserByPhoneOutput
				gofakeit.Struct(&userByPhoneOutput)
				userByPhoneOutput.Phone = requestBody.Phone
				userByPhoneOutput.Password = hashedPassword

				validatorMock.EXPECT().Validate(&requestBody).Return(nil)
				repoMock.EXPECT().GetUserByPhone(ctx, repository.GetUserByPhoneInput{Phone: requestBody.Phone}).Return(userByPhoneOutput, nil)
				passwordMock.EXPECT().Validate(userByPhoneOutput.Password, requestBody.Password).Return(errors.New("crypto/bcrypt: hashedPassword is not the hash of the given password"))

				err := server.LoginUser(echoCtx)

				resBody := generateResponseString(handler.BadRequestErr)
				Expect(err).Should(BeNil())
				Expect(responseRec.Code).Should(Equal(http.StatusBadRequest))
				Expect(responseRec.Body.String()).Should(Equal(resBody))
			})
		})

		When("the user sends input that match the requirements and combination of phone and password are corect", func() {
			Context("there is a problem with the database connection when increasing login counter", func() {
				It("tells the user that internal error occur", func() {
					ctx := echoCtx.Request().Context()

					requestBody.Password = plainPassword

					var userByPhoneOutput repository.GetUserByPhoneOutput
					gofakeit.Struct(&userByPhoneOutput)
					userByPhoneOutput.Phone = requestBody.Phone
					userByPhoneOutput.Password = hashedPassword

					validatorMock.EXPECT().Validate(&requestBody).Return(nil)
					repoMock.EXPECT().GetUserByPhone(ctx, repository.GetUserByPhoneInput{Phone: requestBody.Phone}).Return(userByPhoneOutput, nil)
					passwordMock.EXPECT().Validate(userByPhoneOutput.Password, requestBody.Password).Return(nil)
					repoMock.EXPECT().IncreaseUserLoginCounterByID(ctx, repository.IncreaseUserLoginCounterByIDInput{ID: userByPhoneOutput.ID}).Return(databaseConnErr)

					err := server.LoginUser(echoCtx)

					resBody := generateResponseString(handler.InternalErr)
					Expect(err).Should(BeNil())
					Expect(responseRec.Code).Should(Equal(http.StatusInternalServerError))
					Expect(responseRec.Body.String()).Should(Equal(resBody))
				})
			})

			Context("there is a problem when generating access token", func() {
				It("tells the user that input is invalid", func() {
					ctx := echoCtx.Request().Context()

					requestBody.Password = plainPassword

					var userByPhoneOutput repository.GetUserByPhoneOutput
					gofakeit.Struct(&userByPhoneOutput)
					userByPhoneOutput.Phone = requestBody.Phone
					userByPhoneOutput.Password = hashedPassword

					validatorMock.EXPECT().Validate(&requestBody).Return(nil)
					repoMock.EXPECT().GetUserByPhone(ctx, repository.GetUserByPhoneInput{Phone: requestBody.Phone}).Return(userByPhoneOutput, nil)
					passwordMock.EXPECT().Validate(userByPhoneOutput.Password, requestBody.Password).Return(nil)
					repoMock.EXPECT().IncreaseUserLoginCounterByID(ctx, repository.IncreaseUserLoginCounterByIDInput{ID: userByPhoneOutput.ID}).Return(nil)
					tokenMock.EXPECT().Generate(tokenTTL, model.Token{
						UserID: userByPhoneOutput.ID,
						Name:   userByPhoneOutput.Name,
					}).Return("", errors.New("sign token error"))

					err := server.LoginUser(echoCtx)

					resBody := generateResponseString(handler.BadRequestErr)
					Expect(err).Should(BeNil())
					Expect(responseRec.Code).Should(Equal(http.StatusBadRequest))
					Expect(responseRec.Body.String()).Should(Equal(resBody))
				})
			})

			It("sends user id and access token", func() {
				ctx := echoCtx.Request().Context()

				var userByPhoneOutput repository.GetUserByPhoneOutput
				gofakeit.Struct(&userByPhoneOutput)
				userByPhoneOutput.Phone = requestBody.Phone
				userByPhoneOutput.Password = hashedPassword

				accessToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhIjp7IlVzZXJJRCI6MSwiTmFtZSI6IkpvaG4gTGFyayJ9LCJleHAiOjE3MDcyMjQwMDgsIm5iZiI6MTcwNzIxNjgwOCwiaWF0IjoxNzA3MjE2ODA4fQ.n9hEBwJpPosCTYDZcMH57780HS0i53Nazz3fEqmUUpw-Ccml98SwO3lQlCpDhXbKZbSyZR27aj7CLJ1I69ikvb4PbIlJUMFyxCQHiZhq71agQmvb96xVYNBTmCPkrzYCzDYD5w9vc0vmbBehFu6TlL7pX9DaG3HYMFDFh-3ZXtv9Nn6WjibLh8UUrK2gZSownKkLIK0boJWvU4q5OnoXonJKeWyYPRRusXMRooJPytJ3nOiMr-0OdywJ3fSjYOgpW6mbPF5BxZI_hGZVfo88NZ1tOpQIM2zA-6PVd-hq72MO5E811aBJR5Mt_FdzbDv1DECLFT8gKG7XKQBt5n70Dss4s0NgiJJXfJCZqzJI7pTbA77zpsmMvEJ18UPJJvKnsi_y45Ne_raV00WNGVsuFxrYPB4HtneLzRfRvjN-QGFsxAE-kXcuResqVtiv3PJ04czodhNbiVuIOFIl48jYkrtsKG3qBSNxxcwt8p5YonAOwWaquFsty3_vRb_O9Fq-NGXzCpfy-ve469lkxCiPeRlVlQsz_9Kh_QVQHszkuIHsTl1pFu0zbbUrBAtUi-pmN73C08f6VhX-CkDShxylEj-8Tusn8qsKd95hUlp8mVeu6yXUW2_b0YZmWZ2QCCCfkwbLwFAVbVtfEoGjq1Pa3zSSNTuUUJDYf5j37gqnQKg"

				validatorMock.EXPECT().Validate(&requestBody).Return(nil)
				repoMock.EXPECT().GetUserByPhone(ctx, repository.GetUserByPhoneInput{Phone: requestBody.Phone}).Return(userByPhoneOutput, nil)
				passwordMock.EXPECT().Validate(userByPhoneOutput.Password, requestBody.Password).Return(nil)
				repoMock.EXPECT().IncreaseUserLoginCounterByID(ctx, repository.IncreaseUserLoginCounterByIDInput{ID: userByPhoneOutput.ID}).Return(nil)
				tokenMock.EXPECT().Generate(tokenTTL, model.Token{
					UserID: userByPhoneOutput.ID,
					Name:   userByPhoneOutput.Name,
				}).Return(accessToken, nil)

				err := server.LoginUser(echoCtx)

				resBody := generateResponseString(generated.UserLoginResponse{
					Id:          userByPhoneOutput.ID,
					AccessToken: accessToken,
				})
				Expect(err).Should(BeNil())
				Expect(responseRec.Code).Should(Equal(http.StatusOK))
				Expect(responseRec.Body.String()).Should(Equal(resBody))
			})
		})
	})

	Describe("Get current user profile", func() {
		var (
			session model.Token
		)

		BeforeEach(func() {
			gofakeit.Seed(time.Now().UnixNano())
			mockController = gomock.NewController(GinkgoT())

			repoMock = repository.NewMockRepositoryInterface(mockController)
			tokenMock = token.NewMockTokenInterface(mockController)
			passwordMock = password.NewMockPasswordInterface(mockController)
			validatorMock = validator.NewMockValidatorInterface(mockController)
			opts := handler.NewServerOptions{
				Repository: repoMock,
				Token:      tokenMock,
				Password:   passwordMock,
				Validator:  validatorMock,
			}
			server = handler.NewServer(opts)

			gofakeit.Struct(&session)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/users/profile", nil)
			responseRec = httptest.NewRecorder()
			echoCtx = e.NewContext(req, responseRec)
			echoCtx.Set(constant.AuthContextKey, session)
		})

		AfterEach(func() {
			mockController.Finish()
		})

		When("the user doesn't send access token", func() {
			It("tells the user that access to the resource is forbidden", func() {
				echoCtx.Set(constant.AuthContextKey, nil)

				err := server.GetUserProfile(echoCtx)

				resBody := generateResponseString(handler.ForbiddenErr)
				Expect(err).Should(BeNil())
				Expect(responseRec.Code).Should(Equal(http.StatusForbidden))
				Expect(responseRec.Body.String()).Should(Equal(resBody))
			})
		})

		When("the user send access token with invalid user id data", func() {
			Context("there is a problem with the database connection when retrieving current user", func() {
				It("tells the user that internal error occur", func() {
					ctx := echoCtx.Request().Context()

					repoMock.EXPECT().GetUserByID(ctx, repository.GetUserByIDInput{ID: session.UserID}).Return(repository.GetUserByIDOutput{}, databaseConnErr)

					err := server.GetUserProfile(echoCtx)

					resBody := generateResponseString(handler.InternalErr)
					Expect(err).Should(BeNil())
					Expect(responseRec.Code).Should(Equal(http.StatusInternalServerError))
					Expect(responseRec.Body.String()).Should(Equal(resBody))
				})
			})

			It("tells the user that access to the resource is forbidden", func() {
				ctx := echoCtx.Request().Context()

				repoMock.EXPECT().GetUserByID(ctx, repository.GetUserByIDInput{ID: session.UserID}).Return(repository.GetUserByIDOutput{}, constant.NotFoundErr)

				err := server.GetUserProfile(echoCtx)

				resBody := generateResponseString(handler.ForbiddenErr)
				Expect(err).Should(BeNil())
				Expect(responseRec.Code).Should(Equal(http.StatusForbidden))
				Expect(responseRec.Body.String()).Should(Equal(resBody))
			})
		})

		When("the user send valid access token", func() {
			It("sends user name and phone", func() {
				ctx := echoCtx.Request().Context()

				var userByIDOutput repository.GetUserByIDOutput
				gofakeit.Struct(&userByIDOutput)

				repoMock.EXPECT().GetUserByID(ctx, repository.GetUserByIDInput{ID: session.UserID}).Return(userByIDOutput, nil)

				err := server.GetUserProfile(echoCtx)

				resBody := generateResponseString(generated.UserResponse{
					Name:  userByIDOutput.Name,
					Phone: userByIDOutput.Phone,
				})
				Expect(err).Should(BeNil())
				Expect(responseRec.Code).Should(Equal(http.StatusOK))
				Expect(responseRec.Body.String()).Should(Equal(resBody))
			})
		})
	})

	Describe("Update current user profile", func() {
		var (
			requestBody generated.UpdateUserProfileJSONRequestBody
			session     model.Token
		)

		BeforeEach(func() {
			gofakeit.Seed(time.Now().UnixNano())
			mockController = gomock.NewController(GinkgoT())

			repoMock = repository.NewMockRepositoryInterface(mockController)
			tokenMock = token.NewMockTokenInterface(mockController)
			passwordMock = password.NewMockPasswordInterface(mockController)
			validatorMock = validator.NewMockValidatorInterface(mockController)
			opts := handler.NewServerOptions{
				Repository: repoMock,
				Token:      tokenMock,
				Password:   passwordMock,
				Validator:  validatorMock,
			}
			server = handler.NewServer(opts)

			gofakeit.Struct(&requestBody)
			requestBody.Phone = "+6281234567890"
			gofakeit.Struct(&session)

			requestBodyJSON, _ := json.Marshal(requestBody)

			e := echo.New()
			req := httptest.NewRequest(http.MethodPut, "/users/profile", strings.NewReader(string(requestBodyJSON)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			responseRec = httptest.NewRecorder()
			echoCtx = e.NewContext(req, responseRec)
			echoCtx.Set(constant.AuthContextKey, session)
		})

		AfterEach(func() {
			mockController.Finish()
		})

		When("the user doesn't send access token", func() {
			It("tells the user that access to the resource is forbidden", func() {
				echoCtx.Set(constant.AuthContextKey, nil)

				err := server.UpdateUserProfile(echoCtx)

				resBody := generateResponseString(handler.ForbiddenErr)
				Expect(err).Should(BeNil())
				Expect(responseRec.Code).Should(Equal(http.StatusForbidden))
				Expect(responseRec.Body.String()).Should(Equal(resBody))
			})
		})

		When("the user sends input that does not match the requirements", func() {
			Context("input has different data type", func() {
				It("tells the user that input is invalid", func() {
					ec := echo.New()
					req := httptest.NewRequest(http.MethodPut, "/users/profile", strings.NewReader(`{"name":"Jon Snow","phone":123}`))
					req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
					responseRec = httptest.NewRecorder()
					ecCtx := ec.NewContext(req, responseRec)
					ecCtx.Set(constant.AuthContextKey, session)

					err := server.UpdateUserProfile(ecCtx)

					resBody := generateResponseString(handler.BadRequestErr)
					Expect(err).Should(BeNil())
					Expect(responseRec.Code).Should(Equal(http.StatusBadRequest))
					Expect(responseRec.Body.String()).Should(Equal(resBody))
				})
			})

			It("tells the user that input is invalid", func() {
				validatorMock.EXPECT().Validate(&requestBody).Return(errors.New("Key: 'UpdateUserProfileJSONRequestBody.Name' Error:Field validation for 'Phone' failed on the 'min' tag"))

				err := server.UpdateUserProfile(echoCtx)

				resBody := generateResponseString(handler.BadRequestErr)
				Expect(err).Should(BeNil())
				Expect(responseRec.Code).Should(Equal(http.StatusBadRequest))
				Expect(responseRec.Body.String()).Should(Equal(resBody))
			})
		})

		When("the user sends correct input with invalid access token", func() {
			Context("there is a problem with the database connection when retrieving current user", func() {
				It("tells the user that internal error occur", func() {
					ctx := echoCtx.Request().Context()

					validatorMock.EXPECT().Validate(&requestBody).Return(nil)
					repoMock.EXPECT().GetUserByID(ctx, repository.GetUserByIDInput{ID: session.UserID}).Return(repository.GetUserByIDOutput{}, databaseConnErr)

					err := server.UpdateUserProfile(echoCtx)

					resBody := generateResponseString(handler.InternalErr)
					Expect(err).Should(BeNil())
					Expect(responseRec.Code).Should(Equal(http.StatusInternalServerError))
					Expect(responseRec.Body.String()).Should(Equal(resBody))
				})
			})

			It("tells the user that access to the resource is forbidden", func() {
				ctx := echoCtx.Request().Context()

				validatorMock.EXPECT().Validate(&requestBody).Return(nil)
				repoMock.EXPECT().GetUserByID(ctx, repository.GetUserByIDInput{ID: session.UserID}).Return(repository.GetUserByIDOutput{}, constant.NotFoundErr)

				err := server.UpdateUserProfile(echoCtx)

				resBody := generateResponseString(handler.ForbiddenErr)
				Expect(err).Should(BeNil())
				Expect(responseRec.Code).Should(Equal(http.StatusForbidden))
				Expect(responseRec.Body.String()).Should(Equal(resBody))
			})
		})

		When("the user sends the phone number that belongs to another user", func() {
			It("tells the user that conflict error occur", func() {
				ctx := echoCtx.Request().Context()

				userByIDOutput := repository.GetUserByIDOutput{
					ID:    session.UserID,
					Name:  requestBody.Name,
					Phone: "+6281122334455",
				}
				var userByPhoneOutput repository.GetUserByPhoneOutput
				gofakeit.Struct(&userByPhoneOutput)
				userByPhoneOutput.Phone = requestBody.Phone

				validatorMock.EXPECT().Validate(&requestBody).Return(nil)
				repoMock.EXPECT().GetUserByID(ctx, repository.GetUserByIDInput{ID: session.UserID}).Return(userByIDOutput, nil)
				repoMock.EXPECT().GetUserByPhone(ctx, repository.GetUserByPhoneInput{Phone: requestBody.Phone}).Return(userByPhoneOutput, nil)

				err := server.UpdateUserProfile(echoCtx)

				resBody := generateResponseString(handler.ConflictErr)
				Expect(err).Should(BeNil())
				Expect(responseRec.Code).Should(Equal(http.StatusConflict))
				Expect(responseRec.Body.String()).Should(Equal(resBody))
			})
		})

		When("the user sends name and unique new phone number", func() {
			Context("there is a problem with the database connection when updating current user", func() {
				It("tells the user that internal error occur", func() {
					ctx := echoCtx.Request().Context()

					userByIDOutput := repository.GetUserByIDOutput{
						ID:    session.UserID,
						Name:  requestBody.Name,
						Phone: "+6281122334455",
					}

					validatorMock.EXPECT().Validate(&requestBody).Return(nil)
					repoMock.EXPECT().GetUserByID(ctx, repository.GetUserByIDInput{ID: session.UserID}).Return(userByIDOutput, nil)
					repoMock.EXPECT().GetUserByPhone(ctx, repository.GetUserByPhoneInput{Phone: requestBody.Phone}).Return(repository.GetUserByPhoneOutput{}, nil)
					repoMock.EXPECT().UpdateUserByID(ctx, repository.UpdateUserByIDInput{
						ID:    session.UserID,
						Name:  requestBody.Name,
						Phone: requestBody.Phone,
					}).Return(databaseConnErr)

					err := server.UpdateUserProfile(echoCtx)

					resBody := generateResponseString(handler.InternalErr)
					Expect(err).Should(BeNil())
					Expect(responseRec.Code).Should(Equal(http.StatusInternalServerError))
					Expect(responseRec.Body.String()).Should(Equal(resBody))
				})
			})

			It("sends updated profile", func() {
				ctx := echoCtx.Request().Context()

				userByIDOutput := repository.GetUserByIDOutput{
					ID:    session.UserID,
					Name:  requestBody.Name,
					Phone: "+6281122334455",
				}

				validatorMock.EXPECT().Validate(&requestBody).Return(nil)
				repoMock.EXPECT().GetUserByID(ctx, repository.GetUserByIDInput{ID: session.UserID}).Return(userByIDOutput, nil)
				repoMock.EXPECT().GetUserByPhone(ctx, repository.GetUserByPhoneInput{Phone: requestBody.Phone}).Return(repository.GetUserByPhoneOutput{}, nil)
				repoMock.EXPECT().UpdateUserByID(ctx, repository.UpdateUserByIDInput{
					ID:    session.UserID,
					Name:  requestBody.Name,
					Phone: requestBody.Phone,
				}).Return(nil)

				err := server.UpdateUserProfile(echoCtx)

				resBody := generateResponseString(requestBody)
				Expect(err).Should(BeNil())
				Expect(responseRec.Code).Should(Equal(http.StatusOK))
				Expect(responseRec.Body.String()).Should(Equal(resBody))
			})
		})
	})
})
