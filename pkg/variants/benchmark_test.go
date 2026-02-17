package variants_test

import (
	"math/rand"
	"testing"

	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/variants"
	"github.com/nicholaspcr/GoDE/pkg/variants/best"
	currentToBest "github.com/nicholaspcr/GoDE/pkg/variants/current-to-best"
	"github.com/nicholaspcr/GoDE/pkg/variants/pbest"
	variantsRand "github.com/nicholaspcr/GoDE/pkg/variants/rand"
)

// Benchmark configurations
var (
	benchDimensions   = []int{10, 30, 100}
	benchPopulations  = []int{50, 100, 200}
	benchObjectives   = 3
	benchRankZeroSize = 10
)

// setupBenchmark creates test data for benchmarking variants.
func setupBenchmark(dim, popSize, rankZeroSize int) ([]models.Vector, []models.Vector, variants.Parameters) {
	rng := rand.New(rand.NewSource(42)) //nolint:gosec // Deterministic for benchmarks

	// Create population
	elems := make([]models.Vector, popSize)
	for i := range popSize {
		elements := make([]float64, dim)
		objectives := make([]float64, benchObjectives)
		for j := range dim {
			elements[j] = rng.Float64()
		}
		for j := range benchObjectives {
			objectives[j] = rng.Float64()
		}
		elems[i] = models.Vector{
			Elements:   elements,
			Objectives: objectives,
		}
	}

	// Create rank zero (best individuals)
	rankZero := make([]models.Vector, rankZeroSize)
	for i := range rankZeroSize {
		elements := make([]float64, dim)
		objectives := make([]float64, benchObjectives)
		for j := range dim {
			elements[j] = rng.Float64()
		}
		for j := range benchObjectives {
			objectives[j] = rng.Float64() * 0.5 // Better objectives
		}
		rankZero[i] = models.Vector{
			Elements:   elements,
			Objectives: objectives,
		}
	}

	params := variants.Parameters{
		F:       0.5,
		P:       0.1,
		DIM:     dim,
		CurrPos: 0,
		Random:  rng,
	}

	return elems, rankZero, params
}

// Benchmark rand/1 variant
func BenchmarkRand1(b *testing.B) {
	variant := variantsRand.Rand1()

	for _, dim := range benchDimensions {
		for _, popSize := range benchPopulations {
			elems, rankZero, params := setupBenchmark(dim, popSize, benchRankZeroSize)

			b.Run(benchName("rand/1", dim, popSize), func(b *testing.B) {
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					params.CurrPos = i % popSize
					_, err := variant.Mutate(elems, rankZero, params)
					if err != nil {
						b.Fatal(err)
					}
				}
			})
		}
	}
}

// Benchmark rand/2 variant
func BenchmarkRand2(b *testing.B) {
	variant := variantsRand.Rand2()

	for _, dim := range benchDimensions {
		for _, popSize := range benchPopulations {
			elems, rankZero, params := setupBenchmark(dim, popSize, benchRankZeroSize)

			b.Run(benchName("rand/2", dim, popSize), func(b *testing.B) {
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					params.CurrPos = i % popSize
					_, err := variant.Mutate(elems, rankZero, params)
					if err != nil {
						b.Fatal(err)
					}
				}
			})
		}
	}
}

// Benchmark best/1 variant
func BenchmarkBest1(b *testing.B) {
	variant := best.Best1()

	for _, dim := range benchDimensions {
		for _, popSize := range benchPopulations {
			elems, rankZero, params := setupBenchmark(dim, popSize, benchRankZeroSize)

			b.Run(benchName("best/1", dim, popSize), func(b *testing.B) {
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					params.CurrPos = i % popSize
					_, err := variant.Mutate(elems, rankZero, params)
					if err != nil {
						b.Fatal(err)
					}
				}
			})
		}
	}
}

// Benchmark best/2 variant
func BenchmarkBest2(b *testing.B) {
	variant := best.Best2()

	for _, dim := range benchDimensions {
		for _, popSize := range benchPopulations {
			elems, rankZero, params := setupBenchmark(dim, popSize, benchRankZeroSize)

			b.Run(benchName("best/2", dim, popSize), func(b *testing.B) {
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					params.CurrPos = i % popSize
					_, err := variant.Mutate(elems, rankZero, params)
					if err != nil {
						b.Fatal(err)
					}
				}
			})
		}
	}
}

// Benchmark current-to-best/1 variant
func BenchmarkCurrentToBest1(b *testing.B) {
	variant := currentToBest.CurrToBest1()

	for _, dim := range benchDimensions {
		for _, popSize := range benchPopulations {
			elems, rankZero, params := setupBenchmark(dim, popSize, benchRankZeroSize)

			b.Run(benchName("current-to-best/1", dim, popSize), func(b *testing.B) {
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					params.CurrPos = i % popSize
					_, err := variant.Mutate(elems, rankZero, params)
					if err != nil {
						b.Fatal(err)
					}
				}
			})
		}
	}
}

// Benchmark pbest variant
func BenchmarkPBest(b *testing.B) {
	variant := pbest.Pbest()

	for _, dim := range benchDimensions {
		for _, popSize := range benchPopulations {
			elems, rankZero, params := setupBenchmark(dim, popSize, benchRankZeroSize)

			b.Run(benchName("pbest", dim, popSize), func(b *testing.B) {
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					params.CurrPos = i % popSize
					_, err := variant.Mutate(elems, rankZero, params)
					if err != nil {
						b.Fatal(err)
					}
				}
			})
		}
	}
}

// Benchmark memory allocations
func BenchmarkVariantMemory(b *testing.B) {
	variants := map[string]variants.Interface{
		"rand/1":            variantsRand.Rand1(),
		"rand/2":            variantsRand.Rand2(),
		"best/1":            best.Best1(),
		"best/2":            best.Best2(),
		"current-to-best/1": currentToBest.CurrToBest1(),
		"pbest":             pbest.Pbest(),
	}

	dim := 30
	popSize := 100
	elems, rankZero, params := setupBenchmark(dim, popSize, benchRankZeroSize)

	for name, variant := range variants {
		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				params.CurrPos = i % popSize
				_, err := variant.Mutate(elems, rankZero, params)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// Benchmark parallel execution
func BenchmarkVariantParallel(b *testing.B) {
	variants := map[string]variants.Interface{
		"rand/1":            variantsRand.Rand1(),
		"best/1":            best.Best1(),
		"current-to-best/1": currentToBest.CurrToBest1(),
		"pbest":             pbest.Pbest(),
	}

	dim := 30
	popSize := 100

	for name, variant := range variants {
		b.Run(name, func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				elems, rankZero, params := setupBenchmark(dim, popSize, benchRankZeroSize)
				i := 0
				for pb.Next() {
					params.CurrPos = i % popSize
					_, err := variant.Mutate(elems, rankZero, params)
					if err != nil {
						b.Fatal(err)
					}
					i++
				}
			})
		})
	}
}

// benchName creates a consistent benchmark name.
func benchName(variant string, dim, popSize int) string {
	return variant + "/dim=" + itoa(dim) + "/pop=" + itoa(popSize)
}

// itoa is a simple integer to string converter for benchmark names.
func itoa(i int) string {
	if i == 0 {
		return "0"
	}

	neg := i < 0
	if neg {
		i = -i
	}

	var buf [20]byte
	pos := len(buf)
	for i > 0 {
		pos--
		buf[pos] = byte(i%10) + '0'
		i /= 10
	}

	if neg {
		pos--
		buf[pos] = '-'
	}

	return string(buf[pos:])
}
