package handlers

import (
	"github.com/fgeck/go-register/internal/render"
	"github.com/fgeck/go-register/templates/views"
	"github.com/labstack/echo/v4"
)

func LoginRegisterContainerHandler(c echo.Context) error {
	return render.Render(c, views.LoginRegister())
}

func LoginFormHandler(c echo.Context) error {
	return render.Render(c, views.LoginForm())
}
