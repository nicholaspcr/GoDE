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

	// ErrProgressNotSupported indicates progress tracking is not supported in this store.
	ErrProgressNotSupported = errors.New("progress tracking not supported in database store")

	// ErrCancellationNotSupported indicates cancellation is not supported in this store.
	ErrCancellationNotSupported = errors.New("cancellation not supported in database store")

	// ErrPubSubNotSupported indicates pub/sub is not supported in this store.
	ErrPubSubNotSupported = errors.New("pub/sub not supported in database store")

	// ErrUserNotFound indicates the requested user was not found.
	ErrUserNotFound = errors.New("user not found")
)
