package dtlz

import (
	"fmt"
	"testing"

	"github.com/nicholaspcr/GoDE/pkg/models"
)

func TestDTLZ7FN(t *testing.T) {
	tests := []struct {
		ProblemName string
		x           []float64
		expected    []float64
	}{
		{
			ProblemName: "test_case_1",
			x: []float64{
				0.040971105531507235, 0.550373235584878,
				0.6817311625009819, 0.6274478938025135, 0.9234111071427142,
				0.02499901960750534, 0.136171616578574, 0.9084459589232222,
				0.21089363254881652, 0.08574450529306678, 0.20551052286248087,
				0.43442188671029464,
			},
			expected: []float64{0.04097111, 0.55037324, 17.32742807},
		},
	}

	for _, tt := range tests {
		t.Run(tt.ProblemName, func(t *testing.T) {
			e := &models.Vector{Elements: tt.x}
			err := Dtlz7().Evaluate(e, len(tt.expected))
			if err != nil {
				t.Errorf("failed to run the DTLZ7 func")
			}

			// string representation of the array
			received, want := "", ""

			// rounds up to the 7th decimal case
			for _, obj := range e.Objectives {
				received += fmt.Sprintf("%.8f ", obj)
			}

			// rounds up to the 7th decimal case
			for _, obj := range tt.expected {
				want += fmt.Sprintf("%.8f ", obj)
			}

			// checks the strings
			if received != want {
				t.Errorf(
					"DTLZ7 wrong objs. received %v, expected %v",
					received,
					want,
				)
			}
		})
	}
}
