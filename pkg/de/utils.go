package de

import (
	"context"
	"math"
	"sort"

	"github.com/nicholaspcr/GoDE/pkg/models"
)

// INF is the maximum value used in the crowding distance
var INF = math.MaxFloat64 - 1e5

// ReduceByCrowdDistance - returns np api.elements filtered by rank and crowd
// distance.
func ReduceByCrowdDistance(
	ctx context.Context, elems []models.Vector, np int,
) ([]models.Vector, []models.Vector) {
	ranks := FastNonDominatedRanking(ctx, elems)
	elems = make([]models.Vector, 0, np)

	for i := 0; i < len(ranks); i++ {
		// Check for cancellation and calculate crowding distance
		if err := CalculateCrwdDist(ctx, ranks[i]); err != nil {
			// Return what we have so far on cancellation
			return elems, nil
		}
		sort.SliceStable(ranks[i], func(l, r int) bool {
			return ranks[i][l].CrowdingDistance > ranks[i][r].CrowdingDistance
		})

		elems = append(elems, ranks[i]...)
		if len(elems) >= np {
			elems = elems[:np]
			break
		}
	}

	// Deep copy rank 0 vectors (cannot use builtin copy due to slice fields)
	// Pre-allocate with exact capacity to avoid re-allocations
	zero := make([]models.Vector, len(ranks[0]))
	for idx, v := range ranks[0] {
		zero[idx] = v.Copy()
	}

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

	for p := range len(elems) {
		ithDominated[p] = make([]int, 0) // S_p size 0
		dominatingIth[p] = 0             // N_p = 0

		for q := range len(elems) {
			dominanceTestResult := DominanceTest(
				elems[p].Objectives, elems[q].Objectives,
			)

			switch dominanceTestResult {
			case -1:
				// p dominates q
				// add q to the set of solutions dominated by p
				ithDominated[p] = append(ithDominated[p], q)
			case 1:
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

	for p := range len(elems) {
		counter := 0
		for q := range len(elems) {
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
// api.elements. Supports context cancellation for long-running calculations.
func CalculateCrwdDist(ctx context.Context, elems []models.Vector) error {
	if len(elems) <= 2 {
		for i := range elems {
			elems[i].CrowdingDistance = math.MaxFloat64
		}
		return nil
	}

	// Check for cancellation before starting
	if err := ctx.Err(); err != nil {
		return err
	}

	// resets the crwdst
	for i := range elems {
		elems[i].CrowdingDistance = 0
	}

	szObjectives := len(elems[0].Objectives)

	for m := range szObjectives {
		// Check for cancellation between objectives
		if err := ctx.Err(); err != nil {
			return err
		}

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
				distance /= (objMax - objMin)
			}

			// only adds to the crowdDistance if its smaller than max value
			if elems[i].CrowdingDistance+distance < INF {
				elems[i].CrowdingDistance += distance
			}
		}
	}
	return nil
}

// IncrementalParetoUpdate efficiently integrates new vectors into existing Pareto front.
// More efficient than re-ranking entire set from scratch.
func IncrementalParetoUpdate(
	ctx context.Context,
	currentPareto []models.Vector,
	newVectors []models.Vector,
	maxSize int,
) ([]models.Vector, []models.Vector) {
	// Merge current + new
	merged := make([]models.Vector, 0, len(currentPareto)+len(newVectors))
	merged = append(merged, currentPareto...)
	merged = append(merged, newVectors...)

	// Quick filter for obviously dominated vectors
	filtered := quickDominanceFilter(merged)

	// Full ranking only if needed
	if len(filtered) > maxSize*2 {
		return ReduceByCrowdDistance(ctx, filtered, maxSize)
	}

	// For smaller sets, optimize crowding distance calculation
	ranks := FastNonDominatedRanking(ctx, filtered)
	if err := CalculateCrwdDist(ctx, ranks[0]); err != nil {
		// Return current pareto on cancellation
		return currentPareto, nil
	}
	sort.SliceStable(ranks[0], func(i, j int) bool {
		return ranks[0][i].CrowdingDistance > ranks[0][j].CrowdingDistance
	})

	result := ranks[0]
	if len(result) > maxSize {
		result = result[:maxSize]
	}

	rankZero := make([]models.Vector, len(ranks[0]))
	for idx, v := range ranks[0] {
		rankZero[idx] = v.Copy()
	}

	return result, rankZero
}

// quickDominanceFilter removes obviously dominated vectors in O(n²).
// More efficient than full ranking when we just need basic filtering.
func quickDominanceFilter(vectors []models.Vector) []models.Vector {
	result := make([]models.Vector, 0, len(vectors))

	for i := range vectors {
		dominated := false
		for j := range vectors {
			if i == j {
				continue
			}
			if DominanceTest(vectors[i].Objectives, vectors[j].Objectives) == 1 {
				dominated = true
				break
			}
		}
		if !dominated {
			result = append(result, vectors[i])
		}
	}

	return result
}
