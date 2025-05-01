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
	"github.com/fgeck/go-register/internal/web/middleware"
	"github.com/labstack/echo/v4"
)

const (
	TWENTY_FOUR_HOURS_IN_SECONDS = 24 * 60 * 60
	ISSUER                       = "go-register"
)

func InitServer(echoServer *echo.Echo, queries *repository.Queries, config *config.Config) {
	// Services
	validator := validation.NewValidationService()
	userService := user.NewUserService(queries, validator)
	passwordService := password.NewPasswordService()
	jwtService := jwt.NewJwtService(config.App.JwtSecret, ISSUER, TWENTY_FOUR_HOURS_IN_SECONDS)
	loginRegisterService := loginRegister.NewLoginRegisterService(userService, passwordService, jwtService)
	// Handlers
	registerHandler := NewRegisterHandler(loginRegisterService)
	loginHandler := NewLoginHandler(loginRegisterService)
	// Middlewares
	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	// Setup Server
	echoServer.Static("/", "public")
	echoServer.GET("/", HomeHandler)
	echoServer.GET("/login", loginHandler.LoginRegisterContainerHandler)
	echoServer.GET("/loginForm", loginHandler.LoginFormHandler)
	echoServer.POST("/login", loginHandler.LoginHandler)
	echoServer.GET("/registerForm", registerHandler.RegisterFormHandler)
	echoServer.POST("/register", registerHandler.RegisterUserHandler)
	echoServer.Use(authMiddleware.Authenticate)
	log.Println("All handlers registered")
}
