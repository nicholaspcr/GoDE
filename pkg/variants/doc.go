// Package variants provides Differential Evolution (DE) mutation strategies.
//
// # Overview
//
// DE variants control how new candidate solutions (mutant vectors) are generated
// during the evolutionary process. The choice of variant significantly affects
// the algorithm's exploration-exploitation balance and convergence behavior.
//
// # Available Variants
//
// Random-based variants (exploration-focused):
//   - rand/1: v = r1 + F(r2 - r3)                               [min pop: 4]
//   - rand/2: v = r1 + F(r2 - r3) + F(r4 - r5)                  [min pop: 6]
//
// Best-based variants (exploitation-focused):
//   - best/1: v = best + F(r1 - r2)                             [min pop: 3]
//   - best/2: v = best + F(r1 - r2) + F(r3 - r4)                [min pop: 5]
//
// Hybrid variants (balanced):
//   - pbest:            v = x + F(pbest - x) + F(r1 - r2)       [min pop: 3]
//   - current-to-best/1: v = x + F(best - r1) + F(r2 - r3)     [min pop: 4]
//
// # Parameters
//
// F (Mutation Scaling Factor):
//   - Range: typically [0.5, 1.0]
//   - Controls the amplification of difference vectors
//   - Higher F increases exploration, lower F increases exploitation
//
// P (Selection Parameter, pbest only):
//   - Range: typically [0.05, 0.20]
//   - Percentage of top population to select pbest from
//   - Lower P increases greediness (more exploitation)
//
// CR (Crossover Rate, not defined here but used in DE):
//   - Range: [0.0, 1.0]
//   - Controls how many components are inherited from mutant vs target
//
// # Variant Selection Guide
//
// Choose rand/1 or rand/2 when:
//   - Problem has many local optima
//   - Population diversity is critical
//   - You can afford slower convergence
//
// Choose best/1 or best/2 when:
//   - Fast convergence is needed
//   - Good solutions already discovered
//   - Problem landscape is relatively smooth
//
// Choose pbest or current-to-best/1 when:
//   - Balance between exploration and exploitation is needed
//   - Multi-objective optimization with complex Pareto fronts
//   - Adaptive behavior is desired
//
// # Implementation Details
//
// All variants implement the Interface which requires:
//   - Name() returns the variant identifier
//   - Mutate() generates a new mutant vector
//
// Mutate() receives:
//   - elems: the full population
//   - rankZero: individuals on the Pareto front (for multi-objective)
//   - params: Parameters struct with F, P, dimension, etc.
//
// Mutate() returns:
//   - A new Vector with mutated Elements
//   - An error if population is insufficient or vectors are invalid
//
// # Population Size Requirements
//
// Each variant has a minimum population size based on the number of
// random vectors it requires:
//   - rand/1, current-to-best/1: minimum 4 (target + 3 randoms)
//   - best/1, pbest: minimum 3 (target + 2 randoms)
//   - rand/2: minimum 6 (target + 5 randoms)
//   - best/2: minimum 5 (target + 4 randoms)
//
// These requirements are enforced by pkg/validation/de_config.go.
package variants
