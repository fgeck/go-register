package custom_errors

import "fmt"

type InternalError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (e *InternalError) Error() string {
	return fmt.Sprintf("InternalError: %s (Code: %d)", e.Message, e.Code)
}

func NewInternal(message string, code int) *InternalError {
	return &InternalError{
		Message: message,
		Code:    code,
	}
}
