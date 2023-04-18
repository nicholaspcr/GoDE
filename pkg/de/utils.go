package de

import (
	"math/rand"

	"log"
	"math"
	"sort"

	"github.com/nicholaspcr/GoDE/pkg/api"
)

var (
	// INF is the maximum value used in the crowding distance
	INF = math.MaxFloat64 - 1e5
)

// GeneratePopulation fills the vectors of a given population, does not
// generate the values for its objective functions.
func GeneratePopulation(p *api.Population, params api.PopulationParameters) {
	for i := 0; i < len(p.Vectors); i++ {
		p.Vectors[i].Elements = make([]float64, params.DimensionsSize)

		for j := 0; j < p.DimSize(); j++ {
			// range between floor and ceiling
			constant := p.Ceils()[j] - p.Floors()[j]
			// value varies within [ceil,upper]
			p.Vectors[i].X[j] = rand.Float64()*constant + p.Floors()[j]
		}
	}
}

// ReduceByCrowdDistance - returns NP api.elements filtered by rank and
// crowd distance.
func ReduceByCrowdDistance(
	elems []api.Vector,
	NP int,
) ([]api.Vector, []api.Vector) {

	ranks := FastNonDominatedRanking(elems)
	elems = make([]api.Vector, 0)

	// TODO remove the qtdElems sections
	qtdElems := 0

	for _, r := range ranks {
		qtdElems += len(r)
	}

	if qtdElems < NP {
		log.Println("elems -> ", qtdElems)
		log.Fatal("less api.elements than NP")
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

	zero := make([]api.Vector, len(ranks[0]))
	copy(zero, ranks[0])
	return elems, zero
}

// FastNonDominatedRanking - ranks the api.elements and returns a map with
// api.elements per rank
func FastNonDominatedRanking(
	elems []api.Vector,
) map[int][]api.Vector {

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

	// getting ranked api.elements from their index
	rankedSubList := make(map[int][]api.Vector)
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
//   - '-1': x is best
//   - '1': y is best
//   - '0': nobody dominates
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

// FilterDominated -> returns api.elements that are not dominated in the set
func FilterDominated(
	elems []api.Vector,
) ([]api.Vector, []api.Vector) {
	nonDominated := make([]api.Vector, 0)
	dominated := make([]api.Vector, 0)

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
// api.elements
func CalculateCrwdDist(elems []api.Vector) {
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
