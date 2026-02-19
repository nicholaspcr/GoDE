package wfg

import (
	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/problems"
)

type wfg6 struct{}

// Wfg6 returns the WFG6 test problem, a many-objective benchmark with non-separable reduction.
// Objectives: m (configurable)
func Wfg6() problems.Interface {
	return &wfg6{}
}

func (w *wfg6) Name() string {
	return "wfg6"
}

func (w *wfg6) Evaluate(e *models.Vector, m int) error {
	n_var := len(e.Elements)
	n_obj := m
	k := 2 * (n_obj - 1)

	var y []float64
	xu := arange(2, 2*n_var+1, 2)

	for i := range n_var {
		y = append(y, e.Elements[i]/xu[i])
	}

	y = wfg1_t1(y, n_var, k)
	y = wfg6_t2(y, n_obj, n_var, k)
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
// wfg6 -> t implementations
// ----------------------------------------------------------------------------

func wfg6_t2(X []float64, m, n, k int) []float64 {
	gap := k / (m - 1)
	var ret []float64
	for i := 1; i < m; i++ {
		ret = append(ret, reductionNonSep(X[(i-1)*gap:(i*gap)], gap))
	}
	ret = append(ret, reductionNonSep(X[k:], n-k))
	return ret
}
