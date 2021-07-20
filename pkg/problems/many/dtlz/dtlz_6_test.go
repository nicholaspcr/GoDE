package dtlz

import (
	"fmt"
	"testing"

	"github.com/nicholaspcr/gde3/pkg/problems/models"
)

func TestDTLZ6FN(t *testing.T) {
	tests := []struct {
		name     string
		x        []float64
		expected []float64
	}{
		{
			name: "test_case_1",
			x: []float64{0.040971105531507235, 0.550373235584878,
				0.6817311625009819, 0.6274478938025135, 0.9234111071427142,
				0.02499901960750534, 0.136171616578574, 0.9084459589232222,
				0.21089363254881652, 0.08574450529306678, 0.20551052286248087,
				0.43442188671029464},
			expected: []float64{6.42107287, 7.40537588, 0.63167125},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := models.Elem{
				X: tt.x,
			}
			err := DTLZ6.Fn(&e, len(tt.expected))

			if err != nil {
				t.Errorf("failed to run the DTLZ6 func")
			}

			// string representation of the array
			received, want := "", ""

			// rounds up to the 7th decimal case
			for _, obj := range e.Objs {
				received += fmt.Sprintf("%.8f ", obj)
			}

			// rounds up to the 7th decimal case
			for _, obj := range tt.expected {
				want += fmt.Sprintf("%.8f ", obj)
			}

			// checks the strings
			if received != want {
				t.Errorf(
					"DTLZ6 wrong objs. received %v, expected %v",
					received,
					want,
				)
			}
		})
	}
}
