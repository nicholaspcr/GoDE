package wfg

import (
	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/problems"
)

type wfg1 struct{}

func Wfg1() problems.Interface {
	return &wfg1{}
}

func (w *wfg1) Name() string {
	return "wfg1"
}

func (w *wfg1) Evaluate(e *models.Vector, M int) error {
	n_var := len(e.Elements)
	n_obj := M
	k := 2 * (n_obj - 1)

	var y []float64
	xu := arange(2, 2*n_var+1, 2)

	for i := 0; i < n_var; i++ {
		y = append(y, e.Elements[i]/xu[i])
	}

	y = wfg1_t1(y, n_var, k)
	y = wfg1_t2(y, n_var, k)
	y = wfg1_t3(y, n_var)
	y = wfg1_t4(y, n_obj, n_var, k)

	// post section
	A := _ones(n_obj - 1)
	y = _post(y, A)

	var h []float64
	for m := 0; m < n_obj-1; m++ {
		h = append(h, _shape_convex(y[:(len(y)-1)], m+1))
	}
	h = append(h, _shape_mixed(y[0], 5.0, 1.0))

	S := arange(2, 2*n_obj+1, 2)

	// fmt.Println(y, S, h)
	newObjs := _calculate(y, S, h)

	e.Objectives = make([]float64, len(newObjs))
	copy(e.Objectives, newObjs)
	return nil
}

// ---------------------------------------------------------------------------------------------------------
// wfg1 -> t implementations
// ---------------------------------------------------------------------------------------------------------

// t1 implementations
func wfg1_t1(X []float64, n, k int) []float64 {
	x := make([]float64, len(X))
	copy(x, X)

	for i := k; i < n; i++ {
		x[i] = _transformation_shift_linear(x[i], 0.35)
	}
	return x
}

// t2 implementation
func wfg1_t2(X []float64, n, k int) []float64 {
	x := make([]float64, len(X))
	copy(x, X)

	for i := k; i < n; i++ {
		x[i] = _transformation_bias_flat(x[i], 0.8, 0.75, 0.85)
	}
	return x
}

// t3 implementation
func wfg1_t3(X []float64, n int) []float64 {
	x := make([]float64, len(X))
	copy(x, X)

	for i := 0; i < n; i++ {
		x[i] = _transformation_bias_poly(x[i], 0.02)
	}

	return x
}

func wfg1_t4(X []float64, m, n, k int) []float64 {
	x := make([]float64, len(X))
	copy(x, X)

	w := arange(2, 2*n+1, 2)
	gap := k / (m - 1)
	var t []float64

	for i := 1; i < m; i++ {
		_y := x[(i-1)*gap : (i * gap)]
		_w := w[(i-1)*gap : (i * gap)]
		t = append(t, _reduction_weighted_sum(_y, _w))
	}
	t = append(t, _reduction_weighted_sum(X[k:n], w[k:n]))

	return t
}
