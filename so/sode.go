package so

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"time"
)

// Run -> single objevtive Differential Evolution function that accepts flags
func Run(p Params) {
	allEqts := [...]Equation{Rastrigin}
	allVariants := [...]Variant{Rand1, Rand2}
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
	if p.DIM <= 0 {
		arrDim = append(arrDim, 2, 5, 10)
	} else {
		arrDim = append(arrDim, p.DIM)
	}

	var arrP []float64
	if p.P < 0.0 {
		arrP = append(arrP, 0.01, 0.05, 0.1)
	} else {
		arrP = append(arrP, p.P)
	}

	for i := 0; i < len(inputs); i++ {
		for j := 0; j < len(arrDim); j++ {
			for k := 0; k < len(arrP); k++ {
				DE(p, inputs[i].Eq, inputs[i].Var)
			}
		}
	}

	resultTime := time.Since(testTimer)
	fmt.Println(resultTime)
}

// DE -> Differential Evolution function
func DE(
	p Params,
	equation Equation,
	variant Variant,
) {
	rand.Seed(time.Now().UTC().UnixNano())

	// base path for the sode
	var dbPath string = os.Getenv("HOME") + "/.goDE"
	checkFilePath(dbPath)
	dbPath += "/sode"
	checkFilePath(dbPath)

	var path string = dbPath + "/convergence"
	path = checkUsedPaths(path, equation.fileName, variant.funcName, p.DIM, p.P)
	f, err := os.Create(path)
	checkError(err)

	// TODO: Rename to reflect what it writes, X[...] points of the current best of population
	var extraDataPath string = dbPath + "/extra"
	extraDataPath = checkUsedPaths(extraDataPath, equation.fileName, variant.funcName, p.DIM, p.P)
	fExtraData, err := os.Create(extraDataPath)
	checkError(err)

	var allPointsPath string = dbPath + "/allPoints"
	var fAllPoints *os.File
	if p.DIM == 2 {
		allPointsPath = checkUsedPaths(allPointsPath, equation.fileName, variant.funcName, p.DIM, p.P)
		var err error
		fAllPoints, err = os.Create(allPointsPath)
		checkError(err)
		defer fAllPoints.Close()
	}

	defer f.Close()
	defer fExtraData.Close()

	writeFileHeader(f, fExtraData, fAllPoints, p.GEN, p.DIM, p.NP)

	for currInst := 1; currInst <= p.EXECS; currInst++ {
		fmt.Fprintf(f, "%-30d\t", currInst)
		population := generatePopulation(p.NP, p.DIM, p.FLOOR, p.CEIL, equation.calcFunc)

		for currGen := 0; currGen < p.GEN; currGen++ {
			writeFileBody(
				f, fExtraData, fAllPoints,
				currInst, currGen, p.GEN,
				p.DIM, p.NP,
				population,
			)

			// trial population vector
			trial := elemArrCopy(population)
			for i := 0; i < len(population); i++ {
				mutant, err := variant.makeMutant(population, p.F, p.P, i, p.DIM)
				checkError(err)

				//the experimental
				index := rand.Int() % p.DIM // selecting random index
				for j := 1; j <= p.DIM; j++ {
					changeProbability := rand.Float64()
					if changeProbability < p.CR || j == p.DIM {
						trial[i].X[index] = mutant.X[index]
					}
					index = (index + 1) % p.DIM
				}

				for j := 0; j < p.DIM; j++ {
					if trial[i].X[j] < p.FLOOR {
						trial[i].X[j] = p.FLOOR
					}
					if trial[i].X[j] > p.CEIL {
						trial[i].X[j] = p.CEIL
					}
				}
				trial[i].fit = equation.calcFunc(trial[i].X)

				// the crossover
				if trial[i].fit < population[i].fit {
					population[i] = trial[i].Copy()
				}
			}

			// sorting population
			sort.Sort(byFit(population))
		}
	}
}
