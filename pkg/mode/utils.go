package mo

import (
	"log"
	"math"
	"math/rand"
	"sort"

	"gitlab.com/nicholaspcr/go-de/pkg/problems/models"
)

// GeneratePopulation - creates a population without objs calculates
func GeneratePopulation(p models.Params) models.Elements {
	ret := make(models.Elements, p.NP)
	constant := p.CEIL - p.FLOOR // range between floor and ceiling
	for i := 0; i < p.NP; i++ {
		ret[i].X = make([]float64, p.DIM)

		for j := 0; j < p.DIM; j++ {
			ret[i].X[j] = rand.Float64()*constant + p.FLOOR // value varies within [ceil,upper]
		}
	}
	return ret
}

// todo: create a proper error handler
func checkError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

// ReduceByCrowdDistance - returns NP models.elements filtered by rank and crowd distance
func ReduceByCrowdDistance(elems models.Elements, NP int) (models.Elements, models.Elements) {
	ranks := FastNonDominatedRanking(elems)

	qtdElems := 0
	for _, r := range ranks {
		qtdElems += len(r)
	}
	if qtdElems < NP {
		log.Println("elems -> ", qtdElems)
		log.Fatal("less models.elements than NP")
	}

	elems = models.Elements{} // clears it

	for i := 0; i < len(ranks); i++ {
		CalculateCrwdDist(ranks[i])
		sort.SliceStable(ranks[i], func(l, r int) bool {
			return ranks[i][l].Crwdst > ranks[i][r].Crwdst
		})

		if len(elems)+len(ranks[i]) >= NP {
			counter := 0
			for len(elems) < NP {
				elems = append(elems, ranks[i][counter])
				counter++
			}
			break
		} else {
			elems = append(elems, ranks[i]...)
		}
	}

	newElems := make(models.Elements, len(elems))
	copy(newElems, elems)
	zero := make(models.Elements, len(ranks[0]))
	copy(zero, ranks[0])
	return newElems, zero
}

// FastNonDominatedRanking - ranks the models.elements and returns a map with models.elements per rank
func FastNonDominatedRanking(elems models.Elements) map[int]models.Elements {
	dominatingIth := make([]int, len(elems))
	ithDominated := make([][]int, len(elems))
	fronts := make([][]int, len(elems)+1)

	rand.Shuffle(len(elems), func(l, r int) {
		elems[l], elems[r] = elems[r], elems[l]
	})

	for i := range fronts {
		fronts[i] = make([]int, 0)
	}

	for p := 0; p < len(elems); p++ {
		ithDominated[p] = make([]int, 0) // S_p size 0
		dominatingIth[p] = 0             // N_p = 0

		for q := 0; q < len(elems); q++ {
			dominanceTestResult := DominanceTest(elems[p].Objs, elems[q].Objs)

			if dominanceTestResult == -1 { // p dominates q
				ithDominated[p] = append(ithDominated[p], q)
			} else if dominanceTestResult == 1 { // q dominates p
				dominatingIth[p]++
			}
		}
		if dominatingIth[p] == 0 {
			fronts[0] = append(fronts[0], p)
		}
	}

	for i := 1; i < len(fronts); i++ {
		for _, p := range fronts[i-1] { // for each p in F_i
			for _, q := range ithDominated[p] { // for each q in S_p
				dominatingIth[q]--
				if dominatingIth[q] == 0 {
					fronts[i] = append(fronts[i], q)
				}
			}
		}
	}

	// getting ranked models.elements from their index
	rankedSubList := make(map[int]models.Elements)
	for i := 0; i < len(fronts); i++ {
		for m := range fronts[i] {
			rankedSubList[i] = append(
				rankedSubList[i],
				elems[fronts[i][m]].Copy(),
			)
		}
	}

	return rankedSubList
}

// DominanceTest - results meanings:
//
//  - '-1': x is best
//  - '1': y is best
//  - '0': nobody dominates
func DominanceTest(x, y []float64) int {
	result := 0
	for i := range x {
		if (x)[i] > (y)[i] {
			if result == -1 {
				return 0
			}
			result = 1
		}
		if (y)[i] > (x)[i] {
			if result == 1 {
				return 0
			}
			result = -1
		}
	}
	return result
}

// FilterDominated -> returns models.elements that are not dominated in the set
func FilterDominated(elems models.Elements) (models.Elements, models.Elements) {
	nonDominated := make(models.Elements, 0)
	dominated := make(models.Elements, 0)

	for p := 0; p < len(elems); p++ {
		counter := 0
		for q := 0; q < len(elems); q++ {
			if p == q {
				continue
			}
			// q dominates the p element
			if DominanceTest(elems[p].Objs, elems[q].Objs) == 1 {
				counter++
			}
		}
		if counter == 0 {
			nonDominated = append(nonDominated, elems[p])
		} else {
			dominated = append(dominated, elems[p])
		}
	}

	nd := make(models.Elements, len(nonDominated))
	copy(nd, nonDominated)

	d := make(models.Elements, len(dominated))
	copy(d, dominated)

	return nd, d
}

// CalculateCrwdDist - assumes that the slice is composed of non dominated models.elements
func CalculateCrwdDist(elems models.Elements) {
	if len(elems) <= 2 {
		for i := range elems {
			elems[i].Crwdst = math.MaxFloat32
		}
		return
	}

	for i := range elems {
		elems[i].Crwdst = 0 // resets the crwdst
	}
	szObjs := len(elems[0].Objs)
	for m := 0; m < szObjs; m++ {
		// sort by current objective
		sort.SliceStable(elems, func(i, j int) bool {
			return elems[i].Objs[m] < elems[j].Objs[m]
		})

		objMin := elems[0].Objs[m]
		objMax := elems[len(elems)-1].Objs[m]
		elems[0].Crwdst = math.MaxFloat32
		elems[len(elems)-1].Crwdst = math.MaxFloat32
		for i := 1; i < len(elems)-1; i++ {
			distance := elems[i+1].Objs[m] - elems[i-1].Objs[m]

			op := objMax - objMin
			if op < 0 {
				op *= -1
			}
			if op > 0 {
				elems[i].Crwdst += distance / (objMax - objMin)
			}
		}
	}
}
