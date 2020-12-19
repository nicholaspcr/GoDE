package mo

import (
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"time"
)

// DE -> runs a simple multiObjective DE in the ZDT1 case
func DE(
	NP, M, DIM, GEN int,
	LOWER, UPPER, CR, F float64,
	OUTDIR string,
	f *os.File,
) []Elem {
	if f == nil {
		var err error
		checkFilePath(OUTDIR)
		var path string = OUTDIR + "/paretoFront"
		checkFilePath(path)
		path += "/rand1.csv"
		f, err = os.Create(path)
		checkError(err)
	}
	defer f.Close()

	// Rand Seed
	rand.Seed(time.Now().UTC().UnixNano())

	// setting the test case, example: ZDT1
	evaluate := DTLZ3
	// generates random population
	population := generatePopulation(NP, DIM, LOWER, UPPER)
	for i := range population {
		err := evaluate(&population[i], M)
		checkError(err)
	}

	// fmt.Println(DIM)
	// fmt.Println(len(population[0].X))
	// fmt.Println(len(population[0].objs))

	writeHeader(population, f)
	writeGeneration(population, f)

	// Finish --- Writing on f
	for currGen := 0; currGen < GEN; currGen++ {
		// fmt.Println(len(population[0].X))
		// trial population vector
		trial := make([]Elem, NP)
		for i, p := range population {
			trial[i] = p.makeCpy()
		}
		for i, t := range trial {
			inds := make([]int, 3)
			err := generateIndices(0, NP, inds)
			checkError(err)
			a, b, c := population[inds[0]], population[inds[1]], population[inds[2]]

			// CROSS OVER
			currInd := rand.Int() % DIM
			for j := 0; j < DIM; j++ {
				changeProb := rand.Float64()
				if changeProb < CR || currInd == DIM {
					t.X[currInd] = a.X[currInd] + F*(b.X[currInd]-c.X[currInd])
				}

				if t.X[currInd] < LOWER {
					t.X[currInd] = LOWER
				}
				if t.X[currInd] > UPPER {
					t.X[currInd] = UPPER
				}
				currInd = (currInd + 1) % DIM
			}

			// for ZDT4
			// if t.X[0] > 1.0 {
			// 	t.X[0] = 1.0
			// } else if t.X[0] < 0 {
			// 	t.X[0] = 0.0
			// }

			evalErr := evaluate(&t, M)
			checkError(evalErr)

			if t.dominates(population[i]) {
				population[i] = t.makeCpy()
			} else if !population[i].dominates(t) {
				population = append(population, t.makeCpy())
			}
		}

		population = reduceByCrowdDistance(population, NP)

		writeGeneration(population, f)
	}
	return population
}

// MultiExecutions returns the pareto front of the total of 30 executions of the same problem
func MultiExecutions(
	EXECS, NP, M, DIM, GEN int,
	LOWER, UPPER, CR, F float64,
) {
	outDir := os.Getenv("HOME") + "./goDE/paretoFront"
	checkFilePath(outDir)

	// obtains the union of the points of all executions
	var arrElem []Elem
	for i := 0; i < EXECS; i++ {
		f, err := os.Create(outDir + "/exec-" + strconv.Itoa(i+1) + ".csv")
		checkError(err)
		currSlice := DE(NP, M, DIM, GEN, LOWER, UPPER, CR, F, outDir, f)
		for _, e := range currSlice {
			arrElem = append(arrElem, e)
		}
	}

	// filter those elements who are not dominated
	var result []Elem
	for i, first := range arrElem {
		flag := true
		for j, second := range arrElem {
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
	checkFilePath(outDir)
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
				inRank = append(inRank, pop[i].makeCpy())
			} else {
				notInRank = append(notInRank, pop[i].makeCpy())
			}
		}

		qtdElements += len(inRank)
		for i := 0; i < len(inRank); i++ {
			ranks[rankSize] = append(ranks[rankSize], inRank[i].makeCpy())
		}

		pop = make([]Elem, 0)
		for i := 0; i < len(notInRank); i++ {
			pop = append(pop, notInRank[i].makeCpy())
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
			if len(rank) == 0 {
				continue
			}
			//end of tail
			if i == 0 || i == len(rank)-1 {
				rank[i].crwdst = rank[i].objs[0] + rank[i].objs[1]
			} else {
				rank[i].crwdst = math.Abs(rank[i-1].objs[0]-rank[i+1].objs[0]) +
					math.Abs(rank[i-1].objs[1]-rank[i+1].objs[1])
			}
		}
		sort.Sort(byCrwdst(rank))

	}

	for _, rank := range ranks {
		for _, elem := range rank {
			pop = append(pop, elem.makeCpy())
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

// generates a group of points with random values
func generatePopulation(NP, DIM int, LOWER, UPPER float64) []Elem {
	ret := make([]Elem, NP)
	constant := UPPER - LOWER // range between floor and ceiling
	for i := 0; i < NP; i++ {
		ret[i].X = make([]float64, DIM)

		for j := 0; j < DIM; j++ {
			ret[i].X[j] = rand.Float64()*constant + LOWER // value varies within [lower,upper]
		}

		// for ZDT4
		// ret[i].X[0] = rand.Float64()
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
		os.Mkdir(filePath, os.ModePerm)
	}
}

func checkError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func writeHeader(pop []Elem, f *os.File) {
	for i := range pop {
		fmt.Fprintf(f, "pop[%d]\t", i)
	}
	fmt.Fprintf(f, "\n")
}

func writeGeneration(pop []Elem, f *os.File) {
	qtdObjs := len(pop[0].objs)
	for i := 0; i < qtdObjs; i++ {
		for _, p := range pop {
			fmt.Fprintf(f, "%10.3f\t", p.objs[i])
		}
		fmt.Fprintf(f, "\n")
	}
}
