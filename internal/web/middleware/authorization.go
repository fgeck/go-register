package middleware

import (
	"strings"

	"github.com/fgeck/go-register/internal/service/security/jwt"
	user "github.com/fgeck/go-register/internal/service/user"
	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func RequireAdminMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Get("user").(*gojwt.Token)
			claims := token.Claims.(*jwt.JwtCustomClaims)
			if claims.UserRole == "" || strings.ToUpper(claims.UserRole) != user.UserRoleAdmin.Name {
				return echo.ErrForbidden
			}
			return next(c)
		}
	}
}
