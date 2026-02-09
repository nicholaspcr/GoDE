package de

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithContextExecutionNumber(t *testing.T) {
	t.Run("stores execution number in context", func(t *testing.T) {
		ctx := WithContextExecutionNumber(context.Background(), 42)
		assert.Equal(t, 42, FromContextExecutionNumber(ctx))
	})

	t.Run("overwrites previous value", func(t *testing.T) {
		ctx := WithContextExecutionNumber(context.Background(), 1)
		ctx = WithContextExecutionNumber(ctx, 99)
		assert.Equal(t, 99, FromContextExecutionNumber(ctx))
	})
}

func TestFromContextExecutionNumber(t *testing.T) {
	t.Run("returns 0 when not set", func(t *testing.T) {
		assert.Equal(t, 0, FromContextExecutionNumber(context.Background()))
	})

	t.Run("returns stored value", func(t *testing.T) {
		ctx := WithContextExecutionNumber(context.Background(), 7)
		assert.Equal(t, 7, FromContextExecutionNumber(ctx))
	})

	t.Run("returns 0 for wrong type", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), executionNumberKey{}, "not-an-int")
		assert.Equal(t, 0, FromContextExecutionNumber(ctx))
	})
}
