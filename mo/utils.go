package mo

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
)

// GetProblemByName -> returns the problem function
func GetProblemByName(name string) ProblemFn {
	switch name {
	case "ZDT1":
		return ZDT1
	case "ZDT2":
		return ZDT2
	case "ZDT3":
		return ZDT3
	case "ZDT4":
		return ZDT4
	case "ZDT6":
		return ZDT6
	case "VNT1":
		return VNT1
	case "DTLZ1":
		return DTLZ1
	case "DTLZ2":
		return DTLZ2
	case "DTLZ3":
		return DTLZ3
	case "DTLZ4":
		return DTLZ4
	case "DTLZ5":
		return DTLZ5
	case "DTLZ6":
		return DTLZ6
	case "DTLZ7":
		return DTLZ7
	default:
		return nil
	}
}

// GetVariantByName -> Returns the variant function
func GetVariantByName(name string) VariantFn {
	name = strings.ToLower(name)
	switch name {
	case "rand1":
		return Rand1
	default:
		return nil
	}
}

func generatePopulation(p Params) Elements {
	ret := make(Elements, p.NP)
	constant := p.CEIL - p.FLOOR // range between floor and ceiling
	for i := 0; i < p.NP; i++ {
		ret[i].X = make([]float64, p.DIM)

		for j := 0; j < p.DIM; j++ {
			ret[i].X[j] = rand.Float64()*constant + p.FLOOR // value varies within [ceil,upper]
		}
	}
	return ret
}

// generates random indices in the int slice, r -> it's a pointer
func generateIndices(startInd, NP int, r []int) error {
	if len(r) > NP {
		return errors.New("insufficient elements in population to generate random indices")
	}
	for i := startInd; i < len(r); i++ {
		for done := false; !done; {
			r[i] = rand.Int() % NP
			done = true
			for j := 0; j < i; j++ {
				done = done && r[j] != r[i]
			}
		}
	}
	return nil
}

func checkFilePath(filePath string) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		err = os.Mkdir(filePath, os.ModePerm)
		if err != nil {
			log.Fatalf("error creating file in path: %v", filePath)
		}
	}
}

// todo: create a proper error handler
func checkError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

// todo: maybe remove this and do a separate subcommand to write the result in a .csv file!
func writeHeader(pop []Elem, f *os.File) {
	for i := range pop {
		fmt.Fprintf(f, "pop[%d]\t", i)
	}
	fmt.Fprintf(f, "\n")
}

// todo: maybe remove this and do a separate subcommand to write the result in a .csv file!
func writeGeneration(pop Elements, f *os.File) {
	qtdObjs := len(pop[0].objs)
	for i := 0; i < qtdObjs; i++ {
		for _, p := range pop {
			fmt.Fprintf(f, "%10.3f\t", p.objs[i])
		}
		fmt.Fprintf(f, "\n")
	}
}
