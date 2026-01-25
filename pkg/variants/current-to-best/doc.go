// Package currenttobest implements DE/current-to-best mutation strategies that blend current and best individuals.
//
// Mutation Formula:
//
// DE/current-to-best/1:
//
//	v_i = x_i + F(best - r1) + F(r2 - r3)
//
// where:
//   - x_i is the current (target) individual
//   - best is a randomly selected individual from the Pareto front (rank 0)
//   - r1, r2, r3 are three randomly selected individuals (different from target i)
//   - F is the mutation scaling factor (typically 0.5-1.0)
//   - v_i is the resulting mutant vector
//
// Minimum population: 4 (target + 3 random vectors)
//
// Alternative Interpretation:
// The formula can be rewritten to show the dual nature:
//
//	v_i = x_i + F(best - r1) + F(r2 - r3)
//	    = x_i + F·best - F·r1 + F·r2 - F·r3
//
// This shows:
//   - Attraction toward best individual (exploitation)
//   - Perturbation from random difference vectors (exploration)
//
// Characteristics:
//   - Balances exploration and exploitation
//   - Uses current position as base (more conservative than best/1)
//   - Directed search toward best solutions
//   - Maintains population diversity better than pure best variants
//   - Effective convergence with moderate exploration
//   - Good for multi-objective problems with complex Pareto fronts
package currenttobest
