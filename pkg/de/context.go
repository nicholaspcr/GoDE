package de

import "context"

type executionNumberKey struct{}

// WithContextExecutionNumber returns a context with the execution number.
func WithContextExecutionNumber(ctx context.Context, n int) context.Context {
	return context.WithValue(ctx, executionNumberKey{}, n)
}

// FromContextExecutionNumber returns the execution number from the context.
// Returns 0 if the execution number is not set.
func FromContextExecutionNumber(ctx context.Context) int {
	val := ctx.Value(executionNumberKey{})
	if val == nil {
		return 0
	}
	n, ok := val.(int)
	if !ok {
		return 0
	}
	return n
}
