package logger

import "fmt"

type LuminoError struct {
	Code    int
	Message string
}

func (e *LuminoError) Error() string {
	return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}

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
