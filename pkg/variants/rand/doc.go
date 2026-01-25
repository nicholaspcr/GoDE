// Package rand implements DE/rand mutation strategies using randomly selected individuals.
//
// Mutation Formulas:
//
// DE/rand/1:
//
//	v_i = r1 + F(r2 - r3)
//
// where:
//   - r1, r2, r3 are three randomly selected individuals (all different from target i)
//   - F is the mutation scaling factor (typically 0.5-1.0)
//   - v_i is the resulting mutant vector
//
// Minimum population: 4 (target + 3 random vectors)
//
// DE/rand/2:
//
//	v_i = r1 + F(r2 - r3) + F(r4 - r5)
//
// where:
//   - r1, r2, r3, r4, r5 are five randomly selected individuals (all different from target i)
//   - F is the mutation scaling factor
//   - v_i is the resulting mutant vector
//
// Minimum population: 6 (target + 5 random vectors)
//
// Characteristics:
//   - Highly explorative due to random base vector
//   - Good for avoiding premature convergence
//   - rand/2 provides stronger diversity than rand/1
//   - Convergence speed typically slower than best variants
package rand
