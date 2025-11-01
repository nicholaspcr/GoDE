package wfg

import (
	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/problems"
)

type wfg5 struct{}

func Wfg5() problems.Interface {
	return &wfg5{}
}

func (w *wfg5) Name() string {
	return "wfg5"
}

func (w *wfg5) Evaluate(e *models.Vector, M int) error {
	n_var := len(e.Elements)
	n_obj := M
	k := 2 * (n_obj - 1)

	var y []float64
	xu := arange(2, 2*n_var+1, 2)

	for i := 0; i < len(e.Elements); i++ {
		y = append(y, e.Elements[i]/xu[i])
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

	e.Objectives = make([]float64, len(newObjs))
	copy(e.Objectives, newObjs)
	return nil
}

// ----------------------------------------------------------------------------
// wfg5 -> t implementations
// ----------------------------------------------------------------------------

// wfg5_t1 implementation
func wfg5_t1(X []float64) []float64 {
	var ret []float64
	for _, x := range X {
		ret = append(ret, _transformation_param_deceptive(x, 0.35, 0.001, 0.05))
	}
	return ret
}
