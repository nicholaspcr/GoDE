// Package errors defines sentinel errors for storage operations.
package errors

import "errors"

var (
	// ErrUnsupportedFieldMask indicates a field mask operation is not supported.
	ErrUnsupportedFieldMask = errors.New("unsupported field mask")

	// ErrExecutionNotFound indicates the requested execution was not found.
	ErrExecutionNotFound = errors.New("execution not found")

	// ErrParetoSetNotFound indicates the requested pareto set was not found.
	ErrParetoSetNotFound = errors.New("pareto set not found")
)
