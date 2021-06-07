package wfg

import (
	"github.com/nicholaspcr/gde3/pkg/problems/models"
)

var WFG9 = models.ProblemFn{
	Fn: func(e *models.Elem, M int) error {
		n_var := len(e.X)
		n_obj := M
		k := 2 * (n_obj - 1)

		xu := arrange(2, 2*n_var+1, 2)

		var y []float64
		for i := 0; i < n_var; i++ {
			y = append(y, e.X[i]/xu[i])
		}

		copy(y[:n_var-1], wfg9_t1(y, n_var)) // transfers to these position of the y vector
		y = wfg9_t2(y, n_var, k)
		y = wfg9_t3(y, n_obj, n_var, k)
		y = _post(y, _ones(n_obj-1)) // post

		var h []float64
		for m := 0; m < n_obj; m++ {
			h = append(h, _shape_concave(y[:len(y)-1], m+1))
		}

		S := arrange(2, 2*n_obj+1, 2)
		newObjs := _calculate(y, S, h)

		e.Objs = make([]float64, len(newObjs))
		copy(e.Objs, newObjs)
		return nil
	},
	Name: "wfg9",
}

// ---------------------------------------------------------------------------------------------------------
// wfg8 -> t implementations
// ---------------------------------------------------------------------------------------------------------

func wfg9_t1(X []float64, n int) []float64 {
	x := make([]float64, len(X))
	copy(x, X)
	var ret []float64
	for i := 0; i < n-1; i++ {
		aux := _reduction_weighted_sum_uniform(x[i+1:])
		ret = append(ret, _transformation_param_dependent(x[i], aux, 0.98/49.98, 0.02, 50.0))
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

	// sums `a` and `b`
	sz := len(a)
	if len(b) > sz {
		sz = len(b)
	}

	var ret []float64

	for i := 0; i < sz; i++ {
		var valA, valB float64
		if i < len(a) {
			valA = a[i]
		}
		if i < len(b) {
			valB = b[i]
		}
		ret = append(ret, valA+valB)
	}

	return ret
}

func wfg9_t3(X []float64, m, n, k int) []float64 {
	gap := k / (m - 1)
	var ret []float64

	for i := 1; i < m; i++ {
		ret = append(ret, _reduction_non_sep(X[(m-1)*gap:(m*gap)], gap))
	}
	ret = append(ret, _reduction_non_sep(X[k:], n-k))

	return ret
}
