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

	// Validate population size (minimum 4 for mutation operations)
	if err := ValidateRange(cfg.PopulationSize, int64(4), int64(10000), "population_size"); err != nil {
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
