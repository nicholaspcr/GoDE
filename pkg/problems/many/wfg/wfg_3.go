package wfg

import "github.com/nicholaspcr/gde3/pkg/problems/models"

var WFG3 = models.ProblemFn{
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

		// post section
		A := _ones(n_obj - 1)
		for i := 1; i < len(A); i++ {
			A[i] = 0
		}
		y = _post(y, A)

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
	Name: "wfg3",
}
