//Package errs defines and organizes error types and functions
package errs

import "net/http"

type AppError struct {
	Code    int
	Message string
}

// AsMessage retrieves the message string from current AppError
func (e AppError) AsMessage() *AppError {
	return &AppError{
		Message: e.Message,
	}
}

// NewValidationError defines the parameters for an AppError that occurs when a request is invalid.
func NewValidationError(message string) *AppError {
	return &AppError{
		Message: message,
		Code:    http.StatusUnprocessableEntity,
	}
}

// NewUnexpectedError defines the parameters for an AppError that occurs when an unexpected error occurs.
func NewUnexpectedError(message string) *AppError {
	return &AppError{
		Message: message,
		Code:    http.StatusInternalServerError,
	}
}
