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

	// creates file
	dataF, err := os.Create("population.txt")
	if err != nil {
		log.Fatal("error")
	}
	defer dataF.Close()
	objsF, err := os.Create("objectives.txt")
	if err != nil {
		log.Fatal("error")
	}
	defer objsF.Close()

	// evaluate function
	evaluate := mo.GetProblemByName("dtlz1")

	// writing the population in data
	for i := range elems {
		// variables
		fmt.Fprintln(dataF, elems[i].X)
	}

	for i := range elems {
		// calculating
		evaluate.Fn(&elems[i], p.M)
		fmt.Fprintln(objsF, elems[i].Objs)
	}
}
