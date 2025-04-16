package handlers

import (
	"log"

	"github.com/fgeck/go-register/internal/repository"
	"github.com/labstack/echo/v4"
)

func SetupHandlers(e *echo.Echo, repo *repository.Queries) {
	e.Static("/", "public")
	registerHandler := &RegisterHandler{DB: repo}
	e.GET("/", HomeHandler)
	e.GET("/login", LoginRegisterContainerHandler)
	e.GET("/loginForm", LoginFormHandler)
	e.GET("/registerForm", RegisterFormHandler)
	e.POST("/register", RegisterHandler)
	log.Println("All handlers registered")
}
