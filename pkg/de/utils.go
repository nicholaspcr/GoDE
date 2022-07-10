package de

import (
	"math/rand"

	"log"
	"math"
	"sort"

	"github.com/nicholaspcr/GoDE/pkg/models"
)

var (
	// INF is the maximum value used in the crowding distance
	INF = math.MaxFloat64 - 1e5
)

// GeneratePopulation - creates a population without objs calculates
func GeneratePopulation(p AlgorithmParams) models.Population {
	ret := make(models.Population, p.NP)
	for i := 0; i < p.NP; i++ {
		ret[i].X = make([]float64, p.DIM)

		for j := 0; j < p.DIM; j++ {
			// range between floor and ceiling
			constant := p.CEIL[j] - p.FLOOR[j]
			// value varies within [ceil,upper]
			ret[i].X[j] = rand.Float64()*constant + p.FLOOR[j]
		}
	}
	return ret
}

// ReduceByCrowdDistance - returns NP models.elements filtered by rank and
// crowd distance
func ReduceByCrowdDistance(
	elems models.Population,
	NP int,
) (models.Population, models.Population) {

	ranks := FastNonDominatedRanking(elems)
	elems = make(models.Population, 0)

	// TODO remove the qtdElems sections
	qtdElems := 0

	for _, r := range ranks {
		qtdElems += len(r)
	}

	if qtdElems < NP {
		log.Println("elems -> ", qtdElems)
		log.Fatal("less models.elements than NP")
	}

	for i := 0; i < len(ranks); i++ {
		CalculateCrwdDist(ranks[i])
		sort.SliceStable(ranks[i], func(l, r int) bool {
			return ranks[i][l].Crwdst > ranks[i][r].Crwdst
		})

		elems = append(elems, ranks[i]...)
		if len(elems) >= NP {
			elems = elems[:NP]
			break
		}
	}

	zero := make(models.Population, len(ranks[0]))
	copy(zero, ranks[0])
	return elems, zero
}

// FastNonDominatedRanking - ranks the models.elements and returns a map with
// models.elements per rank
func FastNonDominatedRanking(
	elems models.Population,
) map[int]models.Population {

	// this func is inspired by the DEB_NSGA-II paper
	// a fast and elitist multiobjective genetic algorithm

	dominatingIth := make([]int, len(elems))  // N_p equivalent
	ithDominated := make([][]int, len(elems)) // S_p equivalent
	fronts := make([][]int, 1)                // F equivalent
	fronts[0] = []int{}                       // initializes first front

	// TODO remove section
	//rand.Shuffle(len(elems), func(l, r int) {
	//	elems[l], elems[r] = elems[r], elems[l]
	//})

	for p := 0; p < len(elems); p++ {
		ithDominated[p] = make([]int, 0) // S_p size 0
		dominatingIth[p] = 0             // N_p = 0

		for q := 0; q < len(elems); q++ {
			dominanceTestResult := DominanceTest(elems[p].Objs, elems[q].Objs)

			if dominanceTestResult == -1 {
				// p dominates q
				// add q to the set of solutions dominated by p
				ithDominated[p] = append(ithDominated[p], q)

			} else if dominanceTestResult == 1 {
				// q dominates p
				// increment the domination counter of p
				dominatingIth[p]++
			}
		}
		if dominatingIth[p] == 0 {
			// adds p to the first front
			fronts[0] = append(fronts[0], p)
		}
	}

	// used to go through the existant fronts
	for i := 0; len(fronts[i]) > 0; i++ {
		// slice to be added to the next front
		nextFront := []int{}

		// for each p in F_i
		for _, p := range fronts[i] {

			// for each q in S_p
			for _, q := range ithDominated[p] {

				dominatingIth[q]--

				if dominatingIth[q] == 0 {
					nextFront = append(nextFront, q)
				}
			}
		}

		// adds the next front to the matrix
		fronts = append(fronts, nextFront)
	}

	// TODO remove section
	// previous method with matrix of fronts already instantiated
	//	for i := 1; i < len(fronts); i++ {
	//
	//		// for each p in F_i
	//		for _, p := range fronts[i-1] {
	//
	//			// for each q in S_p
	//			for _, q := range ithDominated[p] {
	//
	//				dominatingIth[q]--
	//
	//				if dominatingIth[q] == 0 {
	//					fronts[i] = append(fronts[i], q)
	//				}
	//			}
	//		}
	//	}

	// getting ranked models.elements from their index
	rankedSubList := make(map[int]models.Population)
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
		if x[i] > y[i] {
			if result == -1 {
				return 0
			}
			result = 1
		}
		if y[i] > x[i] {
			if result == 1 {
				return 0
			}
			result = -1
		}
	}
	return result
}

// FilterDominated -> returns models.elements that are not dominated in the set
func FilterDominated(
	elems models.Population,
) (models.Population, models.Population) {
	nonDominated := make(models.Population, 0)
	dominated := make(models.Population, 0)

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
			nonDominated = append(nonDominated, elems[p].Copy())
		} else {
			dominated = append(dominated, elems[p].Copy())
		}
	}

	return nonDominated, dominated
}

// CalculateCrwdDist - assumes that the slice is composed of non dominated
// models.elements
func CalculateCrwdDist(elems models.Population) {
	if len(elems) <= 2 {
		for i := range elems {
			elems[i].Crwdst = math.MaxFloat64
		}
		return
	}

	// resets the crwdst
	for i := range elems {
		elems[i].Crwdst = 0
	}

	szObjs := len(elems[0].Objs)

	for m := 0; m < szObjs; m++ {
		// sort by current objective
		sort.SliceStable(elems, func(i, j int) bool {
			return elems[i].Objs[m] < elems[j].Objs[m]
		})

		// obtain the extremes of the objective analysed
		objMin := elems[0].Objs[m]
		objMax := elems[len(elems)-1].Objs[m]

		// first and last receive max Crwdst value
		elems[0].Crwdst = INF
		elems[len(elems)-1].Crwdst = INF

		for i := 1; i < len(elems)-1; i++ {

			distance := elems[i+1].Objs[m] - elems[i-1].Objs[m]

			// if difference between extremes is less than 1e-8
			if objMax-objMin > 0 {
				distance = distance / (objMax - objMin)
			}

			// only adds to the crowdDistance if smalled than max value
			if elems[i].Crwdst+distance < INF {
				elems[i].Crwdst += distance
			}
		}
	}
}
