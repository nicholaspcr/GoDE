package de

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/nicholaspcr/GoDE/pkg/algorithms"
	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/problems"
	"github.com/nicholaspcr/GoDE/pkg/writer"
)

// tokens is a counting semaphore use to
// enforce  a limit of 10 concurrent requests
var tokens = make(chan struct{}, 15)

// MultiExecutions returns the pareto front of the total of 30 executions of the
// same problem
func MultiExecutions(
	params models.AlgorithmParams,
	problem problems.Problem,
	variant models.Variant,
	initialPopulation models.Population,
) {
	homePath := os.Getenv("HOME")
	paretoPath := fmt.Sprintf(
		"/.gode/de/paretoFront/%s/%s",
		problem.Name(),
		variant.Name())

	if variant.Name() == "pbest" {
		paretoPath += "/P-" + fmt.Sprint(
			params.P,
		)
	}

	writer.CheckFilePath(
		homePath,
		paretoPath,
	)

	// channel to get elems related to rank[0] pareto
	rankedChan := make(
		chan models.Population,
		params.EXECS,
	)

	// getting the maximum calculated value for each objective
	maximumObjs := make(
		chan []float64,
		params.EXECS,
	)

	wg := &sync.WaitGroup{}

	// runs GDE3 for EXECS amount of times
	for i := 0; i < params.EXECS; i++ {
		filePath := fmt.Sprintf(
			"%s%s/exec-%d.csv",
			homePath,
			paretoPath,
			i+1,
		)

		f, err := os.Create(filePath)
		if err != nil {
			panic(err)
		}

		cpyPopulation := make(
			models.Population,
			len(initialPopulation),
		)
		copy(
			cpyPopulation,
			initialPopulation,
		)

		wg.Add(1)
		// worker
		go func() {
			// adding to concurrent queue
			tokens <- struct{}{}
			// cleaning concurrent queue
			defer func() { <-tokens }()
			// finishing worker
			defer func() { wg.Done() }()
			// running one execution of the GDE3
			algorithms.GDE3().Execute(
				rankedChan,
				maximumObjs,
				params,
				problem,
				variant,
				cpyPopulation,
				f,
			)
		}()
	}
	// closer
	fmt.Println(
		"waiting for the executions to be done",
	)

	go func() {
		wg.Wait()
		close(rankedChan)
		close(maximumObjs)
	}()

	fmt.Printf("execs: ")
	counter := 0
	// gets data from the pareto created by rank[0] of each gen
	var rankedPareto models.Population
	for v := range rankedChan {
		counter++
		fmt.Printf("%d, ", counter)

		rankedPareto = append(
			rankedPareto,
			v...)

		// gets non dominated and filters by crowdingDistance
		_, rankedPareto = algorithms.ReduceByCrowdDistance(
			rankedPareto,
			len(rankedPareto),
		)

		// limits the amounts of dots to 1k
		if len(rankedPareto) > 1000 {
			rankedPareto = rankedPareto[:1000]
		}
	}

	// checks path for the path used to store the details of each generation
	multiExecutionsPath := fmt.Sprintf(
		"/.gode/de/multiExecutions/%s/%s",
		problem.Name(),
		variant.Name(),
	)

	if variant.Name() == "pbest" {
		multiExecutionsPath += "/P-" + fmt.Sprint(
			params.P,
		)
	}
	writer.CheckFilePath(
		homePath,
		multiExecutionsPath,
	)

	// result of the ranked pareto
	f, err := os.Create(
		homePath + multiExecutionsPath + "/rankedPareto.csv",
	)
	if err != nil {
		log.Fatalln(
			"Failed to create the ranked pareto file",
		)
	}

	// creates writer and writes the elements objs
	w := writer.NewWriter(f)
	w.Comma = ';'
	w.WriteHeader(params.M)
	w.ElementsObjs(rankedPareto)

	// getting biggest objs values
	maxObjs := make([]float64, params.M)
	for arr := range maximumObjs {
		for i, obj := range arr {
			if obj > maxObjs[i] {
				maxObjs[i] = obj
			}
		}
	}
	fmt.Println(
		"maximum objective values found",
	)
	fmt.Println(maxObjs)
}
