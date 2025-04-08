package handlers

import (
	"log"

	"github.com/labstack/echo/v4"
)

func RegisterAllHandlers(e *echo.Echo) {
	e.Static("/", "public")
	e.GET("/", HomeHandler)
	e.GET("/login", LoginRegisterContainerHandler)
	e.GET("/loginForm", LoginFormHandler)
	e.GET("/registerForm", RegisterFormHandler)
	log.Println("All handlers registered")
}
