package so

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"time"
)

// DE -> Differential Evolution function
func DE(
	NP, dim, gen, maxInst int,
	lower, upper, CR, F, P float64,
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
	path = checkFullPath(path, equation.fileName, variant.funcName, dim, P)
	f, err := os.Create(path)
	checkError(err)

	// TODO: Rename to reflect what it writes, X[...] points of the current best of population
	var extraDataPath string = dbPath + "/extra"
	extraDataPath = checkFullPath(extraDataPath, equation.fileName, variant.funcName, dim, P)
	fExtraData, err := os.Create(extraDataPath)
	checkError(err)

	var allPointsPath string = dbPath + "/allPoints"
	var fAllPoints *os.File
	if dim == 2 {
		allPointsPath = checkFullPath(allPointsPath, equation.fileName, variant.funcName, dim, P)
		var err error
		fAllPoints, err = os.Create(allPointsPath)
		checkError(err)
		defer fAllPoints.Close()
	}

	defer f.Close()
	defer fExtraData.Close()

	writeFileHeader(f, fExtraData, fAllPoints, gen, dim, NP)

	for currInst := 1; currInst <= maxInst; currInst++ {
		fmt.Fprintf(f, "%-30d\t", currInst)
		population := generatePopulation(NP, dim, lower, upper, equation.calcFunc)

		for currGen := 0; currGen < gen; currGen++ {
			writeFileBody(
				f, fExtraData, fAllPoints,
				currInst, currGen, gen,
				dim, NP,
				population,
			)

			// trial population vector
			trial := elemArrCopy(population)
			for i := 0; i < len(population); i++ {
				mutant, err := variant.makeMutant(population, F, P, i, dim)
				checkError(err)

				//the experimental
				index := rand.Int() % dim // selecting random index
				for j := 1; j <= dim; j++ {
					changeProbability := rand.Float64()
					if changeProbability < CR || j == dim {
						trial[i].X[index] = mutant.X[index]
					}
					index = (index + 1) % dim
				}

				for j := 0; j < dim; j++ {
					if trial[i].X[j] < lower {
						trial[i].X[j] = lower
					}
					if trial[i].X[j] > upper {
						trial[i].X[j] = upper
					}
				}
				trial[i].fit = equation.calcFunc(trial[i].X)

				// the crossover
				if trial[i].fit < population[i].fit {
					population[i] = trial[i].makeCopy()
				}
			}

			// sorting population
			sort.Sort(byFit(population))
		}
	}
}

func checkFilePath(filePath string) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		os.Mkdir(filePath, os.ModePerm)
	}
}

func checkFullPath(basePath, equationName, variantName string, dim int, P float64) string {
	var ret string = basePath
	checkFilePath(ret)
	ret += "/" + equationName
	checkFilePath(ret)
	ret += "/" + "dim-" + strconv.Itoa(dim)
	checkFilePath(ret)
	if variantName == "pbest" {
		ret += "/" + variantName
		checkFilePath(ret)
		s := fmt.Sprintf("%.2f", P)
		ret += "/p-" + s + ".csv"
	} else {
		ret += "/" + variantName + ".csv"
	}
	return ret
}

func writeFileHeader(f, fExtraData, fAllPoints *os.File, gen, dim, NP int) {
	// Writing generation numbers for easier graphs
	for genIndex := 0; genIndex <= gen; genIndex++ {
		if genIndex == gen {
			fmt.Fprintf(f, "%20d\n", genIndex)
		} else {
			fmt.Fprintf(f, "%20d\t", genIndex)
		}
	}

	//Header of the extra Data
	fmt.Fprintf(fExtraData, "inst\tgen\t")
	for i := 0; i < dim; i++ {
		fmt.Fprintf(fExtraData, "X[%d]\t", i)
	}
	fmt.Fprintf(fExtraData, "F(X[...])\n")

	// writes the X and Y of all Elements in population
	if dim == 2 {
		fmt.Fprintf(fAllPoints, "instance\tgeneration\tF[i]\t")
		for i := 0; i < NP; i++ {
			if i == NP-1 {
				fmt.Fprintf(fAllPoints, "Elem-%d\n", i+1)
			} else {
				fmt.Fprintf(fAllPoints, "Elem-%d\t", i+1)
			}

		}
	}
}

func writeFileBody(f, fExtraData, fAllPoints *os.File,
	currInst, currGen, gen int,
	DIM, NP int,
	pop []Elem) {
	//Begin --- Writing on f
	if currGen == gen-1 {
		fmt.Fprintf(f, "%20.3f\n", pop[0].fit)
	} else {
		fmt.Fprintf(f, "%20.3f\t", pop[0].fit)
	}
	//Finish --- Writing on f

	//Begin --- Writing on fExtraData
	fmt.Fprintf(fExtraData, "%5d\t%5d\t", currInst, currGen+1)
	for i := 0; i < DIM; i++ {
		fmt.Fprintf(fExtraData, "%20.3f\t", pop[0].X[i])
	}
	fmt.Fprintf(fExtraData, "%20.3f\n", pop[0].fit)
	// Finish --- Writing on fExtraData

	//Begin --- Writing on f
	if DIM == 2 {
		for i := 0; i < DIM; i++ {
			fmt.Fprintf(fAllPoints, "%5d\t%5d\t%5d\t", currInst, currGen, i)
			for j := 0; j < NP; j++ {
				if j == NP-1 {
					fmt.Fprintf(fAllPoints, "%20.3f\n", pop[j].X[i])
				} else {
					fmt.Fprintf(fAllPoints, "%20.3f\t", pop[j].X[i])
				}

			}
		}
	}
	//Finish --- Writing on f
}

func checkError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

// generates a group of points with random values
func generatePopulation(sz, dim int, lower, upper float64, calcFit func(x []float64) float64) []Elem {
	ret := make([]Elem, sz)
	rang := upper - lower // range between floor and ceiling
	for i := 0; i < sz; i++ {
		x := make([]float64, dim)
		for j := 0; j < dim; j++ {
			x[j] = rand.Float64()*rang + lower // value varies within [lower,upper]
		}
		ret[i] = Elem{
			X:   x,
			fit: calcFit(x),
		}
	}
	sort.Sort(byFit(ret))
	return ret
}
