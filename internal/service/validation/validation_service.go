package validation

import (
	"fmt"
	"net/http"
	"regexp"
	"unicode"

	customErrors "github.com/fgeck/go-register/internal/service/errors"
)

const (
	PASSWORD_MIN_LENGTH = 8
	USERNAME_MIN_LENGTH = 3
	USERNAME_MAX_LENGTH = 30
	USERNAME_REGEX      = `^[a-zA-Z0-9]+$`
	EMAIL_REGEX         = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
)

type ValidationServiceInterface interface {
	ValidateEmail(email string) error
	ValidatePassword(password string) error
	ValidateUsername(username string) error
}

type ValidationService struct{}

func NewValidationService() *ValidationService {
	return &ValidationService{}
}

func (v *ValidationService) ValidateEmail(email string) error {
	matched, _ := regexp.MatchString(EMAIL_REGEX, email)
	if !matched {
		return customErrors.NewUserFacing("invalid email format", http.StatusBadRequest)
	}

	return nil
}

func (v *ValidationService) ValidatePassword(password string) error {
	var hasMinLen, hasUpper, hasLower, hasNumber, hasSpecial bool

	if len(password) >= PASSWORD_MIN_LENGTH {
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
		return customErrors.NewUserFacing(
			fmt.Sprintf(
				"password must be at least %d characters long and include at least 1 uppercase letter, "+
					"1 lowercase letter, 1 number, and 1 special character",
				PASSWORD_MIN_LENGTH,
			),
			http.StatusBadRequest,
		)
	}

	return nil
}

func (v *ValidationService) ValidateUsername(username string) error {
	if len(username) < USERNAME_MIN_LENGTH || len(username) > USERNAME_MAX_LENGTH {
		return customErrors.NewUserFacing("username must be at least 3 characters long", http.StatusBadRequest)
	}

	for _, char := range username {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) {
			return customErrors.NewUserFacing("username can only contain letters and numbers", http.StatusBadRequest)
		}
	}

	return nil
}
