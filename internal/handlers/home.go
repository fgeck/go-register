package handlers

import (
	"github.com/fgeck/go-register/internal/render"
	"github.com/fgeck/go-register/templates/views"
	"github.com/labstack/echo/v4"
)

func HomeHandler(c echo.Context) error {
	return render.Render(c, views.Home("Home"))
}
