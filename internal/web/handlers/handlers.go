package handlers

import (
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
	TWENTY_FOUR_HOURS_IN_SECONDS = 24 * 60 * 60
	ISSUER                       = "go-register"
)

func InitServer(e *echo.Echo, queries *repository.Queries, config *config.Config) {
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

	// Public Routes
	e.Static("/", "public")
	e.GET("/", HomeHandler)
	e.GET("/login", loginHandler.LoginRegisterContainerHandler)
	e.GET("/loginForm", loginHandler.LoginFormHandler)
	e.POST("/api/login", loginHandler.LoginHandler)
	e.GET("/registerForm", registerHandler.RegisterFormHandler)
	e.POST("/api/register", registerHandler.RegisterUserHandler)

	// Protected Routes (requires authentication)
	// protectedGroup := e.Group("/api/protected")
	// protectedGroup.Use(authMiddleware)
	// protectedGroup.GET("/profile", ProfileHandler)

	// Admin Routes (requires "UserRole" == "admin")
	// adminGroup := e.Group("/api/admin")
	// adminGroup.Use(authMiddleware, adminMiddleware)
}
