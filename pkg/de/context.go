package de

import "context"

type executionNumberKey struct{}

// WithContextExecutionNumber returns a context with the execution number.
func WithContextExecutionNumber(ctx context.Context, n int) context.Context {
	return context.WithValue(ctx, executionNumberKey{}, n)
}

// FromContextExecutionNumber returns the execution number from the context.
func FromContextExecutionNumber(ctx context.Context) int {
	return ctx.Value(executionNumberKey{}).(int)
}
