package mo

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"time"
)

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

	// Finish --- Writing on f
	for currGen := 0; currGen < p.GEN; currGen++ {
		// trial population vector
		trial := population.Copy()
		for i, t := range trial {
			e, err := variant(population, p)
			checkError(err)
			// CROSS OVER
			currInd := rand.Int() % p.DIM
			for j := 0; j < p.DIM; j++ {
				changeProb := rand.Float64()
				if changeProb < p.CR || currInd == p.DIM {
					t.X[currInd] = e.X[currInd]
				}

				if t.X[currInd] < p.FLOOR {
					t.X[currInd] = p.FLOOR
				}
				if t.X[currInd] > p.CEIL {
					t.X[currInd] = p.CEIL
				}
				currInd = (currInd + 1) % p.DIM
			}

			// for ZDT4
			// if t.X[0] > 1.0 {
			// 	t.X[0] = 1.0
			// } else if t.X[0] < 0 {
			// 	t.X[0] = 0.0
			// }

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

// MultiExecutions returns the pareto front of the total of 30 executions of the same problem
func MultiExecutions(params Params, prob ProblemFn, variant VariantFn) {
	outDir := os.Getenv("HOME") + "/.goDE/mode"
	checkFilePath(outDir)
	paretoPath := outDir + "/paretoFront"
	checkFilePath(paretoPath)

	// obtains the union of the points of all executions
	var pareto Elements
	// generates random population
	population := generatePopulation(params)
	for i := 0; i < params.EXECS; i++ {
		f, err := os.Create(paretoPath + "/exec-" + strconv.Itoa(i+1) + ".csv")
		checkError(err)
		currSlice := DE(params, prob, variant, population.Copy(), f)
		pareto = append(pareto, currSlice...)

	}

	// filter those elements who are not dominated
	var result []Elem
	for i, first := range pareto {
		flag := true
		for j, second := range pareto {
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

	// creates path and file
	var path string = outDir + "/multiExecutions"
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
	fmt.Println(result[len(result)-1].X)
	fmt.Println(result[len(result)-1].objs)
}

func reduceByCrowdDistance(pop []Elem, NP int) []Elem {
	rankSize := 0
	ranks := make(map[int]([]Elem))
	qtdElements := 0
	for qtdElements < NP {
		inRank := make([]Elem, 0)
		notInRank := make([]Elem, 0)
		for i := 0; i < len(pop); i++ {
			flag := true
			for j := 0; j < len(pop); j++ {
				if i == j {
					continue
				}
				if pop[j].dominates(pop[i]) {
					flag = false
					break
				}
			}
			if flag == true {
				inRank = append(inRank, pop[i].Copy())
			} else {
				notInRank = append(notInRank, pop[i].Copy())
			}
		}

		qtdElements += len(inRank)
		for i := 0; i < len(inRank); i++ {
			ranks[rankSize] = append(ranks[rankSize], inRank[i].Copy())
		}

		pop = make([]Elem, 0)
		for i := 0; i < len(notInRank); i++ {
			pop = append(pop, notInRank[i].Copy())
		}
		rankSize++

		if len(inRank) == 0 && len(notInRank) == 0 {
			fmt.Println("Ranking ERROR -> PopSz = ", len(pop))
			break
		}
	}

	pop = make([]Elem, 0)
	//sorting each rank by crowd distance
	for _, rank := range ranks {
		sort.Sort(byFirstObj(rank))
		// Calculates the distance for the points in the rank
		for i := 0; i < len(rank); i++ {
			for j := range rank[i].objs {
				//end of tail
				if i == 0 {
					rank[i].crwdst += (rank[i+1].objs[j] - rank[i].objs[j]) * (rank[i+1].objs[j] - rank[i].objs[j])
				} else if i == len(rank)-1 {
					rank[i].crwdst += (rank[i-1].objs[j] - rank[i].objs[j]) * (rank[i-1].objs[j] - rank[i].objs[j])
				} else {
					rank[i].crwdst += (rank[i-1].objs[j] - rank[i+1].objs[j]) * (rank[i-1].objs[j] - rank[i+1].objs[j])
				}
				rank[i].crwdst = math.Sqrt(rank[i].crwdst)
			}
		}
		sort.Sort(byCrwdst(rank))

	}

	for _, rank := range ranks {
		for _, elem := range rank {
			pop = append(pop, elem.Copy())
			if len(pop) >= NP {
				break
			}
		}
		if len(pop) >= NP {
			break
		}
	}
	return pop
}
