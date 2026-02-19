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

func ones(n int) []float64 {
	a := make([]float64, n)
	for i := range n {
		a[i] = 1
	}
	return a
}

func post(t, a []float64) []float64 {
	x := make([]float64, 0, len(t))
	lastIndex := len(t) - 1
	for i := range lastIndex {
		x = append(x, math.Max(t[lastIndex], a[i])*(t[i]-0.5)+0.5)
	}
	x = append(x, t[lastIndex])
	return x
}

func calculate(xVec, s, h []float64) []float64 {
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

// correctTo01 handles the values that are between 0 +- 1e-10 and 1 +- e1-10,
// replaces with a fixed value
// instead of leaving floating points
func correctTo01(x float64) float64 {
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

// transformationShiftLinear
func transformationShiftLinear(y, shift float64) float64 {
	return correctTo01(
		math.Abs(y-shift) / math.Abs(math.Floor(shift-y)+shift),
	)
}

func transformationShiftDeceptive(y, A, B, C float64) float64 {
	tmp1 := math.Floor(y-A+B) * (1.0 - C + (A-B)/B) / (A - B)
	tmp2 := math.Floor(A+B-y) * (1.0 - C + (1.0-A-B)/B) / (1.0 - A - B)
	ret := 1.0 + (math.Abs(y-A)-B)*(tmp1+tmp2+1.0/B)
	return correctTo01(ret)
}

func transformationShiftMultiModal(y, A, B, C float64) float64 {
	tmp1 := math.Abs(y-C) / (2.0 * (math.Floor(C-y) + C))
	tmp2 := (4.0*A + 2.0) * math.Pi * (0.5 - tmp1)
	ret := (1.0 + math.Cos(tmp2) + 4.0*B*math.Pow(tmp1, 2.0)) / (B + 2.0)
	return correctTo01(ret)
}

func transformationBiasFlat(y, a, b, c float64) float64 {
	ret := a + math.Min(
		0.0,
		math.Floor(y-b),
	)*(a*(b-y)/b) - math.Min(
		0,
		math.Floor(c-y),
	)*((1.0-a)*(y-c)/(1.0-c))
	return correctTo01(ret)
}

func transformationBiasPoly(y, alpha float64) float64 {
	return correctTo01(math.Pow(y, alpha))
}

func transformationParamDependent(y, yDeg, A, B, C float64) float64 {
	aux := A - (1.0-2.0*yDeg)*math.Abs(math.Floor(0.5-yDeg)+A)
	ret := math.Pow(y, B+(C-B)*aux)
	return correctTo01(ret)
}

func transformationParamDeceptive(y float64, A, B, C float64) float64 {
	tmp1 := math.Floor(y-A+B) * (1.0 - C + (A-B)/B) / (A - B)
	tmp2 := math.Floor(A+B-y) * (1.0 - C + (1.0-A-B)/B) / (1.0 - A - B)
	ret := 1.0 + (math.Abs(y-A)-B)*(tmp1+tmp2+1.0/B)
	return correctTo01(ret)
}

// ----------------------------------------------------------------------------
// REDUCTION
// ----------------------------------------------------------------------------

func reductionWeightedSum(y, w []float64) float64 {
	var internalProduct float64
	var wSum float64
	for i := range w {
		internalProduct += y[i] * w[i]
		wSum += w[i]
	}
	return correctTo01(internalProduct / wSum)
}

func reductionWeightedSumUniform(y []float64) float64 {
	var mean float64
	for _, v := range y {
		mean += v
	}
	mean /= float64(len(y))
	return correctTo01(mean)
}

func reductionNonSep(x []float64, A int) float64 {
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

	return correctTo01(num / denom)
}

// ----------------------------------------------------------------------------
// SHAPE
// ----------------------------------------------------------------------------

func shapeConcave(X []float64, m int) float64 {
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
	return correctTo01(ret)
}

func shapeConvex(X []float64, m int) float64 {
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
	return correctTo01(ret)
}

func shapeLinear(X []float64, m int) float64 {
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
	return correctTo01(ret)
}

func shapeMixed(X, A, alpha float64) float64 {
	aux := 2.0 * A * math.Pi
	ret := math.Pow(1.0-X-(math.Cos(aux*X+0.5*math.Pi)/aux), alpha)
	return correctTo01(ret)
}

func shapeDisconnected(X, alpha, beta, A float64) float64 {
	aux := math.Cos(A * math.Pi * math.Pow(X, beta))
	return correctTo01(1.0 - math.Pow(X, alpha)*aux*aux)
}
