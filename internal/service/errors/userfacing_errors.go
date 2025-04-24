package customErrors

import "fmt"

type UserFacingError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func NewUserFacing(message string, code int) *UserFacingError {
	return &UserFacingError{
		Message: message,
		Code:    code,
	}
}

func (e *UserFacingError) Error() string {
	return fmt.Sprintf("UserFacingError: %s (Code: %d)", e.Message, e.Code)
}
