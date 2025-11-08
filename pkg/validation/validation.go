// Package validation provides comprehensive input validation utilities for the GoDE project.
//
// This package offers a collection of validation functions for common data types and formats,
// including user credentials, email addresses, numeric ranges, and differential evolution
// configuration parameters. All validation functions return structured errors that conform
// to gRPC status codes for consistent API error handling.
//
// Key Features:
//   - Generic type constraints for type-safe numeric validation
//   - Regex-based format validation for emails and usernames
//   - Comprehensive DE (Differential Evolution) configuration validation
//   - User credential validation (username, email, password)
//   - Structured error responses with field context
//
// Example usage:
//
//	// Validate user input
//	if err := validation.ValidateUsername("john_doe"); err != nil {
//	    return err // Returns structured ValidationError
//	}
//
//	// Validate numeric ranges
//	if err := validation.ValidateRange(populationSize, int64(4), int64(10000), "population_size"); err != nil {
//	    return err
//	}
//
//	// Validate DE configuration
//	if err := validation.ValidateDEConfig(deConfig); err != nil {
//	    return err // Returns gRPC-compatible error
//	}
package validation

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	// Email regex (RFC 5322 simplified)
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	// Username regex: alphanumeric, underscore, hyphen only
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
)

// ValidateRange checks if a numeric value is within the specified range (inclusive).
func ValidateRange[T int | int64 | float32 | float64](value T, minVal, maxVal T, fieldName string) error {
	if value < minVal || value > maxVal {
		return NewValidationError(
			fieldName,
			value,
			ErrOutOfRange,
			fmt.Sprintf("must be between %v and %v, got %v", minVal, maxVal, value),
		)
	}
	return nil
}

// ValidatePositive checks if a numeric value is positive (> 0).
func ValidatePositive[T int | int64 | float32 | float64](value T, fieldName string) error {
	if value <= 0 {
		return NewValidationError(
			fieldName,
			value,
			ErrOutOfRange,
			"must be positive",
		)
	}
	return nil
}

// ValidateNonEmpty checks if a string is non-empty after trimming.
func ValidateNonEmpty(value string, fieldName string) error {
	if strings.TrimSpace(value) == "" {
		return NewValidationError(fieldName, value, ErrEmptyField, "cannot be empty")
	}
	return nil
}

// ValidateStringLength checks if a string length is within bounds.
func ValidateStringLength(value string, minLen, maxLen int, fieldName string) error {
	length := len(value)
	if length < minLen || length > maxLen {
		return NewValidationError(
			fieldName,
			value,
			ErrOutOfRange,
			fmt.Sprintf("length must be between %d and %d characters, got %d", minLen, maxLen, length),
		)
	}
	return nil
}

// ValidateEmail validates email format.
func ValidateEmail(email string) error {
	if !emailRegex.MatchString(email) {
		return NewValidationError("email", email, ErrInvalidEmail, "invalid email format")
	}
	return nil
}

// ValidateUsername validates username format and length.
func ValidateUsername(username string) error {
	if err := ValidateNonEmpty(username, "username"); err != nil {
		return err
	}

	if err := ValidateStringLength(username, 3, 64, "username"); err != nil {
		return err
	}

	if !usernameRegex.MatchString(username) {
		return NewValidationError(
			"username",
			username,
			ErrInvalidUsername,
			"must contain only alphanumeric characters, underscores, and hyphens",
		)
	}

	return nil
}

// ValidatePassword validates password requirements.
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return NewValidationError("password", "***", ErrPasswordTooShort, "minimum 8 characters required")
	}

	// bcrypt has a maximum of 72 bytes
	if len(password) > 72 {
		return NewValidationError("password", "***", ErrPasswordTooLong, "maximum 72 characters allowed")
	}

	return nil
}
