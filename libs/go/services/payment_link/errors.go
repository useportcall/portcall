package payment_link

import "fmt"

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string { return e.Message }

func NewValidationError(format string, args ...any) *ValidationError {
	return &ValidationError{Message: fmt.Sprintf(format, args...)}
}
