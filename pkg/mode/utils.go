package mo

import (
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"sort"
	"strings"
)

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

// checks existance of filePath
func checkFilePath(basePath, filePath string) {
	folders := strings.Split(filePath, "/")
	for _, folder := range folders {
		basePath += "/" + folder
		if _, err := os.Stat(basePath); os.IsNotExist(err) {
			err = os.Mkdir(basePath, os.ModePerm)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

// todo: create a proper error handler
func checkError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

// returns NP elements filtered by rank and crwod distance
func reduceByCrowdDistance(elems Elements, pareto *Elements, NP int) Elements {
	ranks := fastNonDominatedRanking(elems)
	fmt.Println(len(elems))
	for _, r := range ranks {
		fmt.Print(fmt.Sprint(len(r)) + " ")
	}
	elems = make(Elements, 0)
	//sorting each rank by crowd distance
	for i := range ranks {
		calculateCrwdDist(ranks[i])
		sort.Sort(byCrwdst(ranks[i]))
	}

	// writes the pareto ranked into the pareto db
	*pareto = append(*pareto, ranks[0]...)

	for _, rank := range ranks {
		elems = append(elems, rank...)
		if len(elems) > NP {
			elems = elems[:50]
			break
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
		ranked, nonRanked := filterDominated(elems)
		ranks[currentRank] = append(ranks[currentRank], ranked...)
		currentRank++
		elems = make(Elements, 0)
		elems = append(elems, nonRanked...)
	}
	return ranks
}

// x is best 	-> -1
// y is best 	-> 	1
// else 			->	0
func dominanceTest(x, y *[]float64) int {
	result := 0
	for i := range *x {
		if (*x)[i] > (*y)[i] {
			if result == -1 {
				return 0
			}
			result = 1
		} else if (*y)[i] > (*x)[i] {
			if result == 1 {
				return 0
			}
			result = -1
		}
	}
	return result
}

func fastNonDominatedRanking(elems Elements) map[int]Elements {
	dominatingIth := make([]int, len(elems))
	ithDominated := make([][]Elem, len(elems))
	front := make([][]int, len(elems)+1)

	for p := 0; p < len(elems)-1; p++ {
		for q := p + 1; q < len(elems); q++ {
			dominanceTestResult := dominanceTest(&elems[p].objs, &elems[q].objs)
			if dominanceTestResult == -1 {
				ithDominated[p] = append(ithDominated[p], elems[q])
				dominatingIth[q]++
			} else if dominanceTestResult == 1 {
				ithDominated[q] = append(ithDominated[q], elems[p])
				dominatingIth[p]++
			}
		}
	}
	for i := 0; i < len(elems); i++ {
		if dominatingIth[i] == 0 {
			front[0] = append(front[0], i)
		}
	}
	i := 0
	for len(front[i]) != 0 {
		i++
		for p := range front[i-1] {
			if p <= len(ithDominated) {
				for q := range ithDominated[p] {
					dominatingIth[q]--
					if dominatingIth[q] == 0 {
						front[i] = append(front[i], q)
					}
				}
			}
		}
	}
	rankedSubList := make(map[int]Elements)
	for j := 0; j < i; j++ {
		for m := range front[j] {
			rankedSubList[j] = append(rankedSubList[j], elems[front[j][m]].Copy())
		}
	}
	return rankedSubList
}

// filterDominated -> returns elements that are not dominated in the set
func filterDominated(elems Elements) (nonDominated, dominated Elements) {
	sort.Sort(byFirstObj(elems))
	nonDom := make(Elements, 0)
	dom := make(Elements, 0)
	for i := len(elems) - 1; i >= 0; i-- {
		flag := true
		for j, second := range elems {
			if i == j {
				continue
			}
			if second.dominates(elems[i]) {
				flag = false
				break
			}
		}
		if flag {
			nonDom = append(nonDom, elems[i])
		} else {
			dom = append(dom, elems[i])
		}
	}
	return nonDom, dom
}

// assumes that the slice is composed of non dominated elements
func calculateCrwdDist(elems Elements) {
	if len(elems) <= 3 {
		return
	}
	for i := range elems {
		elems[i].crwdst = 0 // resets the crwdst
	}
	szObjs := len(elems[0].objs)
	for m := 0; m < szObjs; m++ {
		// sort by current objective
		sort.SliceStable(elems, func(i, j int) bool {
			return elems[i].objs[m] < elems[j].objs[m]
		})

		objMin := elems[0].objs[m]
		objMax := elems[len(elems)-1].objs[m]
		elems[0].crwdst = math.MaxFloat64
		elems[len(elems)-1].crwdst = math.MaxFloat64
		for i := 1; i < len(elems)-1; i++ {
			distance := elems[i+1].objs[m] - elems[i-1].objs[m]
			if math.Abs(objMax-objMin) != 0 {
				distance = distance / (objMax - objMin)
			}
			elems[i].crwdst += distance
		}
	}
}
