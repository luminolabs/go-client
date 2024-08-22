package logger

import "fmt"

// LuminoError represents a custom error type for Lumino
type LuminoError struct {
	Code    int    // Error code
	Message string // Error message
}

// Error returns the string representation of the error
func (e *LuminoError) Error() string {
	return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}

// NewError creates a new LuminoError with the given code and message
func NewError(code int, message string) *LuminoError {
	return &LuminoError{
		Code:    code,
		Message: message,
	}
}

// Define common errors
var (
	ErrInvalidInput   = NewError(1, "Invalid input")
	ErrNetworkFailure = NewError(2, "Network failure")
	ErrUnauthorized   = NewError(3, "Unauthorized action")
	// Add more common errors as needed
)
