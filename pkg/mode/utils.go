package mo

import (
	"errors"
	"log"
	"math"
	"math/rand"
	"sort"
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

// todo: create a proper error handler
func checkError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

// returns NP elements filtered by rank and crwod distance
func reduceByCrowdDistance(elems Elements, NP int) (reduceElements, rankZero Elements) {
	ranks := fastNonDominatedRanking(elems)
	totalRanksElements := 0
	for _, r := range ranks {
		totalRanksElements += len(r)
	}
	elems = make(Elements, 0)
	//sorting each rank by crowd distance
	for i := range ranks {
		calculateCrwdDist(ranks[i])
		sort.Sort(byCrwdst(ranks[i]))
	}

	for _, rank := range ranks {
		elems = append(elems, rank...)
		if len(elems) > NP {
			elems = elems[:50]
			break
		}
	}

	// todo: REVIEW quick fix for the ranking generating less than np elements
	for len(elems) < NP {
		randIndex := rand.Int() % len(elems)
		elems = append(elems, elems[randIndex])
	}
	return elems, ranks[0]
}

func fastNonDominatedRanking(elems Elements) map[int]Elements {
	dominatingIth := make([]int, len(elems))
	ithDominated := make([][]Elem, len(elems))
	fronts := make([][]int, len(elems)+1)

	for p := 0; p < len(elems)-1; p++ {
		for q := p + 1; q < len(elems); q++ {
			// dominanceTestResult := dominanceTest(&elems[p].objs, &elems[q].objs)
			// if dominanceTestResult == -1 {
			// 	ithDominated[p] = append(ithDominated[p], elems[q])
			// } else if dominanceTestResult == 1 {
			// 	dominatingIth[p]++
			// }
			if elems[p].dominates(elems[q]) {
				ithDominated[p] = append(ithDominated[p], elems[q])
			} else if elems[q].dominates(elems[p]) {
				dominatingIth[p]++
			}
		}
		if dominatingIth[p] == 0 {
			fronts[0] = append(fronts[0], p)
		}
	}

	for i := 1; i < len(fronts); i++ {
		for p := range fronts[i-1] {
			for q := range ithDominated[p] {
				dominatingIth[q]--
				if dominatingIth[q] == 0 {
					fronts[i] = append(fronts[i], q)
				}
			}
		}
	}
	rankedSubList := make(map[int]Elements)
	for i := range fronts {
		for m := range fronts[i] {
			rankedSubList[i] = append(rankedSubList[i], elems[fronts[i][m]].Copy())
		}
	}

	return rankedSubList
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
