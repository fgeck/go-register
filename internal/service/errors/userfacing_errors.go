package userfacing_errors

import "fmt"

type UserFacingError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (e *UserFacingError) Error() string {
	return fmt.Sprintf("UserFacingError: %s (Code: %d)", e.Message, e.Code)
}

func New(message string, code int) *UserFacingError {
	return &UserFacingError{
		Message: message,
		Code:    code,
	}
}
