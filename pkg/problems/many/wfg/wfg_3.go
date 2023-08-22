package wfg

import (
	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/problems"
)

type wfg3 struct{}

func Wfg3() problems.Interface {
	return &wfg3{}
}

func (w *wfg3) Name() string {
	return "wfg3"
}

func (w *wfg3) Evaluate(e *models.Vector, M int) error {
	n_var := len(e.Elements)
	n_obj := M
	k := 2 * (n_obj - 1)

	var y []float64
	xu := arange(2, 2*n_var+1, 2)

	for i := 0; i < n_var; i++ {
		y = append(y, e.Elements[i]/xu[i])
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

	S := arange(2, 2*n_obj+1, 2)
	newObjs := _calculate(y, S, h)

	e.Objectives = make([]float64, len(newObjs))
	copy(e.Objectives, newObjs)
	return nil
}
