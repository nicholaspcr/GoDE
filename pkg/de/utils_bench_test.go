package de

import (
	"context"
	"testing"

	"github.com/nicholaspcr/GoDE/pkg/models"
)

// BenchmarkVectorCopy benchmarks different methods of copying vectors
func BenchmarkVectorCopy(b *testing.B) {
	// Create test vectors with realistic sizes
	sizes := []struct {
		name       string
		dimensions int
		objectives int
	}{
		{"small_2D_2Obj", 2, 2},
		{"medium_10D_2Obj", 10, 2},
		{"large_30D_3Obj", 30, 3},
		{"xlarge_100D_5Obj", 100, 5},
	}

	for _, size := range sizes {
		b.Run(size.name, func(b *testing.B) {
			// Create test data
			vectors := make([]models.Vector, 100)
			for i := range vectors {
				vectors[i] = models.Vector{
					Elements:         make([]float64, size.dimensions),
					Objectives:       make([]float64, size.objectives),
					CrowdingDistance: float64(i),
				}
				for j := range vectors[i].Elements {
					vectors[i].Elements[j] = float64(i + j)
				}
				for j := range vectors[i].Objectives {
					vectors[i].Objectives[j] = float64(i * j)
				}
			}

			b.Run("Copy_method", func(b *testing.B) {
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					result := make([]models.Vector, len(vectors))
					for idx, v := range vectors {
						result[idx] = v.Copy()
					}
				}
			})

			b.Run("builtin_copy", func(b *testing.B) {
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					result := make([]models.Vector, len(vectors))
					copy(result, vectors)
				}
			})
		})
	}
}

// BenchmarkReduceByCrowdDistance benchmarks the main function containing the TODO
func BenchmarkReduceByCrowdDistance(b *testing.B) {
	ctx := context.Background()

	// Create realistic population sizes
	populations := []struct {
		name     string
		popSize  int
		dimSize  int
		objSize  int
		reduceNP int
	}{
		{"small_10pop_5dim", 10, 5, 2, 5},
		{"medium_50pop_10dim", 50, 10, 2, 25},
		{"large_100pop_30dim", 100, 30, 3, 50},
	}

	for _, pop := range populations {
		b.Run(pop.name, func(b *testing.B) {
			// Generate test population
			elems := make([]models.Vector, pop.popSize)
			for i := range elems {
				elems[i] = models.Vector{
					Elements:   make([]float64, pop.dimSize),
					Objectives: make([]float64, pop.objSize),
				}
				// Create diverse objectives for realistic ranking
				for j := range elems[i].Objectives {
					elems[i].Objectives[j] = float64(i%10) + float64(j)*0.1
				}
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				// Make a copy for each iteration to avoid mutation
				testElems := make([]models.Vector, len(elems))
				for idx, v := range elems {
					testElems[idx] = v.Copy()
				}
				_, _ = ReduceByCrowdDistance(ctx, testElems, pop.reduceNP)
			}
		})
	}
}

// BenchmarkFastNonDominatedRanking benchmarks the ranking algorithm
func BenchmarkFastNonDominatedRanking(b *testing.B) {
	ctx := context.Background()

	populations := []struct {
		name    string
		popSize int
		objSize int
	}{
		{"10_individuals_2obj", 10, 2},
		{"50_individuals_2obj", 50, 2},
		{"100_individuals_3obj", 100, 3},
	}

	for _, pop := range populations {
		b.Run(pop.name, func(b *testing.B) {
			// Generate test population
			elems := make([]models.Vector, pop.popSize)
			for i := range elems {
				elems[i] = models.Vector{
					Objectives: make([]float64, pop.objSize),
				}
				for j := range elems[i].Objectives {
					elems[i].Objectives[j] = float64(i%10) + float64(j)*0.1
				}
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = FastNonDominatedRanking(ctx, elems)
			}
		})
	}
}

// BenchmarkCalculateCrwdDist benchmarks crowding distance calculation
func BenchmarkCalculateCrwdDist(b *testing.B) {
	sizes := []struct {
		name    string
		popSize int
		objSize int
	}{
		{"10_individuals", 10, 2},
		{"50_individuals", 50, 2},
		{"100_individuals", 100, 3},
	}

	for _, size := range sizes {
		b.Run(size.name, func(b *testing.B) {
			elems := make([]models.Vector, size.popSize)
			for i := range elems {
				elems[i] = models.Vector{
					Objectives: make([]float64, size.objSize),
				}
				for j := range elems[i].Objectives {
					elems[i].Objectives[j] = float64(i) + float64(j)*0.1
				}
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				CalculateCrwdDist(elems)
			}
		})
	}
}

// BenchmarkDominanceTest benchmarks the dominance comparison
func BenchmarkDominanceTest(b *testing.B) {
	x := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	y := []float64{1.5, 2.5, 2.5, 4.5, 4.5}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = DominanceTest(x, y)
	}
}
