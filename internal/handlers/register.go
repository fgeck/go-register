package handlers

import (
	"net/http"

	"github.com/fgeck/go-register/internal/service/render"
	"github.com/fgeck/go-register/internal/service/user"
	"github.com/fgeck/go-register/templates/views"
	echo "github.com/labstack/echo/v4"
)

type RegisterHandler struct {
	userService *user.UserService
}

func NewRegisterHandler(userService *user.UserService) *RegisterHandler {
	return &RegisterHandler{
		userService: userService,
	}
}

func (r *RegisterHandler) RegisterFormHandler(c echo.Context) error {
	return render.Render(c, views.RegisterForm())
}

func (r *RegisterHandler) RegisterUserHandler(c echo.Context) error {
	username := c.FormValue("username")
	email := c.FormValue("email")
	password := c.FormValue("password")

	if err := r.userService.ValidateCreateUserParams(username, email, password); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	userCreatedDto, err := r.userService.CreateUser(c.Request().Context(), username, email, password)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create user")
	}
	return c.JSON(http.StatusOK, userCreatedDto)

}
