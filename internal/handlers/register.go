package handlers

import (
	"net/http"

	userfacing_errors "github.com/fgeck/go-register/internal/service/errors"
	"github.com/fgeck/go-register/internal/service/loginRegister"
	"github.com/fgeck/go-register/internal/service/render"
	"github.com/fgeck/go-register/templates/views"
	echo "github.com/labstack/echo/v4"
)

type RegisterHandler struct {
	loginRegisterService loginRegister.LoginRegisterServiceInterface
}

func NewRegisterHandler(loginRegisterService loginRegister.LoginRegisterServiceInterface) *RegisterHandler {
	return &RegisterHandler{
		loginRegisterService: loginRegisterService,
	}
}

func (r *RegisterHandler) RegisterFormHandler(c echo.Context) error {
	return render.Render(c, views.RegisterForm())
}

func (r *RegisterHandler) RegisterUserHandler(c echo.Context) error {
	username := c.FormValue("username")
	email := c.FormValue("email")
	password := c.FormValue("password")

	user, err := r.loginRegisterService.RegisterUser(c.Request().Context(), username, email, password)
	if err != nil {
		if userfacingErr, ok := err.(*userfacing_errors.UserFacingError); ok {
			return c.JSON(userfacingErr.Code, map[string]string{"error": userfacingErr.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to register user"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"user": user})
}
