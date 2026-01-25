package validation

import (
	"fmt"

	"github.com/nicholaspcr/GoDE/pkg/api/v1"
)

// ValidateDEConfig validates differential evolution configuration.
func ValidateDEConfig(cfg *api.DEConfig) error {
	if cfg == nil {
		return NewValidationError("de_config", nil, ErrEmptyField, "DE config is nil")
	}

	// Validate executions
	if err := ValidateRange(cfg.Executions, int64(1), int64(100), "executions"); err != nil {
		return err
	}

	// Validate generations
	if err := ValidateRange(cfg.Generations, int64(1), int64(10000), "generations"); err != nil {
		return err
	}

	// Validate population size (minimum 3 - required by best/1 and pbest variants)
	// Variant-specific validation (e.g., rand/2 needs 6) is done in ValidateRunAsyncRequest
	if err := ValidateRange(cfg.PopulationSize, int64(3), int64(10000), "population_size"); err != nil {
		return err
	}

	// Validate dimensions size
	if err := ValidateRange(cfg.DimensionsSize, int64(1), int64(1000), "dimensions_size"); err != nil {
		return err
	}

	// Validate objectives size
	if err := ValidateRange(cfg.ObjectivesSize, int64(1), int64(10), "objectives_size"); err != nil {
		return err
	}

	// Validate floor and ceil limiters
	if cfg.FloorLimiter >= cfg.CeilLimiter {
		return NewValidationError(
			"floor_limiter",
			cfg.FloorLimiter,
			ErrOutOfRange,
			fmt.Sprintf("floor_limiter (%v) must be less than ceil_limiter (%v)",
				cfg.FloorLimiter, cfg.CeilLimiter),
		)
	}

	// Validate total memory allocation (dimensions × population)
	// Prevent memory bombs from excessively large configurations
	const maxTotalElements = int64(10_000_000) // 10 million elements max
	totalElements := cfg.DimensionsSize * cfg.PopulationSize
	if totalElements > maxTotalElements {
		return NewValidationError(
			"population_size",
			cfg.PopulationSize,
			ErrOutOfRange,
			fmt.Sprintf("dimensions_size (%d) × population_size (%d) = %d exceeds maximum allowed (%d)",
				cfg.DimensionsSize, cfg.PopulationSize, totalElements, maxTotalElements),
		)
	}

	// Validate GDE3 config if present
	if gde3 := cfg.GetGde3(); gde3 != nil {
		if err := ValidateGDE3Config(gde3); err != nil {
			return err
		}
	}

	return nil
}

// ValidateGDE3Config validates GDE3-specific parameters.
func ValidateGDE3Config(cfg *api.GDE3Config) error {
	if cfg == nil {
		return nil // GDE3 config is optional
	}

	// Validate CR (Crossover Rate): [0.0, 1.0]
	if err := ValidateRange(cfg.Cr, float32(0.0), float32(1.0), "cr"); err != nil {
		return err
	}

	// Validate F (Scaling Factor): [0.0, 2.0]
	if err := ValidateRange(cfg.F, float32(0.0), float32(2.0), "f"); err != nil {
		return err
	}

	// Validate P (Selection Parameter): [0.0, 1.0]
	if err := ValidateRange(cfg.P, float32(0.0), float32(1.0), "p"); err != nil {
		return err
	}

	return nil
}

// ValidateRunAsyncRequest validates a DE run request including variant-specific constraints.
func ValidateRunAsyncRequest(algorithm, variant, problem string, cfg *api.DEConfig) error {
	// Validate DE config first
	if err := ValidateDEConfig(cfg); err != nil {
		return err
	}

	// Validate variant-specific population size requirements
	minPopulation := getMinPopulationForVariant(variant)
	if cfg.PopulationSize < int64(minPopulation) {
		return NewValidationError(
			"population_size",
			cfg.PopulationSize,
			ErrOutOfRange,
			fmt.Sprintf("variant %s requires minimum population size of %d, got %d",
				variant, minPopulation, cfg.PopulationSize),
		)
	}

	return nil
}

// getMinPopulationForVariant returns the minimum population size required for a given variant.
//
// Different DE variants have different minimum population requirements based on the number
// of random vectors they need to select during mutation:
//   - rand/1: needs CurrPos + 3 random vectors = 4 minimum
//   - rand/2: needs CurrPos + 5 random vectors = 6 minimum
//   - best/1: needs CurrPos + 2 random vectors = 3 minimum
//   - best/2: needs CurrPos + 4 random vectors = 5 minimum
//   - pbest: needs CurrPos + 2 random vectors = 3 minimum
//   - current-to-best/1: needs CurrPos + 3 random vectors = 4 minimum
func getMinPopulationForVariant(variant string) int {
	switch variant {
	case "rand/1", "current-to-best/1":
		return 4
	case "rand/2":
		return 6
	case "best/1", "pbest":
		return 3
	case "best/2":
		return 5
	default:
		// Conservative fallback for unknown variants
		return 4
	}
}
