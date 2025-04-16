package handlers

import (
	"log"

	"github.com/fgeck/go-register/internal/repository"
	"github.com/fgeck/go-register/internal/service/user"
	"github.com/fgeck/go-register/internal/service/validation"
	"github.com/labstack/echo/v4"
)

func SetupHandlers(e *echo.Echo, queries *repository.Queries) {
	e.Static("/", "public")
	validator := validation.NewValidationService()
	userService := user.NewUserService(queries, validator)
	registerHandler := NewRegisterHandler(userService)
	e.GET("/", HomeHandler)
	e.GET("/login", LoginRegisterContainerHandler)
	e.GET("/loginForm", LoginFormHandler)
	e.GET("/registerForm", registerHandler.RegisterFormHandler)
	e.POST("/register", registerHandler.RegisterUserHandler)
	log.Println("All handlers registered")
}
