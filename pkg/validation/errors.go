package validation

import (
	"errors"
	"fmt"
)

var (
	// ErrEmptyField indicates a required field was empty or whitespace-only.
	ErrEmptyField       = errors.New("field cannot be empty")
	// ErrInvalidFormat indicates a field value does not match the expected format.
	ErrInvalidFormat    = errors.New("invalid format")
	// ErrOutOfRange indicates a numeric value is outside the valid range.
	ErrOutOfRange       = errors.New("value out of valid range")
	// ErrInvalidEmail indicates an email address format is invalid.
	ErrInvalidEmail     = errors.New("invalid email format")
	// ErrPasswordTooShort indicates a password is shorter than the minimum length.
	ErrPasswordTooShort = errors.New("password must be at least 8 characters")
	// ErrPasswordTooLong indicates a password exceeds the maximum length.
	ErrPasswordTooLong  = errors.New("password exceeds maximum length")
	// ErrInvalidUsername indicates a username contains invalid characters.
	ErrInvalidUsername  = errors.New("username contains invalid characters")
)

// ValidationError represents a validation error with context.
type ValidationError struct {
	Value   interface{}
	Err     error
	Field   string
	Message string
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
