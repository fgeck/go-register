package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"unicode"

	"github.com/fgeck/go-register/internal/repository"
	"github.com/fgeck/go-register/internal/service/render"
	"github.com/fgeck/go-register/templates/views"
	echo "github.com/labstack/echo/v4"
)

type RegisterHandler struct {
	userService *service.UserService
}

func (r *RegisterHandler) RegisterFormHandler(c echo.Context) error {
	return render.Render(c, views.RegisterForm())
}

func (r *RegisterHandler) RegisterUserHandler(c echo.Context) error {
	username := c.FormValue("username")
	email := c.FormValue("email")
	password := c.FormValue("password")

	if err := r.validateEmail(email); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	if err := r.validatePassword(password); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	if err := r.validateUsername(username); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	exists, err := r.repo.UserExistsByEmail(c.Request().Context(), email)
	if err != nil {
		return err
	}
	if exists {
		return c.String(http.StatusBadRequest, "Email already exists")
	}
	passwordHash, err :=

		r.repo.CreateUser(c.Request().Context(), repository.CreateUserParams{
			Username:     username,
			Email:        email,
			PasswordHash: passwordHash,
		})

	return nil
}

func (r *RegisterHandler) validateEmail(email string) error {
	// Simple regex for email validation
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(emailRegex, email)
	if !matched {
		return errors.New("invalid email format")
	}
	return nil
}

// Validate password (min 8 chars, 1 upper, 1 lower, 1 number, 1 special char)
func (r *RegisterHandler) validatePassword(password string) error {
	var hasMinLen, hasUpper, hasLower, hasNumber, hasSpecial bool

	if len(password) >= minPasswordLength {
		hasMinLen = true
	}

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasMinLen || !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return fmt.Errorf("password must be at least %d characters long and include at least 1 uppercase letter, 1 lowercase letter, 1 number, and 1 special character", minPasswordLength)
	}
	return nil
}

func (r *RegisterHandler) validateUsername(username string) error {
	if len(username) < 3 {
		return errors.New("username must be at least 3 characters long")
	}

	for _, char := range username {
		if !(unicode.IsLetter(char) || unicode.IsDigit(char)) {
			return errors.New("username can only contain letters and numbers")
		}
	}
	return nil
}
