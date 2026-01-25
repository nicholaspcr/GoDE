// Package best implements DE/best mutation strategies using the best individual in the population.
//
// Mutation Formulas:
//
// DE/best/1:
//
//	v_i = best + F(r1 - r2)
//
// where:
//   - best is a randomly selected individual from the Pareto front (rank 0)
//   - r1, r2 are two randomly selected individuals (different from target i)
//   - F is the mutation scaling factor (typically 0.5-1.0)
//   - v_i is the resulting mutant vector
//
// Minimum population: 3 (target + 2 random vectors)
//
// DE/best/2:
//
//	v_i = best + F(r1 - r2) + F(r3 - r4)
//
// where:
//   - best is a randomly selected individual from the Pareto front (rank 0)
//   - r1, r2, r3, r4 are four randomly selected individuals (different from target i)
//   - F is the mutation scaling factor
//   - v_i is the resulting mutant vector
//
// Minimum population: 5 (target + 4 random vectors)
//
// Characteristics:
//   - Highly exploitative due to best base vector
//   - Fast convergence speed
//   - best/2 provides more diversity than best/1
//   - Risk of premature convergence if population loses diversity
//   - Effective when good solutions are already discovered
package best
