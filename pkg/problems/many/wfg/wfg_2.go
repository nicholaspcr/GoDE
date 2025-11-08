package wfg

import (
	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/problems"
)

type wfg2 struct{}

// Wfg2 returns the WFG2 test problem, a many-objective benchmark with disconnected Pareto regions.
// Objectives: m (configurable)
func Wfg2() problems.Interface {
	return &wfg2{}
}

func (w *wfg2) Name() string {
	return "wfg2"
}

func (w *wfg2) Evaluate(e *models.Vector, m int) error {
	n_var := len(e.Elements)
	n_obj := m
	k := 2 * (n_obj - 1)

	var y []float64
	xu := arange(2, 2*n_var+1, 2)

	for i := 0; i < n_var; i++ {
		y = append(y, e.Elements[i]/xu[i])
	}

	y = wfg1_t1(y, n_var, k)
	y = wfg2_t2(y, n_var, k)
	y = wfg2_t3(y, n_obj, n_var, k)
	// post section
	a := _ones(n_obj - 1)
	y = _post(y, a)

	var h []float64
	for m := 0; m < n_obj-1; m++ {
		h = append(h, _shape_convex(y[:len(y)-1], m+1))
	}
	h = append(h, _shape_disconnected(y[0], 1, 1, 5))

	s := arange(2, 2*n_obj+1, 2)
	newObjs := _calculate(y, s, h)

	e.Objectives = make([]float64, len(newObjs))
	copy(e.Objectives, newObjs)
	return nil
}

// ----------------------------------------------------------------------------
// wfg2 -> t2-t3 implementations
// ----------------------------------------------------------------------------

// wfg2_t2 implementation
func wfg2_t2(X []float64, n, k int) []float64 {
	x := make([]float64, len(X[:k]))
	copy(x, X[:k])

	l := n - k
	ind_non_sep := k + l/2

	i := k + 1
	for i <= ind_non_sep {
		head := k + 2*(i-k) - 2
		tail := k + 2*(i-k)

		x = append(x, _reduction_non_sep(X[head:tail], 2))
		i++
	}
	return x
}

// wfg2_t3 implementation
func wfg2_t3(X []float64, m, n, k int) []float64 {
	ind_r_sum := k + (n-k)/2
	gap := k / (m - 1)

	var t []float64
	for i := 1; i < m; i++ {
		t = append(t, _reduction_weighted_sum_uniform(X[((i-1)*gap):(i*gap)]))
	}
	t = append(t, _reduction_weighted_sum_uniform(X[k:ind_r_sum]))

	return t
}
