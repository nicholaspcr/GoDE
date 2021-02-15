package mo

import (
	"encoding/csv"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

// incode variable to decide to write or not to write files
var writeAllPoints = true

// MultiExecutions returns the pareto front of the total of 30 executions of the same problem
func MultiExecutions(params Params, prob ProblemFn, variant VariantFn, disablePlot bool) {
	homePath := os.Getenv("HOME")
	paretoPath := "/.go-de/mode/paretoFront/" + prob.Name + "/" + variant.Name
	if variant.Name == "pbest" {
		paretoPath += "/P-" + fmt.Sprint(params.P)
	}

	fmt.Println(prob.Name)

	checkFilePath(homePath, paretoPath)

	startTimer := time.Now()                 //	timer start
	rand.Seed(time.Now().UTC().UnixNano())   // Rand Seed
	population := GeneratePopulation(params) // random generated population

	lastGenChan := make(chan Elements, params.EXECS) // channel to get elems related to the last gen
	rankedChan := make(chan Elements, params.EXECS)  // channel to get elems related to rank[0] pareto

	// getting the maximum calculated value for each objective
	execsObjsValues := make([][]float64, params.EXECS)
	for i := range execsObjsValues {
		// todo: this only works with dtlz i think
		// using M value to set the amoung of objectives
		execsObjsValues[i] = make([]float64, params.M)
	}

	var wg sync.WaitGroup // number of working go routines

	// runs GDE3 for EXECS amount of times
	for i := 0; i < params.EXECS; i++ {
		filePath := homePath + paretoPath + "/exec-" + strconv.Itoa(i+1) + ".csv"
		var f *os.File
		var err error

		f, err = os.Create(filePath)
		checkError(err)

		wg.Add(1)

		// worker
		go func(i int) {
			defer wg.Done()
			GD3(
				lastGenChan,
				rankedChan,
				&execsObjsValues[i],
				params,
				prob.Fn,
				variant,
				population.Copy(),
				f,
			)
		}(i)
	}
	// closer
	fmt.Println("waiting for the executions to be done")
	wg.Wait()
	close(lastGenChan)
	close(rankedChan)

	if !disablePlot {
		// gets data from the pareto created in the last generation
		var lastGenPareto Elements
		for v := range lastGenChan {
			lastGenPareto = append(lastGenPareto, v...)
			lastGenPareto, _ = FilterDominated(lastGenPareto)
			rand.Shuffle(len(lastGenPareto), func(i, j int) {
				lastGenPareto[i], lastGenPareto[j] = lastGenPareto[j].Copy(), lastGenPareto[i].Copy()
			})
			// puts a cap on the solution's amount of points
			if len(lastGenPareto) > 500 {
				lastGenPareto = lastGenPareto[:500]
			}
		}

		counter := 0
		// gets data from the pareto created by rank[0] of each gen
		var rankedPareto Elements
		for v := range rankedChan {
			fmt.Printf("exec-%d\n", counter)
			counter++

			rankedPareto = append(rankedPareto, v...)
			rankedPareto, _ = FilterDominated(rankedPareto)
			rand.Shuffle(len(rankedPareto), func(i, j int) {
				rankedPareto[i], rankedPareto[j] = rankedPareto[j].Copy(), rankedPareto[i].Copy()
			})
			if len(rankedPareto) > 1000 {
				rankedPareto = rankedPareto[:1000]
			}
		}

		// checks path for the path used to store the details of each generation
		multiExecutionsPath := "/.go-de/mode/multiExecutions/" + prob.Name + "/" + variant.Name
		if variant.Name == "pbest" {
			multiExecutionsPath += "/P-" + fmt.Sprint(params.P)
		}
		checkFilePath(homePath, multiExecutionsPath)

		// results of the normal pareto
		writeResult(
			homePath+multiExecutionsPath+"/lastPareto.csv",
			lastGenPareto,
		)

		// result of the ranked pareto
		writeResult(
			homePath+multiExecutionsPath+"/rankedPareto.csv",
			rankedPareto,
		)

		fmt.Println(rankedPareto[0].X)
	}
	fmt.Println("Done writing file!")
	timeSpent := time.Since(startTimer)
	fmt.Println("time spend on executions: ", timeSpent)
	// getting biggest objs values
	maxObjs := make([]float64, len(execsObjsValues[0]))
	for i := range maxObjs {
		for j := range execsObjsValues {
			maxObjs[i] = math.Max(maxObjs[i], execsObjsValues[j][i])
		}
	}
	fmt.Println("maximum objective values found")
	fmt.Println(maxObjs)
}

// tokens is a counting semaphore use to
// enforce  a limit of 10 concurrent requests
var tokens = make(chan struct{}, 10)

// GD3 -> runs a simple multiObjective DE in the ZDT1 case
func GD3(
	normalCh chan<- Elements,
	rankedCh chan<- Elements,
	maximumObjs *[]float64,
	p Params,
	evaluate func(e *Elem, M int) error,
	variant VariantFn,
	population Elements,
	f *os.File,
) {
	// adding to  concurrent queue
	tokens <- struct{}{}
	defer f.Close()

	// gets new random token
	rand.Seed(time.Now().UnixNano())

	// var writer *csv.Writer
	writer := csv.NewWriter(f)
	writer.Comma = '\t'

	// calculates the objs of the inital population
	for i := range population {
		err := evaluate(&population[i], p.M)
		checkError(err)
	}

	writeHeader(population, writer)
	writeGeneration(population, writer)

	// stores the rank[0] of each generation
	bestElems := make(Elements, 0)

	var genRankZero Elements
	var bestInGen Elements

	for g := 0; g < p.GEN; g++ {
		genRankZero, _ = FilterDominated(population)

		for i := 0; i < len(population); i++ {
			vr, err := variant.fn(
				population,
				genRankZero,
				varParams{
					currPos: i,
					DIM:     p.DIM,
					F:       p.F,
					P:       p.P,
				})
			checkError(err)

			// trial element
			trial := population[i].Copy()

			// CROSS OVER
			currInd := rand.Int() % p.DIM
			randLucky := rand.Int() % p.DIM
			for j := 0; j < p.DIM; j++ {
				changeProb := rand.Float64()
				if changeProb < p.CR || currInd == randLucky {
					trial.X[currInd] = vr.X[currInd]
				}
				if trial.X[currInd] < p.FLOOR {
					trial.X[currInd] = p.FLOOR
				}
				if trial.X[currInd] > p.CEIL {
					trial.X[currInd] = p.CEIL
				}
				currInd = (currInd + 1) % p.DIM
			}

			evalErr := evaluate(&trial, p.M)
			checkError(evalErr)

			// SELECTION
			// if trial.dominates(population[i]) {
			// 	population[i] = trial.Copy()
			// } else if !population[i].dominates(trial) {
			// 	population = append(population, trial.Copy())
			// }

			comp := DominanceTest(&population[i].Objs, &trial.Objs)
			if comp == 1 {
				population[i] = trial.Copy()
			} else if comp == 0 {
				population = append(population, trial.Copy())
			}
		}

		population, bestInGen = ReduceByCrowdDistance(population, p.NP)
		bestElems = append(bestElems, bestInGen...)

		writeGeneration(population, writer)

		// checks for the biggest objective
		for _, p := range population {
			for i := range p.Objs {
				if p.Objs[i] > (*maximumObjs)[i] {
					(*maximumObjs)[i] = p.Objs[i]
				}
			}
		}
	}
	normalCh <- population
	rankedCh <- bestElems

	// clearing concurrent queue
	<-tokens
}
