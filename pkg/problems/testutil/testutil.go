// Package testutil provides testing utilities for problem evaluation tests.
package testutil

import (
	"fmt"
	"testing"

	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/problems"
)

// AssertObjectivesEqual evaluates a problem and asserts that the resulting objectives
// match the expected values (rounded to 7 decimal places).
// This is a convenience wrapper for AssertObjectivesEqualWithPrecision using 7 decimal places.
func AssertObjectivesEqual(t *testing.T, problem problems.Interface, input []float64, expected []float64, problemName string) {
	AssertObjectivesEqualWithPrecision(t, problem, input, expected, problemName, 7)
}

// AssertObjectivesEqualWithPrecision evaluates a problem and asserts that the resulting objectives
// match the expected values (rounded to the specified decimal precision).
func AssertObjectivesEqualWithPrecision(t *testing.T, problem problems.Interface, input []float64, expected []float64, problemName string, precision int) {
	t.Helper()

	vector := &models.Vector{Elements: input}
	err := problem.Evaluate(vector, len(expected))
	if err != nil {
		t.Errorf("failed to evaluate %s: %v", problemName, err)
		return
	}

	// Format both received and expected as strings with specified decimal precision
	formatStr := fmt.Sprintf("%%.%df ", precision)
	received, want := "", ""
	for _, obj := range vector.Objectives {
		received += fmt.Sprintf(formatStr, obj)
	}
	for _, obj := range expected {
		want += fmt.Sprintf(formatStr, obj)
	}

	// Compare string representations
	if received != want {
		t.Errorf(
			"%s wrong objectives.\nreceived: %v\nexpected: %v",
			problemName,
			received,
			want,
		)
	}
}
