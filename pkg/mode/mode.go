package mo

import (
	"encoding/csv"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"
)

// MultiExecutions returns the pareto front of the total of 30 executions of the same problem
func MultiExecutions(params Params, prob ProblemFn, variant VariantFn) {
	homePath := os.Getenv("HOME")
	paretoPath := "/.go-de/mode/paretoFront/" + prob.Name + "/" + variant.Name
	checkFilePath(homePath, paretoPath)

	startTimer := time.Now()                 //	timer start
	rand.Seed(time.Now().UTC().UnixNano())   // Rand Seed
	population := generatePopulation(params) // random generated population

	var wg sync.WaitGroup                            // number of working go routines
	lastGenChan := make(chan Elements, params.EXECS) // channel to get elems related to the last gen
	rankedChan := make(chan Elements, params.EXECS)  // channel to get elems related to rank[0] pareto

	// getting the maximum calculated value for each objective
	execsObjsValues := make([][]float64, params.EXECS)
	for i := range execsObjsValues {
		// todo: this only works with dtlz i think
		// using M value to set the amoung of objectives
		execsObjsValues[i] = make([]float64, params.M)
	}

	// runs GDE3 for EXECS amount of times
	for i := 0; i < params.EXECS; i++ {
		f, err := os.Create(
			homePath +
				paretoPath +
				"/exec-" +
				strconv.Itoa(i+1) +
				".csv")
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
				prob.fn,
				variant,
				population.Copy(),
				f,
			)
		}(i)
	}
	// closer
	go func() {
		wg.Wait()
		close(lastGenChan)
		close(rankedChan)
	}()

	// gets data from the pareto created in the last generation
	var lastGenPareto Elements
	for v := range lastGenChan {
		lastGenPareto = append(lastGenPareto, v...)
		lastGenPareto, _ = filterDominated(lastGenPareto)

		calculateCrwdDist(lastGenPareto)

		// leaves first the non border points in the rank zero
		sort.SliceStable(lastGenPareto, func(i, j int) bool {
			// smaller than max float32
			if lastGenPareto[i].crwdst >= math.MaxFloat32-10.0 {
				return false
			}
			return lastGenPareto[i].crwdst > lastGenPareto[j].crwdst
		})

		// puts a cap on the solution's amount of points
		if len(lastGenPareto) > 2*params.NP {
			lastGenPareto = lastGenPareto[:2*params.NP]
		}
	}

	// gets data from the pareto created by rank[0] of each gen
	var rankedPareto Elements
	for v := range rankedChan {
		rankedPareto = append(rankedPareto, v...)
		rankedPareto, _ = filterDominated(rankedPareto)

		calculateCrwdDist(rankedPareto)

		// leaves first the non border points in the rank zero
		sort.SliceStable(rankedPareto, func(i, j int) bool {
			// smaller than max float32
			if rankedPareto[i].crwdst >= math.MaxFloat32-10.0 {
				return false
			}
			return rankedPareto[i].crwdst > rankedPareto[j].crwdst
		})

		// puts a cap on the solution's amount of points
		if len(rankedPareto) > 2*params.NP {
			rankedPareto = rankedPareto[:2*params.NP]
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
// enforce  a limit of 5 concurrent requests
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

	writer := csv.NewWriter(f)
	writer.Comma = '\t'
	defer f.Close()

	// calculates the objs of the inital population
	for i := range population {
		err := evaluate(&population[i], p.M)
		checkError(err)
	}
	writeHeader(population, writer)
	writeGeneration(population, writer)

	// stores the rank[0] of each generation
	bestElems := make(Elements, 0)
	// genRankZero -> stores the previous generation rank zero
	// it is used in the variants best1, best2 and currToBest1
	_, genRankZero := filterDominated(population)

	for ; p.GEN > 0; p.GEN-- {
		trial := population.Copy() // trial population slice

		// leaves first the non border points in the rank zero
		sort.SliceStable(genRankZero, func(i, j int) bool {
			// smaller than max float32
			if genRankZero[i].crwdst >= math.MaxFloat32-10.0 {
				return false
			}
			return genRankZero[i].crwdst > genRankZero[j].crwdst
		})
		for i, t := range trial {
			v, err := variant.fn(
				population,
				genRankZero,
				varParams{
					currPos: i,
					DIM:     p.DIM,
					F:       p.F,
				})
			checkError(err)

			// CROSS OVER
			currInd := rand.Int() % p.DIM
			for j := 0; j < p.DIM; j++ {
				changeProb := rand.Float64()
				if changeProb < p.CR || currInd == p.DIM-1 {
					t.X[currInd] = v.X[currInd]
				}
				if t.X[currInd] < p.FLOOR {
					t.X[currInd] = p.FLOOR
				}
				if t.X[currInd] > p.CEIL {
					t.X[currInd] = p.CEIL
				}
				currInd = (currInd + 1) % p.DIM
			}

			evalErr := evaluate(&t, p.M)
			checkError(evalErr)

			// SELECTION
			if t.dominates(population[i]) {
				population[i] = t.Copy()
			} else if !population[i].dominates(t) {
				population = append(population, t.Copy())
			}
		}

		population, genRankZero = reduceByCrowdDistance(population, p.NP)
		bestElems = append(bestElems, genRankZero...)

		// testing the filter of the bestElemetns here
		calculateCrwdDist(bestElems)

		// leaves first the non border points in the rank zero
		sort.SliceStable(bestElems, func(i, j int) bool {
			// smaller than max float32
			if bestElems[i].crwdst >= math.MaxFloat32-10.0 {
				return false
			}
			return bestElems[i].crwdst > bestElems[j].crwdst
		})

		// puts a cap on the solution's amount of points
		if len(bestElems) > p.NP {
			bestElems = bestElems[:p.NP]
		}

		writeGeneration(population, writer)

		// checks for the biggest objective
		for _, p := range population {
			for i := range p.objs {
				if p.objs[i] > (*maximumObjs)[i] {
					(*maximumObjs)[i] = p.objs[i]
				}
			}
		}
	}
	normalCh <- population
	rankedCh <- bestElems

	// clearing concurrent queue
	<-tokens
}
