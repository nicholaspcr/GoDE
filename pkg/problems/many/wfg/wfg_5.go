package wfg

import "github.com/nicholaspcr/gde3/pkg/problems/models"

var WFG5 = models.ProblemFn{
	Fn: func(e *models.Elem, M int) error {
		n_var := len(e.X)
		n_obj := M
		k := 2 * (n_obj - 1)

		var y []float64
		xu := arange(2, 2*n_var+1, 2)

		for i := 0; i < n_var; i++ {
			y = append(y, e.X[i]/xu[i])
		}

		y = wfg5_t1(y)
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
	},
	Name: "wfg5",
}

// ---------------------------------------------------------------------------------------------------------
// wfg5 -> t implementations
// ---------------------------------------------------------------------------------------------------------

// wfg5_t1 implementation
func wfg5_t1(X []float64) []float64 {

	var ret []float64
	for _, x := range X {
		ret = append(ret, _transformation_param_deceptive(x, 0.35, 0.001, 0.05))
	}
	return ret
}
