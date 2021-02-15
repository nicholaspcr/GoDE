package main

import (
	"fmt"
	"log"
	"os"

	mo "gitlab.com/nicholaspcr/go-de/pkg/mode"
)

func main() {
	p := mo.Params{
		EXECS: 1,
		NP:    100,
		M:     3,
		DIM:   7,
		GEN:   100,
		FLOOR: 0,
		CEIL:  1,
		CR:    0.2,
		F:     0.2,
		P:     0,
	}

	// generate elements
	elems := mo.GeneratePopulation(p)

	// evaluate function
	evaluate := mo.GetProblemByName("dtlz1")

	f, err := os.Create("population.txt")
	if err != nil {
		log.Fatal("failed to create file")
	}

	// writing the population in data
	for i := range elems {
		evaluate.Fn(&elems[i], p.M)
		fmt.Fprintln(f, elems[i])
	}

	f, err = os.Create("filtered.txt")
	if err != nil {
		log.Fatal("failed to create file")
	}

	nd, _ := mo.FilterDominated(elems)
	// writing the population in data
	for i := range nd {
		fmt.Fprintln(f, elems[i])
	}
}
