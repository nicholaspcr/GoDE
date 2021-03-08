package mo

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"strconv"
	"sync"
	"time"
)

// MultiExecutions returns the pareto front of the total of 30 executions of the same problem
func MultiExecutions(
	params Params,
	prob ProblemFn,
	variant VariantFn,
	disablePlot bool,
) {
	if params.CPUProf != "" {
		f, err := os.Create(params.CPUProf)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}
	if params.MemProf != "" {
		f, err := os.Create(params.MemProf)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		runtime.GC()    // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}

	homePath := os.Getenv("HOME")
	paretoPath := "/.go-de/mode/paretoFront/" + prob.Name + "/" + variant.Name

	if variant.Name == "pbest" {
		paretoPath += "/P-" + fmt.Sprint(params.P)
	}

	checkFilePath(homePath, paretoPath)

	startTimer := time.Now()                 //	timer start
	rand.Seed(time.Now().UnixNano())         // Rand Seed
	population := GeneratePopulation(params) // random generated population

	rankedChan := make(chan Elements, params.EXECS) // channel to get elems related to rank[0] pareto

	// getting the maximum calculated value for each objective
	maximumObjs := make(chan []float64, params.EXECS)

	wg := &sync.WaitGroup{}

	// runs GDE3 for EXECS amount of times
	for i := 0; i < params.EXECS; i++ {
		filePath := homePath + paretoPath + "/exec-" + strconv.Itoa(i+1) + ".csv"

		f, err := os.Create(filePath)
		checkError(err)

		wg.Add(1)
		// worker
		go GD3(
			wg,
			rankedChan,
			maximumObjs,
			params,
			prob.Fn,
			variant,
			population.Copy(),
			f,
		)
	}
	// closer
	fmt.Println("waiting for the executions to be done")

	go func() {
		wg.Wait()
		close(rankedChan)
		close(maximumObjs)
	}()

	fmt.Printf("execs: ")
	counter := 0
	// gets data from the pareto created by rank[0] of each gen
	var rankedPareto Elements
	for v := range rankedChan {
		counter++
		fmt.Printf("%d, ", counter)

		rankedPareto = append(rankedPareto, v...)
		rankedPareto, _ = FilterDominated(rankedPareto)
		if len(rankedPareto) > 1000 {
			rankedPareto = rankedPareto[:1000]
		}
	}
	fmt.Printf("\n")

	// checks path for the path used to store the details of each generation
	multiExecutionsPath := "/.go-de/mode/multiExecutions/" + prob.Name + "/" + variant.Name
	if variant.Name == "pbest" {
		multiExecutionsPath += "/P-" + fmt.Sprint(params.P)
	}
	checkFilePath(homePath, multiExecutionsPath)

	// result of the ranked pareto
	writeResult(
		homePath+multiExecutionsPath+"/rankedPareto.csv",
		rankedPareto,
	)

	fmt.Println("Done writing file!")
	timeSpent := time.Since(startTimer)
	fmt.Println("time spend on executions: ", timeSpent)

	// getting biggest objs values
	maxObjs := make([]float64, params.M)
	for arr := range maximumObjs {
		for i, obj := range arr {
			maxObjs[i] = math.Max(maxObjs[i], obj)
		}
	}
	fmt.Println("maximum objective values found")
	fmt.Println(maxObjs)
}
