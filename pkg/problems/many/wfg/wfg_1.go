package wfg

import (
	"github.com/nicholaspcr/gde3/pkg/problems/models"
)

var WFG1 = models.ProblemFn{
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
		y = wfg1_t2(y, n_var, k)
		y = wfg1_t3(y, n_var)
		y = wfg1_t4(y, n_obj, n_var, k)

		// post
		y = _post(y, _ones(n_obj-1))

		var h []float64
		for m := 0; m < n_obj-1; m++ {
			h = append(h, _shape_convex(y[:len(y)-1], m+1))
		}
		h = append(h, _shape_mixed(y[0], 1, 5))

		S := arrange(2, 2*n_obj+1, 2)
		newObjs := _calculate(y, S, h)

		e.Objs = make([]float64, len(newObjs))
		copy(e.Objs, newObjs)
		return nil
	},
	Name: "wfg1",
}

// ---------------------------------------------------------------------------------------------------------
// wfg1 -> t1-t4 implementations
// ---------------------------------------------------------------------------------------------------------

// t1 implementations
func wfg1_t1(X []float64, n, k int) []float64 {
	x := make([]float64, len(X))
	copy(x, X)

	for i := k; i < n; i++ {
		x[i] = _shiftLinear(x[i], 0.0)
	}
	return x
}

// t2 implementation
func wfg1_t2(X []float64, n, k int) []float64 {
	x := make([]float64, len(X))
	copy(x, X)

	for i := k; i < n; i++ {
		x[i] = _biasFlat(x[i], 0.8, 0.75, 0.85)
	}
	return x
}

// t3 implementation
func wfg1_t3(X []float64, n int) []float64 {
	x := make([]float64, len(X))
	copy(x, X)

	for i := 0; i < n; i++ {
		x[i] = _biasPoly(x[i], 0.02)
	}

	return x
}

func wfg1_t4(X []float64, n_obj, n_var, k int) []float64 {
	x := make([]float64, len(X))
	copy(x, X)

	w := arrange(2, 2*n_var+1, 2)
	gap := k / (n_obj - 1)
	t := make([]float64, 0)

	for i := 1; i < n_obj; i++ {
		_y := x[(i-1)*gap : (i * gap)]
		_w := w[(i-1)*gap : (i * gap)]
		t = append(t, _reduction_weighted_sum(_y, _w))
	}
	t = append(t, _reduction_weighted_sum(x[k:n_var], w[k:n_var]))

	return t
}
