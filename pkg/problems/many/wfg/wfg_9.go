package wfg

import (
	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/problems"
)

type wfg9 struct{}

// Wfg9 returns the WFG9 test problem, a many-objective benchmark with deceptive and non-separable properties.
// Objectives: m (configurable)
func Wfg9() problems.Interface {
	return &wfg9{}
}

func (w *wfg9) Name() string {
	return "wfg9"
}

func (w *wfg9) Evaluate(e *models.Vector, m int) error {
	n_var := len(e.Elements)
	n_obj := m
	k := 2 * (n_obj - 1)

	var y []float64
	xu := arange(2, 2*n_var+1, 2)

	for i := range n_var {
		y = append(y, e.Elements[i]/xu[i])
	}

	copy(
		y[:n_var-1],
		wfg9_t1(y, n_var),
	) // transfers to these position of the y vector
	y = wfg9_t2(y, n_var, k)
	y = wfg9_t3(y, n_obj, n_var, k)
	y = post(y, ones(n_obj-1)) // post

	var h []float64
	for m := range n_obj {
		h = append(h, shapeConcave(y[:len(y)-1], m+1))
	}

	s := arange(2, 2*n_obj+1, 2)
	newObjs := calculate(y, s, h)

	e.Objectives = make([]float64, len(newObjs))
	copy(e.Objectives, newObjs)
	return nil
}

// ----------------------------------------------------------------------------
// wfg8 -> t implementations
// ----------------------------------------------------------------------------

func wfg9_t1(X []float64, n int) []float64 {
	x := make([]float64, len(X))
	copy(x, X)
	var ret []float64
	for i := 0; i < n-1; i++ {
		aux := reductionWeightedSumUniform(x[i+1:])
		ret = append(
			ret,
			transformationParamDependent(x[i], aux, 0.98/49.98, 0.02, 50.0),
		)
	}
	return ret
}

func wfg9_t2(X []float64, n, k int) []float64 {
	x := make([]float64, len(X))
	copy(x, X)
	var a, b []float64

	for i := range k {
		a = append(a, transformationShiftDeceptive(x[i], 0.35, 0.001, 0.05))
	}
	for i := k; i < n; i++ {
		b = append(b, transformationShiftMultiModal(x[i], 30.0, 95.0, 0.35))
	}

	var ret []float64
	ret = append(ret, a...)
	ret = append(ret, b...)

	return ret
}

func wfg9_t3(X []float64, m, n, k int) []float64 {
	gap := k / (m - 1)
	var ret []float64

	for i := 1; i < m; i++ {
		ret = append(ret, reductionNonSep(X[(i-1)*gap:(i*gap)], gap))
	}
	ret = append(ret, reductionNonSep(X[k:], n-k))

	return ret
}
