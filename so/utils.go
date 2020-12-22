package so

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"strconv"
)

func checkFilePath(filePath string) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		err = os.Mkdir(filePath, os.ModePerm)
		if err != nil {
			log.Fatalf("Error creating a folder in the path: %v", filePath)
		}
	}
}

func checkUsedPaths(basePath, equationName, variantName string, dim int, P float64) string {
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
func generatePopulation(sz, dim int, lower, upper float64, calcFit func(x []float64) float64) Elements {
	ret := make(Elements, sz)
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
