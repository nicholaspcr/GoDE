package wfg

import (
	"github.com/nicholaspcr/GoDE/pkg/models"
)

type wfg4 struct{}

func Wfg4() models.Problem {
	return &wfg4{}
}

func (w *wfg4) Name() string {
	return "wfg4"
}

func (w *wfg4) Evaluate(e *models.Vector, M int) error {

	n_var := len(e.X)
	n_obj := M
	k := 2 * (n_obj - 1)

	var y []float64
	xu := arange(2, 2*n_var+1, 2)

	for i := 0; i < n_var; i++ {
		y = append(y, e.X[i]/xu[i])
	}

	y = wfg4_t1(y, n_var, k)
	y = wfg4_t2(y, n_obj, k)

	y = _post(y, _ones(n_obj-1)) // post

	var h []float64
	for m := 0; m < n_obj; m++ {
		h = append(h, _shape_concave(y[:len(y)-1], m+1))
	}

	S := arange(2, 2*n_obj+1, 2)
	newObjs := _calculate(y, S, h)

	e.Objs = make([]float64, len(newObjs))
	copy(e.Objs, newObjs)
	return nil
}

// ----------------------------------------------------------------------------
// wfg4 -> t implementations
// ----------------------------------------------------------------------------

// wfg4_t1 implementation
func wfg4_t1(X []float64, n, k int) []float64 {
	var ret []float64
	for _, x := range X {
		ret = append(
			ret,
			_transformation_shift_multi_modal(x, 30.0, 10.0, 0.35),
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
		t = append(t, _reduction_weighted_sum_uniform(x[(i-1)*gap:(i*gap)]))
	}
	t = append(t, _reduction_weighted_sum_uniform(x[k:]))

	return t
}
