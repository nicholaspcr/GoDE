package wfg

import (
	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/problems"
)

type wfg8 struct{}

func Wfg8() problems.Interface {
	return &wfg8{}
}

func (w *wfg8) Name() string {
	return "wfg8"
}

func (w *wfg8) Evaluate(e *models.Vector, M int) error {

	n_var := len(e.Elements)
	n_obj := M
	k := 2 * (n_obj - 1)

	var y []float64
	xu := arange(2, 2*n_var+1, 2)

	for i := 0; i < n_var; i++ {
		y = append(y, e.Elements[i]/xu[i])
	}

	t_temp := wfg8_t1(y, k, n_var)
	copy(y[k:n_var], t_temp) // transfers to these position of the y vector

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
// wfg8 -> t implementations
// ----------------------------------------------------------------------------

func wfg8_t1(X []float64, k, n int) []float64 {
	x := make([]float64, len(X))
	copy(x, X)
	var ret []float64
	for i := k; i < n; i++ {
		aux := _reduction_weighted_sum_uniform(x[:i])
		ret = append(
			ret,
			_transformation_param_dependent(x[i], aux, 0.98/49.98, 0.02, 50.0),
		)
	}
	return ret
}
