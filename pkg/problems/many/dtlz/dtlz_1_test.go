package dtlz

import (
	"testing"

	"github.com/nicholaspcr/GoDE/pkg/problems/testutil"
)

func TestDTLZ1FN(t *testing.T) {
	tests := []struct {
		ProblemName string
		x           []float64
		expected    []float64
	}{
		{
			ProblemName: "test_case_1",
			x: []float64{
				0.040971105531507235, 0.550373235584878, 0.6817311625009819,
				0.6274478938025135, 0.9234111071427142, 0.02499901960750534,
				0.136171616578574, 0.9084459589232222, 0.21089363254881652,
				0.08574450529306678, 0.20551052286248087, 0.43442188671029464,
			},
			expected: []float64{9.86928923, 8.06270417, 419.7421470},
		},
	}

	for _, tt := range tests {
		t.Run(tt.ProblemName, func(t *testing.T) {
			testutil.AssertObjectivesEqual(t, Dtlz1(), tt.x, tt.expected, "DTLZ1")
		})
	}
}
