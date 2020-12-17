package so

import (
	"flag"
	"fmt"
	"time"
)

// Run -> single objevtive Differential Evolution function that accepts flags
func Run() {
	lower := flag.Float64("floor", -5, "Lower bound of the range of values.\n")
	upper := flag.Float64("ceil", 5, "Upper bound of the range of values.\n")
	NP := flag.Int("NP", 100, "Size of the population.\n")
	dim := flag.Int("dim", 0, "Dimension of the function, number of variables.\nDefault is to apply dimensions [2, 5, 10].\n")
	CR := flag.Float64("CR", 0.9, "Crossover Factor.\n")
	F := flag.Float64("F", 0.5, "Factor of the mutant vector.\n")
	P := flag.Float64("P", -1, "Factor of the p-best variant.\n")
	gen := flag.Int("gen", 500, "Quantity of generations to be processed.\n")
	maxInst := flag.Int("maxInst", 30, "number of times the differential evolution will be runned.\n")
	outputDir := flag.String("outDir", "./output/soDE", "path to the folder in which the output will be stored\n")
	flag.Parse()

	allEqts := [...]Equation{rastrigin}
	allVariants := [...]Variant{rand1, rand2}
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
	if *dim <= 0 {
		arrDim = append(arrDim, 2, 5, 10)
	} else {
		arrDim = append(arrDim, *dim)
	}

	var arrP []float64
	if *P < 0.0 {
		arrP = append(arrP, 0.01, 0.05, 0.1)
	} else {
		arrP = append(arrP, *P)
	}

	for i := 0; i < len(inputs); i++ {
		for j := 0; j < len(arrDim); j++ {
			for k := 0; k < len(arrP); k++ {
				DE(
					*NP, arrDim[j], *gen, *maxInst,
					*lower, *upper, *CR, *F, arrP[k],
					*outputDir,
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
