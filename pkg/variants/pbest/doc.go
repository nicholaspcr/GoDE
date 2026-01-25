// Package pbest implements DE/pbest mutation strategies using top-performing individuals.
//
// Mutation Formula:
//
// DE/pbest:
//
//	v_i = x_i + F(pbest - x_i) + F(r1 - r2)
//
// where:
//   - x_i is the current (target) individual
//   - pbest is a randomly selected individual from the top P% of the population
//   - r1, r2 are two randomly selected individuals (different from target i)
//   - F is the mutation scaling factor (typically 0.5-1.0)
//   - P is the selection parameter (typically 0.05-0.20, meaning top 5-20%)
//   - v_i is the resulting mutant vector
//
// Minimum population: 3 (target + 2 random vectors)
//
// Selection Process:
//   - Calculate indexLimit = ⌈population_size × P⌉
//   - Select pbest randomly from rankZero[0:indexLimit]
//   - This ensures selection from top-performing individuals while maintaining diversity
//
// Characteristics:
//   - Balances exploration and exploitation
//   - More adaptive than best/1 due to probabilistic selection
//   - P parameter controls exploitation pressure (lower P = more greedy)
//   - Maintains diversity better than best/1
//   - Effective for multi-objective optimization
package pbest
