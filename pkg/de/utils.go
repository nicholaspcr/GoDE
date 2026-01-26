package de

import (
	"context"
	"math"
	"sort"

	"github.com/nicholaspcr/GoDE/pkg/models"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// INF is the maximum value used in the crowding distance
var INF = math.MaxFloat64 - 1e5

// ReduceByCrowdDistance - returns np api.elements filtered by rank and crowd
// distance.
func ReduceByCrowdDistance(
	ctx context.Context, elems []models.Vector, np int,
) ([]models.Vector, []models.Vector) {
	tracer := otel.Tracer("de")
	ctx, span := tracer.Start(ctx, "de.ReduceByCrowdDistance",
		trace.WithAttributes(
			attribute.Int("input_size", len(elems)),
			attribute.Int("target_size", np),
		),
	)
	defer span.End()

	ranks := FastNonDominatedRanking(ctx, elems)

	// Deep copy rank 0 vectors early to avoid reprocessing
	zero := make([]models.Vector, len(ranks[0]))
	for idx, v := range ranks[0] {
		zero[idx] = v.Copy()
	}

	// Optimization: Copy directly to result to avoid intermediate collection
	result := make([]models.Vector, 0, np)
	for i := 0; i < len(ranks); i++ {
		// Check for cancellation and calculate crowding distance
		if err := CalculateCrwdDist(ctx, ranks[i]); err != nil {
			// Return what we have so far on cancellation
			return result, nil
		}
		sort.SliceStable(ranks[i], func(l, r int) bool {
			return ranks[i][l].CrowdingDistance > ranks[i][r].CrowdingDistance
		})

		// Copy vectors directly to result
		for _, v := range ranks[i] {
			if len(result) >= np {
				span.SetAttributes(
					attribute.Int("result_size", len(result)),
					attribute.Int("rank_zero_size", len(zero)),
					attribute.Int("num_ranks", len(ranks)),
				)
				return result, zero
			}
			result = append(result, v.Copy())
		}
	}

	span.SetAttributes(
		attribute.Int("result_size", len(result)),
		attribute.Int("rank_zero_size", len(zero)),
		attribute.Int("num_ranks", len(ranks)),
	)
	return result, zero
}

// FastNonDominatedRanking ranks solutions into Pareto fronts using the NSGA-II algorithm.
//
// This implements the fast non-dominated sorting procedure from:
// "A Fast and Elitist Multiobjective Genetic Algorithm: NSGA-II" by Deb et al. (2002)
//
// The algorithm assigns each solution to a front (rank):
//   - Front 0 (Rank 0): Non-dominated solutions (Pareto optimal in the current population)
//   - Front 1 (Rank 1): Solutions dominated only by Front 0
//   - Front 2 (Rank 2): Solutions dominated only by Fronts 0 and 1
//   - And so on...
//
// Algorithm steps:
//  1. For each solution p, calculate:
//     - S_p: Set of solutions that p dominates
//     - N_p: Number of solutions that dominate p
//  2. Front 0 = all solutions with N_p = 0
//  3. For each front F_i, create next front F_(i+1):
//     - For each p in F_i, for each q in S_p, decrement N_q
//     - If N_q becomes 0, add q to F_(i+1)
//
// Time complexity: O(M * N^2) where M is number of objectives, N is population size
//
// Returns: Map of rank -> solutions in that rank (deep copies)
func FastNonDominatedRanking(
	ctx context.Context, elems []models.Vector,
) map[int][]models.Vector {
	tracer := otel.Tracer("de")
	ctx, span := tracer.Start(ctx, "de.FastNonDominatedRanking",
		trace.WithAttributes(
			attribute.Int("population_size", len(elems)),
		),
	)
	defer span.End()

	if len(elems) > 0 {
		span.SetAttributes(attribute.Int("objectives_count", len(elems[0].Objectives)))
	}

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

	span.SetAttributes(
		attribute.Int("num_fronts", len(fronts)-1), // -1 because last front is always empty
		attribute.Int("front_0_size", len(rankedSubList[0])),
	)

	return rankedSubList
}

// DominanceTest determines Pareto dominance relationship between two objective vectors.
//
// In multi-objective optimization, solution x dominates solution y if:
//  1. x is no worse than y in all objectives (x[i] <= y[i] for all i)
//  2. x is strictly better than y in at least one objective (x[j] < y[j] for some j)
//
// Note: This implementation assumes MAXIMIZATION objectives (higher is better).
// For minimization problems, the comparison logic is inverted.
//
// Returns:
//   - -1: x dominates y (x is better)
//   - 1:  y dominates x (y is better)
//   - 0:  neither dominates (non-dominated, both are Pareto optimal relative to each other)
//
// Example:
//
//	x = [0.8, 0.2]  y = [0.6, 0.4]
//	DominanceTest(x, y) = 0  (x is better in first objective, y is better in second)
//
//	x = [0.8, 0.6]  y = [0.5, 0.4]
//	DominanceTest(x, y) = -1 (x dominates: better in both objectives)
func DominanceTest(x, y []float64) int {
	result := 0
	for i := range x {
		if x[i] > y[i] {
			// x is better in this objective
			if result == -1 {
				// But y was better in a previous objective, so neither dominates
				return 0
			}
			result = 1
		}
		if y[i] > x[i] {
			// y is better in this objective
			if result == 1 {
				// But x was better in a previous objective, so neither dominates
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
	tracer := otel.Tracer("de")
	_, span := tracer.Start(context.Background(), "de.FilterDominated",
		trace.WithAttributes(
			attribute.Int("population_size", len(elems)),
		),
	)
	defer span.End()

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

	span.SetAttributes(
		attribute.Int("non_dominated_count", len(nonDominated)),
		attribute.Int("dominated_count", len(dominated)),
	)

	return nonDominated, dominated
}

// CalculateCrwdDist calculates crowding distance for solutions in the same Pareto front.
//
// Crowding distance is a density estimation metric from NSGA-II that measures how close
// a solution is to its neighbors in objective space. It's used to maintain diversity
// in the population by preferring solutions in less crowded regions.
//
// Algorithm (for each objective m):
//  1. Sort solutions by objective m
//  2. Assign infinite distance to boundary solutions (best and worst)
//  3. For each interior solution i:
//     distance[i] += (obj[i+1] - obj[i-1]) / (obj_max - obj_min)
//
// The final crowding distance is the sum across all objectives, normalized by the
// range of each objective. Higher values indicate more "isolated" solutions.
//
// Formula for solution i in objective m:
//
//	distance[i][m] = (f_m(i+1) - f_m(i-1)) / (f_m_max - f_m_min)
//
// Where f_m(i) is the value of objective m for solution i after sorting.
//
// Special cases:
//   - Populations with ≤2 solutions: all get infinite distance
//   - Boundary solutions in each objective: get infinite distance
//   - Distance capped at INF to prevent overflow
//
// Time complexity: O(M * N * log N) where M is objectives, N is population size
//
// Modifies elems in-place by setting the CrowdingDistance field.
// Supports context cancellation for long-running calculations.
func CalculateCrwdDist(ctx context.Context, elems []models.Vector) error {
	tracer := otel.Tracer("de")
	ctx, span := tracer.Start(ctx, "de.CalculateCrwdDist",
		trace.WithAttributes(
			attribute.Int("population_size", len(elems)),
		),
	)
	defer span.End()

	if len(elems) <= 2 {
		for i := range elems {
			elems[i].CrowdingDistance = math.MaxFloat64
		}
		span.SetAttributes(attribute.Bool("small_population", true))
		return nil
	}

	// Check for cancellation before starting
	if err := ctx.Err(); err != nil {
		span.RecordError(err)
		return err
	}

	if len(elems) > 0 {
		span.SetAttributes(attribute.Int("objectives_count", len(elems[0].Objectives)))
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

