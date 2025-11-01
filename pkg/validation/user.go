package validation

import (
	"strings"

	"github.com/nicholaspcr/GoDE/pkg/api/v1"
)

// ValidateUser validates all user fields.
func ValidateUser(user *api.User) error {
	if user == nil {
		return NewValidationError("user", nil, ErrEmptyField, "user object is nil")
	}

	if user.Ids == nil {
		return NewValidationError("user.ids", nil, ErrEmptyField, "user IDs are required")
	}

	// Validate username
	if err := ValidateUsername(user.Ids.Username); err != nil {
		return err
	}

	// Validate email
	email := strings.TrimSpace(strings.ToLower(user.Email))
	if err := ValidateNonEmpty(email, "email"); err != nil {
		return err
	}
	if err := ValidateEmail(email); err != nil {
		return err
	}

	// Validate password
	if err := ValidatePassword(user.Password); err != nil {
		return err
	}

	return nil
}
