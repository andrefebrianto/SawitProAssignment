package middleware

import (
	"net/http"
	"strings"

	"github.com/SawitProRecruitment/UserService/constant"
	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/token"

	"github.com/labstack/echo/v4"
)

const (
	authHeader = "Authorization"
	prefix     = "Bearer "
)

func ValidateJWT(tokenHandler *token.Token) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			value := c.Request().Header.Values(authHeader)

			if len(value) > 0 && strings.HasPrefix(value[0], prefix) {
				accessToken := strings.TrimPrefix(value[0], prefix)
				claim, err := tokenHandler.Validate(accessToken)
				if err != nil {
					return c.JSON(http.StatusUnauthorized, &generated.Error{Message: "The request lacks valid authentication credentials for the requested resource"})
				}

				c.Set(constant.AuthContextKey, claim)
			}

			return next(c)
		}
	}
}
