package middleware

import (
	"github.com/fgeck/go-register/internal/service/security/jwt"
	gojwt "github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func JwtAuthMiddleware(jwtSecret string) echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(jwtSecret),
		TokenLookup: "cookie:token",
		NewClaimsFunc: func(c echo.Context) gojwt.Claims {
			return new(jwt.JwtCustomClaims)
		},
	})
}
