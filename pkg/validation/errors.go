package validation

import (
	"errors"
	"fmt"
)

var (
	ErrEmptyField       = errors.New("field cannot be empty")
	ErrInvalidFormat    = errors.New("invalid format")
	ErrOutOfRange       = errors.New("value out of valid range")
	ErrInvalidEmail     = errors.New("invalid email format")
	ErrPasswordTooShort = errors.New("password must be at least 8 characters")
	ErrPasswordTooLong  = errors.New("password exceeds maximum length")
	ErrInvalidUsername  = errors.New("username contains invalid characters")
)

// ValidationError represents a validation error with context.
type ValidationError struct {
	Field   string
	Value   interface{}
	Message string
	Err     error
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed for field '%s': %s", e.Field, e.Message)
}

func (e *ValidationError) Unwrap() error {
	return e.Err
}

// NewValidationError creates a new validation error.
func NewValidationError(field string, value interface{}, err error, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Value:   value,
		Message: message,
		Err:     err,
	}
}
