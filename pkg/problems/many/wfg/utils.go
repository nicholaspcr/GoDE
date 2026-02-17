package wfg

import (
	"math"
)

// ----------------------------------------------------------------------------
// WFG init
// ----------------------------------------------------------------------------

func arange(start, end, steps int) []float64 {
	size := (end - start + steps - 1) / steps
	s := make([]float64, 0, size)
	for i := start; i < end; i += steps {
		s = append(s, float64(i))
	}
	return s
}

func _ones(n int) []float64 {
	a := make([]float64, n)
	for i := range n {
		a[i] = 1
	}
	return a
}

func _post(t, a []float64) []float64 {
	x := make([]float64, 0, len(t))
	lastIndex := len(t) - 1
	for i := range lastIndex {
		x = append(x, math.Max(t[lastIndex], a[i])*(t[i]-0.5)+0.5)
	}
	x = append(x, t[lastIndex])
	return x
}

func _calculate(xVec, s, h []float64) []float64 {
	x := make([]float64, 0, len(h))

	// debug printing
	// fmt.Println("xVec -> ", xVec)
	// fmt.Println("s -> ", s)
	// fmt.Println("h -> ", h)
	for i := range h {
		x = append(x, xVec[len(xVec)-1]+s[i]*h[i])
	}
	return x
}

// ----------------------------------------------------------------------------
// utils
// ----------------------------------------------------------------------------

// _correct_to_01 handles the values that are between 0 +- 1e-10 and 1 +- e1-10,
// replaces with a fixed value
// instead of leaving floating points
func _correct_to_01(x float64) float64 {
	epsilon := 1e-10
	if x < 0.0 && x >= 0-epsilon {
		x = 0
	}
	if x > 1 && x <= 1+epsilon {
		x = 1
	}
	return x
}

// ----------------------------------------------------------------------------
// transformations
// ----------------------------------------------------------------------------

// _transformation_shift_linear
func _transformation_shift_linear(y, shift float64) float64 {
	return _correct_to_01(
		math.Abs(y-shift) / math.Abs(math.Floor(shift-y)+shift),
	)
}

func _transformation_shift_deceptive(y, A, B, C float64) float64 {
	tmp1 := math.Floor(y-A+B) * (1.0 - C + (A-B)/B) / (A - B)
	tmp2 := math.Floor(A+B-y) * (1.0 - C + (1.0-A-B)/B) / (1.0 - A - B)
	ret := 1.0 + (math.Abs(y-A)-B)*(tmp1+tmp2+1.0/B)
	return _correct_to_01(ret)
}

func _transformation_shift_multi_modal(y, A, B, C float64) float64 {
	tmp1 := math.Abs(y-C) / (2.0 * (math.Floor(C-y) + C))
	tmp2 := (4.0*A + 2.0) * math.Pi * (0.5 - tmp1)
	ret := (1.0 + math.Cos(tmp2) + 4.0*B*math.Pow(tmp1, 2.0)) / (B + 2.0)
	return _correct_to_01(ret)
}

func _transformation_bias_flat(y, a, b, c float64) float64 {
	ret := a + math.Min(
		0.0,
		math.Floor(y-b),
	)*(a*(b-y)/b) - math.Min(
		0,
		math.Floor(c-y),
	)*((1.0-a)*(y-c)/(1.0-c))
	return _correct_to_01(ret)
}

func _transformation_bias_poly(y, alpha float64) float64 {
	return _correct_to_01(math.Pow(y, alpha))
}

func _transformation_param_dependent(y, y_deg, A, B, C float64) float64 {
	aux := A - (1.0-2.0*y_deg)*math.Abs(math.Floor(0.5-y_deg)+A)
	ret := math.Pow(y, B+(C-B)*aux)
	return _correct_to_01(ret)
}

func _transformation_param_deceptive(y float64, A, B, C float64) float64 {
	tmp1 := math.Floor(y-A+B) * (1.0 - C + (A-B)/B) / (A - B)
	tmp2 := math.Floor(A+B-y) * (1.0 - C + (1.0-A-B)/B) / (1.0 - A - B)
	ret := 1.0 + (math.Abs(y-A)-B)*(tmp1+tmp2+1.0/B)
	return _correct_to_01(ret)
}

// ----------------------------------------------------------------------------
// REDUCTION
// ----------------------------------------------------------------------------

func _reduction_weighted_sum(_y, _w []float64) float64 {
	var internal_product float64
	var w_sum float64
	for i := range _w {
		internal_product += _y[i] * _w[i]
		w_sum += _w[i]
	}
	return _correct_to_01(internal_product / w_sum)
}

func _reduction_weighted_sum_uniform(y []float64) float64 {
	var mean float64
	for _, v := range y {
		mean += v
	}
	mean /= float64(len(y))
	return _correct_to_01(mean)
}

func _reduction_non_sep(x []float64, A int) float64 {
	val := math.Ceil(float64(A) / 2.0)

	var num float64
	m := len(x)

	for i := range x {
		num += x[i]
		for k := 0; k < A-1; k++ {
			num += math.Abs(x[i] - x[(1+i+k)%m])
		}
	}

	denom := float64(m) * val * (1.0 + 2.0*float64(A) - 2.0*val) / float64(A)

	return _correct_to_01(num / denom)
}

// ----------------------------------------------------------------------------
// SHAPE
// ----------------------------------------------------------------------------

func _shape_concave(X []float64, m int) float64 {
	n := len(X)
	var ret = 1.0
	switch {
	case m == 1:
		for _, x := range X[:n] {
			ret *= math.Sin(0.5 * x * math.Pi)
		}
	case 1 < m && m <= n:
		for _, x := range X[:(n - m + 1)] {
			ret *= math.Sin(0.5 * x * math.Pi)
		}
		ret *= math.Cos(0.5 * X[n-m+1] * math.Pi)
	default:
		ret *= math.Cos(0.5 * X[0] * math.Pi)
	}
	return _correct_to_01(ret)
}

func _shape_convex(X []float64, m int) float64 {
	n := len(X)
	var ret = 1.0
	switch {
	case m == 1:
		for _, x := range X[:n] {
			ret *= 1.0 - math.Cos(0.5*x*math.Pi)
		}
	case m > 1 && m <= n:
		for _, x := range X[:n-m+1] {
			ret *= (1.0 - math.Cos(0.5*x*math.Pi))
		}
		ret *= (1.0 - math.Sin(0.5*X[n-m+1]*math.Pi))
	default:
		ret = 1.0 - math.Sin(0.5*X[0]*math.Pi)
	}
	return _correct_to_01(ret)
}

func _shape_linear(X []float64, m int) float64 {
	n := len(X)
	var ret = 1.0
	switch {
	case m == 1:
		// prod
		for _, v := range X {
			ret *= v
		}
	case m > 1 && m <= n:
		// prod
		for _, x := range X[:n-m+1] {
			ret *= x
		}
		ret *= 1.0 - X[n-m+1]
	default:
		ret = 1.0 - X[0]
	}
	return _correct_to_01(ret)
}

func _shape_mixed(X, A, alpha float64) float64 {
	aux := 2.0 * A * math.Pi
	ret := math.Pow(1.0-X-(math.Cos(aux*X+0.5*math.Pi)/aux), alpha)
	return _correct_to_01(ret)
}

func _shape_disconnected(X, alpha, beta, A float64) float64 {
	aux := math.Cos(A * math.Pi * math.Pow(X, beta))
	return _correct_to_01(1.0 - math.Pow(X, alpha)*aux*aux)
}
