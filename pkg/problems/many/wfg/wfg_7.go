package wfg

import (
	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/problems"
)

type wfg7 struct{}

func Wfg7() problems.Interface {
	return &wfg7{}
}

func (w *wfg7) Name() string {
	return "wfg7"
}

func (w *wfg7) Evaluate(e *models.Vector, M int) error {
	n_var := len(e.Elements)
	n_obj := M
	k := 2 * (n_obj - 1)

	var y []float64
	xu := arange(2, 2*n_var+1, 2)

	for i := 0; i < n_var; i++ {
		y = append(y, e.Elements[i]/xu[i])
	}

	y = wfg7_t1(y, k)
	y = wfg1_t1(y, n_var, k)
	y = wfg4_t2(y, n_obj, k)
	y = _post(y, _ones(n_obj-1)) // post

	var h []float64
	for m := 0; m < n_obj; m++ {
		h = append(h, _shape_concave(y[:len(y)-1], m+1))
	}

	S := arange(2, 2*n_obj+1, 2)
	newObjs := _calculate(y, S, h)

	e.Objectives = make([]float64, len(newObjs))
	copy(e.Objectives, newObjs)
	return nil
}

// ----------------------------------------------------------------------------
// wfg7 -> t implementations
// ----------------------------------------------------------------------------

func wfg7_t1(X []float64, k int) []float64 {
	x := make([]float64, len(X))
	copy(x, X)

	for i := 0; i < k; i++ {
		aux := _reduction_weighted_sum_uniform(x[i+1:])
		x[i] = _transformation_param_dependent(
			x[i],
			aux,
			0.98/49.98,
			0.02,
			50.0,
		)
	}
	return x
}
