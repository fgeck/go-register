package handlers

import (
	"github.com/fgeck/go-register/internal/service/render"
	views "github.com/fgeck/go-register/templates/views"
	echo "github.com/labstack/echo/v4"
)

func HomeHandler(c echo.Context) error {
	return render.Render(c, views.Home())
}
