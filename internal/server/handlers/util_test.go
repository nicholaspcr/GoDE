package handlers

import (
	"errors"
	"testing"

	"github.com/nicholaspcr/GoDE/pkg/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestValidationErrorToStatus(t *testing.T) {
	t.Run("nil error returns nil", func(t *testing.T) {
		assert.Nil(t, ValidationErrorToStatus(nil))
	})

	t.Run("non-ValidationError returns simple InvalidArgument", func(t *testing.T) {
		err := ValidationErrorToStatus(errors.New("some error"))
		require.Error(t, err)

		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Equal(t, "some error", st.Message())
	})

	t.Run("ValidationError includes field violation details", func(t *testing.T) {
		valErr := &validation.ValidationError{
			Field:   "email",
			Message: "invalid format",
		}

		err := ValidationErrorToStatus(valErr)
		require.Error(t, err)

		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())

		// Check field violation details
		details := st.Details()
		require.Len(t, details, 1)
		br, ok := details[0].(*errdetails.BadRequest)
		require.True(t, ok)
		require.Len(t, br.FieldViolations, 1)
		assert.Equal(t, "email", br.FieldViolations[0].Field)
		assert.Equal(t, "invalid format", br.FieldViolations[0].Description)
	})
}

func TestNewValidationError(t *testing.T) {
	err := NewValidationError("username", "must be at least 3 characters")
	require.Error(t, err)

	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Equal(t, "validation failed", st.Message())

	details := st.Details()
	require.Len(t, details, 1)
	br, ok := details[0].(*errdetails.BadRequest)
	require.True(t, ok)
	require.Len(t, br.FieldViolations, 1)
	assert.Equal(t, "username", br.FieldViolations[0].Field)
	assert.Equal(t, "must be at least 3 characters", br.FieldViolations[0].Description)
}

func TestNewMultiFieldValidationError(t *testing.T) {
	t.Run("returns nil for empty violations", func(t *testing.T) {
		assert.Nil(t, NewMultiFieldValidationError(map[string]string{}))
	})

	t.Run("single violation", func(t *testing.T) {
		err := NewMultiFieldValidationError(map[string]string{
			"email": "invalid format",
		})
		require.Error(t, err)

		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())

		details := st.Details()
		require.Len(t, details, 1)
		br, ok := details[0].(*errdetails.BadRequest)
		require.True(t, ok)
		require.Len(t, br.FieldViolations, 1)
	})

	t.Run("multiple violations", func(t *testing.T) {
		err := NewMultiFieldValidationError(map[string]string{
			"username": "too short",
			"email":    "invalid format",
		})
		require.Error(t, err)

		st, ok := status.FromError(err)
		require.True(t, ok)

		details := st.Details()
		require.Len(t, details, 1)
		br, ok := details[0].(*errdetails.BadRequest)
		require.True(t, ok)
		assert.Len(t, br.FieldViolations, 2)
	})
}
