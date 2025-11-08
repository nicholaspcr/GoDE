package wfg

import (
	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/problems"
)

type wfg9 struct{}

// Wfg9 returns the WFG9 test problem, a many-objective benchmark with deceptive and non-separable properties.
// Objectives: m (configurable)
func Wfg9() problems.Interface {
	return &wfg9{}
}

func (w *wfg9) Name() string {
	return "wfg9"
}

func (w *wfg9) Evaluate(e *models.Vector, m int) error {
	n_var := len(e.Elements)
	n_obj := m
	k := 2 * (n_obj - 1)

	var y []float64
	xu := arange(2, 2*n_var+1, 2)

	for i := 0; i < n_var; i++ {
		y = append(y, e.Elements[i]/xu[i])
	}

	copy(
		y[:n_var-1],
		wfg9_t1(y, n_var),
	) // transfers to these position of the y vector
	y = wfg9_t2(y, n_var, k)
	y = wfg9_t3(y, n_obj, n_var, k)
	y = _post(y, _ones(n_obj-1)) // post

	var h []float64
	for m := 0; m < n_obj; m++ {
		h = append(h, _shape_concave(y[:len(y)-1], m+1))
	}

	s := arange(2, 2*n_obj+1, 2)
	newObjs := _calculate(y, s, h)

	e.Objectives = make([]float64, len(newObjs))
	copy(e.Objectives, newObjs)
	return nil
}

// ----------------------------------------------------------------------------
// wfg8 -> t implementations
// ----------------------------------------------------------------------------

func wfg9_t1(X []float64, n int) []float64 {
	x := make([]float64, len(X))
	copy(x, X)
	var ret []float64
	for i := 0; i < n-1; i++ {
		aux := _reduction_weighted_sum_uniform(x[i+1:])
		ret = append(
			ret,
			_transformation_param_dependent(x[i], aux, 0.98/49.98, 0.02, 50.0),
		)
	}
	return ret
}

func wfg9_t2(X []float64, n, k int) []float64 {
	x := make([]float64, len(X))
	copy(x, X)
	var a, b []float64

	for i := 0; i < k; i++ {
		a = append(a, _transformation_shift_deceptive(x[i], 0.35, 0.001, 0.05))
	}
	for i := k; i < n; i++ {
		b = append(b, _transformation_shift_multi_modal(x[i], 30.0, 95.0, 0.35))
	}

	var ret []float64
	ret = append(ret, a...)
	ret = append(ret, b...)

	return ret
}

func wfg9_t3(X []float64, m, n, k int) []float64 {
	gap := k / (m - 1)
	var ret []float64

	for i := 1; i < m; i++ {
		ret = append(ret, _reduction_non_sep(X[(i-1)*gap:(i*gap)], gap))
	}
	ret = append(ret, _reduction_non_sep(X[k:], n-k))

	return ret
}
