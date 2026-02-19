package wfg

import (
	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/problems"
)

type wfg4 struct{}

// Wfg4 returns the WFG4 test problem, a many-objective benchmark with multi-modality.
// Objectives: m (configurable)
func Wfg4() problems.Interface {
	return &wfg4{}
}

func (w *wfg4) Name() string {
	return "wfg4"
}

func (w *wfg4) Evaluate(e *models.Vector, m int) error {
	n_var := len(e.Elements)
	n_obj := m
	k := 2 * (n_obj - 1)

	var y []float64
	xu := arange(2, 2*n_var+1, 2)

	for i := range n_var {
		y = append(y, e.Elements[i]/xu[i])
	}

	y = wfg4_t1(y, n_var, k)
	y = wfg4_t2(y, n_obj, k)

	y = post(y, ones(n_obj-1)) // post

	var h []float64
	for m := range n_obj {
		h = append(h, shapeConcave(y[:len(y)-1], m+1))
	}

	s := arange(2, 2*n_obj+1, 2)
	newObjs := calculate(y, s, h)

	e.Objectives = make([]float64, len(newObjs))
	copy(e.Objectives, newObjs)
	return nil
}

// ----------------------------------------------------------------------------
// wfg4 -> t implementations
// ----------------------------------------------------------------------------

// wfg4_t1 implementation
func wfg4_t1(X []float64, n, k int) []float64 {
	//nolint:prealloc // Dynamic slice growth is intentional for clarity
	var ret []float64
	for _, x := range X {
		ret = append(
			ret,
			transformationShiftMultiModal(x, 30.0, 10.0, 0.35),
		)
	}
	return ret
}

// wfg4_t2 implementation
func wfg4_t2(X []float64, m, k int) []float64 {
	x := make([]float64, len(X))
	copy(x, X)

	gap := k / (m - 1)

	var t []float64
	for i := 1; i < m; i++ {
		t = append(t, reductionWeightedSumUniform(x[(i-1)*gap:(i*gap)]))
	}
	t = append(t, reductionWeightedSumUniform(x[k:]))

	return t
}
