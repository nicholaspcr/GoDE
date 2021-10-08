package wfg

import "github.com/nicholaspcr/gde3/pkg/models"

type wfg6 struct{}

func Wfg6() models.Problem {
	return &wfg6{}
}

func (w *wfg6) Name() string {
	return "wfg6"
}

func (w *wfg6) Evaluate(e *models.Vector, M int) error {

	n_var := len(e.X)
	n_obj := M
	k := 2 * (n_obj - 1)

	var y []float64
	xu := arange(2, 2*n_var+1, 2)

	for i := 0; i < n_var; i++ {
		y = append(y, e.X[i]/xu[i])
	}

	y = wfg1_t1(y, n_var, k)
	y = wfg6_t2(y, n_obj, n_var, k)
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
// wfg6 -> t implementations
// ----------------------------------------------------------------------------

func wfg6_t2(X []float64, m, n, k int) []float64 {
	gap := k / (m - 1)
	var ret []float64
	for i := 1; i < m; i++ {
		ret = append(ret, _reduction_non_sep(X[(i-1)*gap:(i*gap)], gap))
	}
	ret = append(ret, _reduction_non_sep(X[k:], n-k))
	return ret
}
