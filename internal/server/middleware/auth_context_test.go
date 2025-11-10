package middleware

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUsernameFromContext(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		expected string
	}{
		{
			name: "username present in context",
			ctx: func() context.Context {
				//nolint:staticcheck // SA1029: Using empty struct as context key is intentional for testing
				return context.WithValue(context.Background(), usernameCtxKey, "testuser")
			}(),
			expected: "testuser",
		},
		{
			name:     "username not present in context",
			ctx:      context.Background(),
			expected: "",
		},
		{
			name: "wrong type in context",
			ctx: func() context.Context {
				//nolint:staticcheck // SA1029: Using empty struct as context key is intentional for testing
				return context.WithValue(context.Background(), usernameCtxKey, 12345)
			}(),
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UsernameFromContext(tt.ctx)
			assert.Equal(t, tt.expected, result)
		})
	}
}
