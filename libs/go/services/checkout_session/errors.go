package checkout_session

import "fmt"

// ValidationError represents a client-facing validation error (e.g., bad input).
// HTTP handlers should map these to 400 responses.
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string { return e.Message }

// NewValidationError creates a new validation error.
func NewValidationError(format string, args ...any) *ValidationError {
	return &ValidationError{Message: fmt.Sprintf(format, args...)}
}
