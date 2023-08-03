package de

import (
	"context"
	"math"
	"sort"

	"github.com/nicholaspcr/GoDE/pkg/models"
)

// INF is the maximum value used in the crowding distance
var INF = math.MaxFloat64 - 1e5

// ReduceByCrowdDistance - returns NP api.elements filtered by rank and crowd
// distance.
func ReduceByCrowdDistance(
	ctx context.Context, elems []models.Vector, NP int,
) ([]models.Vector, []models.Vector) {
	ranks := FastNonDominatedRanking(ctx, elems)
	elems = make([]models.Vector, 0)

	// TODO remove the qtdElems sections
	qtdElems := 0

	for _, r := range ranks {
		qtdElems += len(r)
	}

	for i := 0; i < len(ranks); i++ {
		CalculateCrwdDist(ranks[i])
		sort.SliceStable(ranks[i], func(l, r int) bool {
			return ranks[i][l].CrowdingDistance > ranks[i][r].CrowdingDistance
		})

		elems = append(elems, ranks[i]...)
		if len(elems) >= NP {
			elems = elems[:NP]
			break
		}
	}

	zero := make([]models.Vector, len(ranks[0]))

	// TODO NICK: is this the best method for copying the vectors?
	for idx, v := range ranks[0] {
		zero[idx] = v.Copy()
	}
	//copy(zero, ranks[0])

	return elems, zero
}

// FastNonDominatedRanking - ranks the API.elements and returns a map with
// api.elements per rank
func FastNonDominatedRanking(
	ctx context.Context, elems []models.Vector,
) map[int][]models.Vector {
	// This function is inspired by the DEB_NSGA-II paper a fast and elitist
	// multi-objective genetic algorithm

	dominatingIth := make([]int, len(elems))  // N_p equivalent
	ithDominated := make([][]int, len(elems)) // S_p equivalent
	fronts := make([][]int, 1)                // F equivalent
	fronts[0] = []int{}                       // initializes first front

	for p := 0; p < len(elems); p++ {
		ithDominated[p] = make([]int, 0) // S_p size 0
		dominatingIth[p] = 0             // N_p = 0

		for q := 0; q < len(elems); q++ {
			dominanceTestResult := DominanceTest(
				elems[p].Objectives, elems[q].Objectives,
			)

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

	// used to go through the existent fronts
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

	// getting ranked api.elements from their index
	rankedSubList := make(map[int][]models.Vector)
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
	elems []models.Vector,
) ([]models.Vector, []models.Vector) {
	nonDominated := make([]models.Vector, 0)
	dominated := make([]models.Vector, 0)

	for p := 0; p < len(elems); p++ {
		counter := 0
		for q := 0; q < len(elems); q++ {
			if p == q {
				continue
			}
			// q dominates the p element
			if DominanceTest(elems[p].Objectives, elems[q].Objectives) == 1 {
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
func CalculateCrwdDist(elems []models.Vector) {
	if len(elems) <= 2 {
		for i := range elems {
			elems[i].CrowdingDistance = math.MaxFloat64
		}
		return
	}

	// resets the crwdst
	for i := range elems {
		elems[i].CrowdingDistance = 0
	}

	szObjectives := len(elems[0].Objectives)

	for m := 0; m < szObjectives; m++ {
		// sort by current objective
		sort.SliceStable(elems, func(i, j int) bool {
			return elems[i].Objectives[m] < elems[j].Objectives[m]
		})

		// obtain the extremes of the objective analysed
		objMin := elems[0].Objectives[m]
		objMax := elems[len(elems)-1].Objectives[m]

		// first and last receive max CrowdingDistance value
		elems[0].CrowdingDistance = INF
		elems[len(elems)-1].CrowdingDistance = INF

		for i := 1; i < len(elems)-1; i++ {

			distance := elems[i+1].Objectives[m] - elems[i-1].Objectives[m]

			// if difference between extremes is less than 1e-8
			if objMax-objMin > 0 {
				distance = distance / (objMax - objMin)
			}

			// only adds to the crowdDistance if its smaller than max value
			if elems[i].CrowdingDistance+distance < INF {
				elems[i].CrowdingDistance += distance
			}
		}
	}
}
