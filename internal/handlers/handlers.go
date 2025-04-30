package handlers

import (
	"log"

	"github.com/fgeck/go-register/internal/repository"
	"github.com/fgeck/go-register/internal/service/config"
	"github.com/fgeck/go-register/internal/service/loginRegister"
	"github.com/fgeck/go-register/internal/service/security/jwt"
	"github.com/fgeck/go-register/internal/service/security/password"
	"github.com/fgeck/go-register/internal/service/user"
	"github.com/fgeck/go-register/internal/service/validation"
	"github.com/labstack/echo/v4"
)

const (
	TWENTY_FOUR_HOURS_IN_MS = 24 * 1000 * 60 * 60
	ISSUER                  = "go-register"
)

func InitServer(echoServer *echo.Echo, queries *repository.Queries, config *config.Config) {
	validator := validation.NewValidationService()
	userService := user.NewUserService(queries, validator)
	passwordService := password.NewPasswordService()
	jwtService := jwt.NewJwtService(config.App.JwtSecret, ISSUER, TWENTY_FOUR_HOURS_IN_MS)
	loginRegisterService := loginRegister.NewLoginRegisterService(userService, passwordService, jwtService)
	registerHandler := NewRegisterHandler(loginRegisterService)
	loginHandler := NewLoginHandler(loginRegisterService)

	echoServer.Static("/", "public")
	echoServer.GET("/", HomeHandler)
	echoServer.GET("/login", loginHandler.LoginRegisterContainerHandler)
	echoServer.GET("/loginForm", loginHandler.LoginFormHandler)
	echoServer.GET("/registerForm", registerHandler.RegisterFormHandler)
	echoServer.POST("/register", registerHandler.RegisterUserHandler)
	log.Println("All handlers registered")
}
