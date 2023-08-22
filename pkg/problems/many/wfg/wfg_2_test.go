package wfg

import (
	"fmt"
	"testing"

	"github.com/nicholaspcr/GoDE/pkg/models"
)

func TestWFG2FN(t *testing.T) {
	tests := []struct {
		ProblemName string
		x           []float64
		expected    []float64
	}{
		{
			ProblemName: "test_case_1",
			x: []float64{
				0.24199364597771478, 0.06294085809752699,
				0.682979237196795, 0.20919587856003843, 0.8615217135283674,
				0.7476546016437432, 0.9409038322828246, 0.1680378421996956,
				0.5659362315602098, 0.9162810921996075, 0.4917771593035209,
				0.9919917334762469, 0.8452736699652191, 0.2720135716900983,
				0.4772027893543616, 0.7957435210039454, 0.4802668984201683,
				0.6262800875490805, 0.29995487782600794, 0.24415475358707514,
				0.9175107784830833, 0.05072118152238865, 0.8066710784368301,
				0.8210562785104756,
			},
			expected: []float64{0.64677672, 0.66722164, 6.55349009},
		},
		{
			ProblemName: "test_case_2",
			x: []float64{
				0.046812816038915586, 0.27965700782202974,
				0.755529270669409, 0.17423804874084414, 0.0601426129551884,
				0.5775324962743565, 0.10860410926044652, 0.5838059492224881,
				0.9223086754458868, 0.45421259093686484, 0.9437717291882289,
				0.5916451110680568, 0.5142652818480095, 0.27333526956383175,
				0.03112196519344138, 0.9981705782913868, 0.3677294655821704,
				0.4464982721068232, 0.4889880005853355, 0.3278044045109515,
				0.8551705425894127, 0.1877510146533291, 0.11597382226325655,
				0.59450013267575,
			},
			expected: []float64{0.66618332, 0.6756434, 6.51150653},
		},
		{
			ProblemName: "test_case_3",
			x: []float64{
				0.046812816038915586, 0.27965700782202974,
				0.755529270669409, 0.17423804874084414, 0.0601426129551884,
				0.5775324962743565, 0.10860410926044652, 0.5838059492224881,
				0.9223086754458868, 0.45421259093686484, 0.9437717291882289,
				0.5916451110680568, 0.5142652818480095, 0.27333526956383175,
				0.03112196519344138, 0.9981705782913868, 0.3677294655821704,
				0.4464982721068232, 0.4889880005853355, 0.3278044045109515,
				0.8551705425894127, 0.1877510146533291, 0.11597382226325655,
				0.59450013267575,
			},
			expected: []float64{0.66618332, 0.6756434, 6.51150653},
		},
	}

	for _, tt := range tests {
		t.Run(tt.ProblemName, func(t *testing.T) {
			e := *models.Vector{
				Elements: tt.x,
			}
			err := Wfg2().Evaluate(e, len(tt.expected))
			if err != nil {
				t.Errorf("failed to run the WFG2 func")
			}

			// string representation of the array
			received, want := "", ""

			// rounds up to the 8th decimal case
			for _, obj := range e.Objectives {
				received += fmt.Sprintf("%.8f ", obj)
			}

			// rounds up to the 8th decimal case
			for _, obj := range tt.expected {
				want += fmt.Sprintf("%.8f ", obj)
			}

			// checks the strings
			if received != want {
				t.Errorf(
					"WFG2 wrong objs. received %v, expected %v",
					received,
					want,
				)
			}
		})
	}
}
