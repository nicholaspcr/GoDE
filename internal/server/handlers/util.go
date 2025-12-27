package handlers

import (
	"github.com/nicholaspcr/GoDE/pkg/validation"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ValidationErrorToStatus converts a validation error to a gRPC status with field violations.
// This provides structured error information that clients can parse programmatically.
func ValidationErrorToStatus(err error) error {
	if err == nil {
		return nil
	}

	// Check if it's a ValidationError
	valErr, ok := err.(*validation.ValidationError)
	if !ok {
		// For non-ValidationError, return a simple InvalidArgument status
		return status.Error(codes.InvalidArgument, err.Error())
	}

	// Create a BadRequest with field violations
	br := &errdetails.BadRequest{
		FieldViolations: []*errdetails.BadRequest_FieldViolation{
			{
				Field:       valErr.Field,
				Description: valErr.Message,
			},
		},
	}

	// Create the status with details
	st := status.New(codes.InvalidArgument, err.Error())
	stWithDetails, err := st.WithDetails(br)
	if err != nil {
		// If we can't add details, return the basic status
		return st.Err()
	}

	return stWithDetails.Err()
}

// NewValidationError creates a gRPC status error with field violation details.
func NewValidationError(field, description string) error {
	br := &errdetails.BadRequest{
		FieldViolations: []*errdetails.BadRequest_FieldViolation{
			{
				Field:       field,
				Description: description,
			},
		},
	}

	st := status.New(codes.InvalidArgument, "validation failed")
	stWithDetails, err := st.WithDetails(br)
	if err != nil {
		return st.Err()
	}

	return stWithDetails.Err()
}

// NewMultiFieldValidationError creates a gRPC status error with multiple field violations.
func NewMultiFieldValidationError(violations map[string]string) error {
	if len(violations) == 0 {
		return nil
	}

	fieldViolations := make([]*errdetails.BadRequest_FieldViolation, 0, len(violations))
	for field, desc := range violations {
		fieldViolations = append(fieldViolations, &errdetails.BadRequest_FieldViolation{
			Field:       field,
			Description: desc,
		})
	}

	br := &errdetails.BadRequest{
		FieldViolations: fieldViolations,
	}

	st := status.New(codes.InvalidArgument, "validation failed")
	stWithDetails, err := st.WithDetails(br)
	if err != nil {
		return st.Err()
	}

	return stWithDetails.Err()
}
