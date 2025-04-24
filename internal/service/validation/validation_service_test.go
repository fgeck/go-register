//go:build unittest

package validation_test

import (
	"testing"

	userfacing_errors "github.com/fgeck/go-register/internal/service/errors"
	validation "github.com/fgeck/go-register/internal/service/validation"
	"github.com/stretchr/testify/assert"
)

func TestValidateEmail(t *testing.T) {
	vs := validation.NewValidationService()

	tests := []struct {
		email    string
		expected *userfacing_errors.UserFacingError
	}{
		{"valid.email@example.com", nil},
		{"invalid-email", userfacing_errors.NewUserFacing("invalid email format", 400)},
		{"", userfacing_errors.NewUserFacing("invalid email format", 400)},
	}

	for _, test := range tests {
		err := vs.ValidateEmail(test.email)
		if test.expected == nil {
			assert.NoError(t, err, "expected no error for email: %s", test.email)
		} else {
			ufe, ok := err.(*userfacing_errors.UserFacingError)
			assert.True(t, ok, "expected a UserFacingError for email: %s", test.email)
			assert.Equal(t, test.expected.Message, ufe.Message, "unexpected error message for email: %s", test.email)
			assert.Equal(t, test.expected.Code, ufe.Code, "unexpected error code for email: %s", test.email)
		}
	}
}

func TestValidatePassword(t *testing.T) {
	vs := validation.NewValidationService()

	tests := []struct {
		password string
		expected *userfacing_errors.UserFacingError
	}{
		{"SuperVal!d1@", nil},
		{"Valid1@", userfacing_errors.NewUserFacing("password must be at least 8 characters long and include at least 1 uppercase letter, 1 lowercase letter, 1 number, and 1 special character", 400)},
		{"short", userfacing_errors.NewUserFacing("password must be at least 8 characters long and include at least 1 uppercase letter, 1 lowercase letter, 1 number, and 1 special character", 400)},
		{"NoSpecialChar1", userfacing_errors.NewUserFacing("password must be at least 8 characters long and include at least 1 uppercase letter, 1 lowercase letter, 1 number, and 1 special character", 400)},
		{"nouppercase1@", userfacing_errors.NewUserFacing("password must be at least 8 characters long and include at least 1 uppercase letter, 1 lowercase letter, 1 number, and 1 special character", 400)},
	}

	for _, test := range tests {
		err := vs.ValidatePassword(test.password)
		if test.expected == nil {
			assert.NoError(t, err, "expected no error for password: %s", test.password)
		} else {
			ufe, ok := err.(*userfacing_errors.UserFacingError)
			assert.True(t, ok, "expected a UserFacingError for password: %s", test.password)
			assert.Equal(t, test.expected.Message, ufe.Message, "unexpected error message for password: %s", test.password)
			assert.Equal(t, test.expected.Code, ufe.Code, "unexpected error code for password: %s", test.password)
		}
	}
}

func TestValidateUsername(t *testing.T) {
	vs := validation.NewValidationService()

	tests := []struct {
		username string
		expected *userfacing_errors.UserFacingError
	}{
		{"validUser", nil},
		{"val1dUs3r", nil},
		{"ab", userfacing_errors.NewUserFacing("username must be at least 3 characters long", 400)},
		{"invalid_user!", userfacing_errors.NewUserFacing("username can only contain letters and numbers", 400)},
	}

	for _, test := range tests {
		err := vs.ValidateUsername(test.username)
		if test.expected == nil {
			assert.NoError(t, err, "expected no error for username: %s", test.username)
		} else {
			ufe, ok := err.(*userfacing_errors.UserFacingError)
			assert.True(t, ok, "expected a UserFacingError for username: %s", test.username)
			assert.Equal(t, test.expected.Message, ufe.Message, "unexpected error message for username: %s", test.username)
			assert.Equal(t, test.expected.Code, ufe.Code, "unexpected error code for username: %s", test.username)
		}
	}
}
