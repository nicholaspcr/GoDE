package mo

import (
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
	outDir := os.Getenv("HOME") + "/.goDE/mode"
	checkFilePath(outDir)
	paretoPath := outDir + "/paretoFront"
	checkFilePath(paretoPath)

	population := generatePopulation(params) // random generated population
	var wg sync.WaitGroup                    // number of working go routines
	elemChan := make(chan Elements)
	for i := 0; i < params.EXECS; i++ {
		f, err := os.Create(paretoPath + "/exec-" + strconv.Itoa(i+1) + ".csv")
		checkError(err)
		wg.Add(1)
		// worker
		go func() {
			defer wg.Done()
			elemChan <- DE(params, prob, variant, population.Copy(), f)
		}()
	}
	// closer
	go func() {
		wg.Wait()
		close(elemChan)
	}()

	var pareto Elements // DE pareto front
	for i := 0; i < params.EXECS; i++ {
		v, ok := <-elemChan
		if !ok {
			fmt.Println("one of the goroutine workers didn't work")
		}
		pareto = append(pareto, v...)
	}
	result := filterDominated(pareto)             // non dominated set
	var path string = outDir + "/multiExecutions" // file path
	checkFilePath(path)
	path += "/rand1.csv"
	f, err := os.Create(path)
	checkError(err)
	defer f.Close()
	// writes in file
	for i := range result {
		fmt.Fprintf(f, "elem[%d]\t", i)
	}
	fmt.Fprintf(f, "\n")
	for i := 0; i < len(result[0].objs); i++ {
		for _, r := range result {
			fmt.Fprintf(f, "%10.3f\t", r.objs[i])
		}
		fmt.Fprintf(f, "\n")
	}
	fmt.Println("Done writing file!")
}

// DE -> runs a simple multiObjective DE in the ZDT1 case
func DE(
	p Params,
	evaluate ProblemFn,
	variant VariantFn,
	population Elements,
	f *os.File,
) Elements {
	defer f.Close()
	// Rand Seed
	rand.Seed(time.Now().UTC().UnixNano())
	for i := range population {
		err := evaluate(&population[i], p.M)
		checkError(err)
	}
	writeHeader(population, f)
	writeGeneration(population, f)

	for ; p.GEN > 0; p.GEN-- {
		trial := population.Copy() // trial population slice
		for i, t := range trial {
			v, err := variant(population, p)
			checkError(err)
			// CROSS OVER
			currInd := rand.Int() % p.DIM
			for j := 0; j < p.DIM; j++ {
				changeProb := rand.Float64()
				if changeProb < p.CR || currInd == p.DIM {
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
			if t.dominates(population[i]) {
				population[i] = t.Copy()
			} else if !population[i].dominates(t) {
				population = append(population, t.Copy())
			}
		}

		population = reduceByCrowdDistance(population, p.NP)
		writeGeneration(population, f)
	}
	return population
}

// returns NP elements filtered by rank and crwod distance
func reduceByCrowdDistance(elems Elements, NP int) Elements {
	ranks := rankElements(elems)
	elems = make(Elements, 0)
	//sorting each rank by crowd distance
	for i := range ranks {
		calcCrwdst(ranks[i])
		sort.Sort(byCrwdst(ranks[i]))
	}
	for _, rank := range ranks {
		for _, elem := range rank {
			elems = append(elems, elem.Copy())
			if len(elems) >= NP {
				return elems
			}
		}
	}
	return elems
}

// rankElements returna  map of dominating elements in ascending order
// destroys the slice provided
func rankElements(elems Elements) map[int]Elements {
	ranks := make(map[int]Elements)
	currentRank := 0
	for len(elems) > 0 {
		inRank := make(Elements, 0)
		notInRank := make(Elements, 0)
		for i := range elems {
			flag := true
			for j := range elems {
				if i == j {
					continue
				}
				if elems[j].dominates(elems[i]) {
					flag = false
					break
				}
			}
			if flag == true {
				inRank = append(inRank, elems[i].Copy())
			} else {
				notInRank = append(notInRank, elems[i].Copy())
			}
		}
		ranks[currentRank] = append(ranks[currentRank], inRank...)
		currentRank++
		elems = make(Elements, 0)
		elems = append(elems, notInRank...)
	}
	return ranks
}

// calcCrwdst -> calculates the crowd distance between elements
func calcCrwdst(elems Elements) {
	sort.Sort(byFirstObj(elems))
	for i := range elems {
		for j := range elems[i].objs {
			//end of tail
			if i == 0 {
				elems[i].crwdst += (elems[i+1].objs[j] - elems[i].objs[j]) * (elems[i+1].objs[j] - elems[i].objs[j])
			} else if i == len(elems)-1 {
				elems[i].crwdst += (elems[i-1].objs[j] - elems[i].objs[j]) * (elems[i-1].objs[j] - elems[i].objs[j])
			} else {
				elems[i].crwdst += (elems[i-1].objs[j] - elems[i+1].objs[j]) * (elems[i-1].objs[j] - elems[i+1].objs[j])
			}
			elems[i].crwdst = math.Sqrt(elems[i].crwdst)
		}
	}
}

// filterDominated -> returns elements that are not dominated in the set
func filterDominated(elems Elements) Elements {
	result := make(Elements, 0)
	for i, first := range elems {
		flag := true
		for j, second := range elems {
			if i == j {
				continue
			}
			if second.dominates(first) {
				flag = false
				break
			}
		}
		if flag {
			result = append(result, first)
		}
	}
	return result
}
