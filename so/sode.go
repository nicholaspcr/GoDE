package so

import (
	"fmt"
	"time"
)

// Run -> single objevtive Differential Evolution function that accepts flags
func Run(np, dim, gen, maxInst int, lower, upper, CR, F, P float64) {
	allEqts := [...]Equation{Rastrigin}
	allVariants := [...]Variant{Rand1, rand2}
	// , best1, best2, currToBestv1, currToBestv2, pBest

	testTimer := time.Now()
	var inputs []Input

	for i := 0; i < len(allEqts); i++ {
		for j := 0; j < len(allVariants); j++ {
			inputs = append(inputs, Input{
				Eq:  allEqts[i],
				Var: allVariants[j],
			})
		}
	}

	// Adding standard dimension values in case is not specified
	var arrDim []int
	if dim <= 0 {
		arrDim = append(arrDim, 2, 5, 10)
	} else {
		arrDim = append(arrDim, dim)
	}

	var arrP []float64
	if P < 0.0 {
		arrP = append(arrP, 0.01, 0.05, 0.1)
	} else {
		arrP = append(arrP, P)
	}

	for i := 0; i < len(inputs); i++ {
		for j := 0; j < len(arrDim); j++ {
			for k := 0; k < len(arrP); k++ {
				DE(
					np, arrDim[j], gen, maxInst,
					lower, upper, CR, F, arrP[k],
					inputs[i].Eq, inputs[i].Var,
				)
			}
		}
	}

	resultTime := time.Since(testTimer)
	fmt.Println(resultTime)
}

// Input -> input used in main
type Input struct {
	Eq  Equation
	Var Variant
}
