package validation

import (
	"errors"
	"fmt"
	"regexp"
	"unicode"
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

const (
	minPasswordLength = 8
)

func (v *ValidationService) ValidateEmail(email string) error {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(emailRegex, email)
	if !matched {
		return errors.New("invalid email format")
	}
	return nil
}

func (v *ValidationService) ValidatePassword(password string) error {
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

func (v *ValidationService) ValidateUsername(username string) error {
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
