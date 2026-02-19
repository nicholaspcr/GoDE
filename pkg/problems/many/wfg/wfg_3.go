package wfg

import (
	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/problems"
)

type wfg3 struct{}

// Wfg3 returns the WFG3 test problem, a many-objective benchmark with a linear degenerate Pareto front.
// Objectives: m (configurable)
func Wfg3() problems.Interface {
	return &wfg3{}
}

func (w *wfg3) Name() string {
	return "wfg3"
}

func (w *wfg3) Evaluate(e *models.Vector, m int) error {
	n_var := len(e.Elements)
	n_obj := m
	k := 2 * (n_obj - 1)

	var y []float64
	xu := arange(2, 2*n_var+1, 2)

	for i := range n_var {
		y = append(y, e.Elements[i]/xu[i])
	}

	y = wfg1_t1(y, n_var, k)
	y = wfg2_t2(y, n_var, k)
	y = wfg2_t3(y, n_obj, n_var, k)

	// post section
	a := ones(n_obj - 1)
	for i := 1; i < len(a); i++ {
		a[i] = 0
	}
	y = post(y, a)

	var h []float64
	for m := range n_obj {
		h = append(h, shapeLinear(y[:len(y)-1], m+1))
	}

	s := arange(2, 2*n_obj+1, 2)
	newObjs := calculate(y, s, h)

	e.Objectives = make([]float64, len(newObjs))
	copy(e.Objectives, newObjs)
	return nil
}
