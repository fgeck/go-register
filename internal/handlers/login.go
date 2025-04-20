package handlers

import (
	"net/http"

	loginregister "github.com/fgeck/go-register/internal/service/loginRegister"
	"github.com/fgeck/go-register/internal/service/render"

	"github.com/fgeck/go-register/templates/views"
	"github.com/labstack/echo/v4"
)

type LoginHandlerInterface interface {
	LoginRegisterContainerHandler(c echo.Context) error
	LoginFormHandler(c echo.Context) error
	LoginHandler(c echo.Context) error
}

type LoginHandler struct {
	loginRegisterService loginregister.LoginRegisterServiceInterface
}

func NewLoginHandler(loginRegisterService loginregister.LoginRegisterServiceInterface) *LoginHandler {
	return &LoginHandler{
		loginRegisterService: loginRegisterService,
	}
}

func (h *LoginHandler) LoginRegisterContainerHandler(c echo.Context) error {
	return render.Render(c, views.LoginRegister())
}

func (h *LoginHandler) LoginFormHandler(c echo.Context) error {
	return render.Render(c, views.LoginForm())
}

func (h *LoginHandler) LoginHandler(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	token, err := h.loginRegisterService.LoginUser(c.Request().Context(), username, password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to login user"})
	}
	return c.JSON(http.StatusOK, map[string]string{"token": token})
}
