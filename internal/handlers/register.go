package handlers

import (
	"github.com/fgeck/go-register/internal/render"
	"github.com/fgeck/go-register/templates/views"
	echo "github.com/labstack/echo/v4"
)

func RegisterFormHandler(c echo.Context) error {
	return render.Render(c, views.RegisterForm())
}
