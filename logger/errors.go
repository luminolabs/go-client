// Package logger error handling and custom error types
package logger

import "fmt"

// LuminoError defines a custom error type for the Lumino client.
// Includes an error code and message for better error handling
// and reporting.
type LuminoError struct {
	Code    int    // Error code
	Message string // Error message
}

// Error implements the error interface for LuminoError.
// Returns a formatted string containing both the error code
// and message.
func (e *LuminoError) Error() string {
	return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}

// NewError creates a new LuminoError instance with the specified
// code and message. Used to generate consistent error types
// throughout the application.
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
