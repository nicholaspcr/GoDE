package wfg

import (
	"fmt"
	"testing"

	"github.com/nicholaspcr/gde3/pkg/problems/models"
)

func TestWFG1FN(t *testing.T) {
	tests := []struct {
		name     string
		x        []float64
		expected []float64
	}{
		{
			name: "test_case_1",
			x: []float64{0.24199364597771478, 0.06294085809752699,
				0.682979237196795, 0.20919587856003843, 0.8615217135283674,
				0.7476546016437432, 0.9409038322828246, 0.1680378421996956,
				0.5659362315602098, 0.9162810921996075, 0.4917771593035209,
				0.9919917334762469, 0.8452736699652191, 0.2720135716900983,
				0.4772027893543616, 0.7957435210039454, 0.4802668984201683,
				0.6262800875490805, 0.29995487782600794, 0.24415475358707514,
				0.9175107784830833, 0.05072118152238865, 0.8066710784368301,
				0.8210562785104756},
			expected: []float64{2.62459718, 1.01342557, 1.23516052},
		},
		//		{
		//			name:     "test_case_2",
		//			x:        []float64{},
		//			expected: []float64{},
		//		},
		//		{
		//			name:     "test_case_3",
		//			x:        []float64{},
		//			expected: []float64{},
		//		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := models.Elem{
				X: tt.x,
			}
			err := WFG1.Fn(&e, len(tt.expected))

			if err != nil {
				t.Errorf("failed to run the WFG1 func")
			}

			// string representation of the array
			received, want := "", ""

			// rounds up to the 8th decimal case
			for _, obj := range e.Objs {
				received += fmt.Sprintf("%.8f ", obj)
			}

			// rounds up to the 8th decimal case
			for _, obj := range tt.expected {
				want += fmt.Sprintf("%.8f ", obj)
			}

			// checks the strings
			if received != want {
				t.Errorf(
					"WFG1 wrong objs. received %v, expected %v",
					received,
					want,
				)
			}
		})
	}
}
