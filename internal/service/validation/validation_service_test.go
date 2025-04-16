package service_test

import (
	"testing"

	validation "github.com/fgeck/go-register/internal/service/validation"
	"github.com/stretchr/testify/assert"
)

func TestValidateEmail(t *testing.T) {
	vs := validation.NewValidationService()

	tests := []struct {
		email    string
		expected string
	}{
		{"valid.email@example.com", ""},
		{"invalid-email", "invalid email format"},
		{"", "invalid email format"},
	}

	for _, test := range tests {
		err := vs.ValidateEmail(test.email)
		if test.expected == "" {
			assert.NoError(t, err, "expected no error for email: %s", test.email)
		} else {
			assert.EqualError(t, err, test.expected, "expected error for email: %s", test.email)
		}
	}
}

func TestValidatePassword(t *testing.T) {
	vs := validation.NewValidationService()

	tests := []struct {
		password string
		expected string
	}{
		{"SuperVal!d1@", ""},
		{"Valid1@", "password must be at least 8 characters long and include at least 1 uppercase letter, 1 lowercase letter, 1 number, and 1 special character"},
		{"short", "password must be at least 8 characters long and include at least 1 uppercase letter, 1 lowercase letter, 1 number, and 1 special character"},
		{"NoSpecialChar1", "password must be at least 8 characters long and include at least 1 uppercase letter, 1 lowercase letter, 1 number, and 1 special character"},
		{"nouppercase1@", "password must be at least 8 characters long and include at least 1 uppercase letter, 1 lowercase letter, 1 number, and 1 special character"},
	}

	for _, test := range tests {
		err := vs.ValidatePassword(test.password)
		if test.expected == "" {
			assert.NoError(t, err, "expected no error for password: %s", test.password)
		} else {
			assert.EqualError(t, err, test.expected, "expected error for password: %s", test.password)
		}
	}
}

func TestValidateUsername(t *testing.T) {
	vs := validation.NewValidationService()

	tests := []struct {
		username string
		expected string
	}{
		{"validUser", ""},
		{"val1dUs3r", ""},
		{"ab", "username must be at least 3 characters long"},
		{"invalid_user!", "username can only contain letters and numbers"},
	}

	for _, test := range tests {
		err := vs.ValidateUsername(test.username)
		if test.expected == "" {
			assert.NoError(t, err, "expected no error for username: %s", test.username)
		} else {
			assert.EqualError(t, err, test.expected, "expected error for username: %s", test.username)
		}
	}
}
