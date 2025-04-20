package handlers

import (
	"log"

	"github.com/fgeck/go-register/internal/repository"
	"github.com/fgeck/go-register/internal/service/loginRegister"
	"github.com/fgeck/go-register/internal/service/security/jwt"
	"github.com/fgeck/go-register/internal/service/security/password"
	"github.com/fgeck/go-register/internal/service/user"
	"github.com/fgeck/go-register/internal/service/validation"
	"github.com/labstack/echo/v4"
)

func SetupHandlers(e *echo.Echo, queries *repository.Queries) {
	validator := validation.NewValidationService()
	userService := user.NewUserService(queries, validator)
	passwordService := password.NewPasswordService()
	jwtService := jwt.NewJwtService("secret", "issuer", 24*1000*60*60) // 24 hours
	loginRegisterService := loginRegister.NewLoginRegisterService(userService, passwordService, jwtService)
	registerHandler := NewRegisterHandler(loginRegisterService)
	loginHandler := NewLoginHandler(loginRegisterService)

	e.Static("/", "public")
	e.GET("/", HomeHandler)
	e.GET("/login", loginHandler.LoginRegisterContainerHandler)
	e.GET("/loginForm", loginHandler.LoginFormHandler)
	e.GET("/registerForm", registerHandler.RegisterFormHandler)
	e.POST("/register", registerHandler.RegisterUserHandler)
	log.Println("All handlers registered")
}
