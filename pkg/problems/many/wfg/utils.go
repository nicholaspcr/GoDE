package wfg

import "math"

// ---------------------------------------------------------------------------------------------------------
// utils
// ---------------------------------------------------------------------------------------------------------

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

func _absolutes(X []float64) []float64 {
	x := make([]float64, len(X))
	copy(x, X)

	for i := range x {
		x[i] = math.Abs(x[i])
	}
	return x
}

// ---------------------------------------------------------------------------------------------------------
// transformations
// ---------------------------------------------------------------------------------------------------------

func _shiftLinear(value, shift float64) float64 {
	if shift == 0.0 {
		shift = 0.35
	}
	return _correct_to_01(math.Abs(value-shift) / math.Abs(math.Floor(shift-value)+shift))
}

func _biasFlat(value, a, b, c float64) float64 {
	ret := math.Min(0.0, math.Floor(value-b))*(a*(b-value)/b) - math.Min(0, math.Floor(c-value)*(1.0-a)*(value-c)/(1.0-c))
	return _correct_to_01(ret)
}

func _biasPoly(value, alpha float64) float64 {
	return _correct_to_01(math.Pow(value, alpha))
}

// ---------------------------------------------------------------------------------------------------------
// WFG init
// ---------------------------------------------------------------------------------------------------------

func arrange(start, end, steps int) []float64 {
	s := make([]float64, 0)
	for i := start; i < end; i += steps {
		s = append(s, float64(i))
	}
	return s
}

func createOnes(n int) []int {
	a := make([]int, 0)
	for i := 0; i < n; i++ {
		a = append(a, 1)
	}
	return a
}

// ---------------------------------------------------------------------------------------------------------
// REDUCTION
// ---------------------------------------------------------------------------------------------------------

func _reduction_weighted_sum(_y, _w []float64) float64 {
	var internal_product float64
	var w_sum float64
	for i := range _y {
		internal_product += _y[i] * _w[i]
		w_sum += _w[i]
	}
	return _correct_to_01(internal_product / w_sum)
}

// def _reduction_weighted_sum(y, w):
//     return correct_to_01(np.dot(y, w) / w.sum())

// def _reduction_weighted_sum_uniform(y):
//     return correct_to_01(y.mean(axis=1))

// def _reduction_non_sep(y, A):
//     n, m = y.shape
//     val = np.ceil(A / 2.0)

//     num = np.zeros(n)
//     for j in range(m):
//         num += y[:, j]
//         for k in range(A - 1):
//             num += np.fabs(y[:, j] - y[:, (1 + j + k) % m])

//     denom = m * val * (1.0 + 2.0 * A - 2 * val) / A

//     return correct_to_01(num / denom)

// ---------------------------------------------------------------------------------------------------------
// SHAPE
// ---------------------------------------------------------------------------------------------------------

func _shape_convex(X [][]float64, m int) []float64 {
	shape := len(X)
	var ret []float64
	if shape == 1 {
		for i := 0; i < shape; i++ {
			ret = append(ret, math.Sin(0.5*X[0][i]*math.Pi))
		}
	} else if m >= 1 && m <= shape {
		for i := 0; i < shape-m+1; i++ {
			ret = append(ret, math.Sin(0.5*X[0][i]*math.Pi))
		}
	} else {

	}
	return ret
}

// def _shape_concave(x, m):
//     M = x.shape[1]
//     if m == 1:
//         ret = np.prod(np.sin(0.5 * x[:, :M] * np.pi), axis=1)
//     elif 1 < m <= M:
//         ret = np.prod(np.sin(0.5 * x[:, :M - m + 1] * np.pi), axis=1)
//         ret *= np.cos(0.5 * x[:, M - m + 1] * np.pi)
//     else:
//         ret = np.cos(0.5 * x[:, 0] * np.pi)
//     return correct_to_01(ret)

// def _shape_convex(x, m):
//     M = x.shape[1]
//     if m == 1:
//         ret = np.prod(1.0 - np.cos(0.5 * x[:, :M] * np.pi), axis=1)
//     elif 1 < m <= M:
//         ret = np.prod(1.0 - np.cos(0.5 * x[:, :M - m + 1] * np.pi), axis=1)
//         ret *= 1.0 - np.sin(0.5 * x[:, M - m + 1] * np.pi)
//     else:
//         ret = 1.0 - np.sin(0.5 * x[:, 0] * np.pi)
//     return correct_to_01(ret)

// def _shape_linear(x, m):
//     M = x.shape[1]
//     if m == 1:
//         ret = np.prod(x, axis=1)
//     elif 1 < m <= M:
//         ret = np.prod(x[:, :M - m + 1], axis=1)
//         ret *= 1.0 - x[:, M - m + 1]
//     else:
//         ret = 1.0 - x[:, 0]
//     return correct_to_01(ret)

// def _shape_mixed(x, A=5.0, alpha=1.0):
//     aux = 2.0 * A * np.pi
//     ret = np.power(1.0 - x - (np.cos(aux * x + 0.5 * np.pi) / aux), alpha)
//     return correct_to_01(ret)

// def _shape_disconnected(x, alpha=1.0, beta=1.0, A=5.0):
//     aux = np.cos(A * np.pi * x ** beta)
//     return correct_to_01(1.0 - x ** alpha * aux ** 2)
