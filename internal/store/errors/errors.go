// Package errors defines sentinel errors for storage operations.
package errors

import "errors"

var (
	// ErrUnsupportedFieldMask indicates a field mask operation is not supported.
	ErrUnsupportedFieldMask = errors.New("unsupported field mask")
)
