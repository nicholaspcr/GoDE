package wfg

import "github.com/nicholaspcr/gde3/pkg/problems/models"

var WFG4 = models.ProblemFn{
	Fn: func(e *models.Elem, M int) error {
		n_var := len(e.X)
		n_obj := M
		k := 2 * (n_obj - 1)

		xu := arrange(1, n_var+1, 1)
		for i := range xu {
			xu[i] *= 2
		}

		var y []float64
		for i := 0; i < n_var; i++ {
			y = append(y, e.X[i]/xu[i])
		}

		y = wfg1_t1(y, n_var, k)
		y = wfg2_t2(y, n_var, k)
		y = wfg2_t3(y, n_obj, n_var, k)
		y = _post(y, _ones(n_obj-1)) // post

		var h []float64
		for m := 0; m < n_obj; m++ {
			h = append(h, _shape_linear(y[:len(y)-1], m+1))
		}

		S := arrange(2, 2*n_obj+1, 2)
		newObjs := _calculate(y, S, h)

		e.Objs = make([]float64, len(newObjs))
		copy(e.Objs, newObjs)
		return nil
	},
	Name: "wfg4",
}

// ---------------------------------------------------------------------------------------------------------
// wfg4 -> t1-t2 implementations
// ---------------------------------------------------------------------------------------------------------

// wfg4_t1 implementation
func wfg4_t1(X []float64, n, k int) []float64 {
	x := make([]float64, len(X[:k]))
	copy(x, X[:k])
	return _transformation_shift_multi_modal(X, 30.0, 10.0, 0.35)
}

// wfg4_t2 implementation
func wfg4_t2(X []float64, m, n, k int) []float64 {
	x := make([]float64, len(X))
	copy(x, X)

	ind_r_sum := k + (n-k)/2
	gap := k / (m - 1)

	var t []float64
	for i := 1; i < m; i++ {
		t = append(t, _reduction_weighted_sum_uniform(x[(m-1)*gap:(m*gap)]))
	}
	t = append(t, _reduction_weighted_sum_uniform(x[k:ind_r_sum]))

	return t
}
