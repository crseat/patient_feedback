//Package errs defines and organizes error types and functions
package errs

type AppError struct {
	Code    int
	Message string
}

// NewAppError defines the parameters for an AppError that occurs when a request is invalid.
func NewAppError(message string) *AppError {
	return &AppError{
		Code:    1,
		Message: message,
	}
}
