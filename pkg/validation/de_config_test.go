package validation

import (
	"testing"

	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"github.com/stretchr/testify/assert"
)

func TestValidateDEConfig(t *testing.T) {
	tests := []struct {
		config  *api.DEConfig
		name    string
		wantErr bool
	}{
		{
			name: "valid config",
			config: &api.DEConfig{
				Executions:     10,
				Generations:    100,
				PopulationSize: 100,
				DimensionsSize: 30,
				ObjectivesSize:  2,
				FloorLimiter:   0.0,
				CeilLimiter:    1.0,
			},
			wantErr: false,
		},
		{
			name: "valid minimum values",
			config: &api.DEConfig{
				Executions:     1,
				Generations:    1,
				PopulationSize: 3, // minimum for best/1 and pbest variants
				DimensionsSize: 1,
				ObjectivesSize:  1,
				FloorLimiter:   -10.0,
				CeilLimiter:    10.0,
			},
			wantErr: false,
		},
		{
			name: "valid maximum values",
			config: &api.DEConfig{
				Executions:     100,
				Generations:    10000,
				PopulationSize: 10000,
				DimensionsSize: 1000,
				ObjectivesSize:  10,
				FloorLimiter:   0.0,
				CeilLimiter:    100.0,
			},
			wantErr: false,
		},
		{
			name:    "nil config",
			config:  nil,
			wantErr: true,
		},
		{
			name: "invalid executions (0)",
			config: &api.DEConfig{
				Executions:     0,
				Generations:    100,
				PopulationSize: 100,
				DimensionsSize: 30,
				ObjectivesSize:  2,
				FloorLimiter:   0.0,
				CeilLimiter:    1.0,
			},
			wantErr: true,
		},
		{
			name: "invalid executions (negative)",
			config: &api.DEConfig{
				Executions:     -1,
				Generations:    100,
				PopulationSize: 100,
				DimensionsSize: 30,
				ObjectivesSize:  2,
				FloorLimiter:   0.0,
				CeilLimiter:    1.0,
			},
			wantErr: true,
		},
		{
			name: "invalid executions (too large)",
			config: &api.DEConfig{
				Executions:     101,
				Generations:    100,
				PopulationSize: 100,
				DimensionsSize: 30,
				ObjectivesSize:  2,
				FloorLimiter:   0.0,
				CeilLimiter:    1.0,
			},
			wantErr: true,
		},
		{
			name: "invalid generations (0)",
			config: &api.DEConfig{
				Executions:     10,
				Generations:    0,
				PopulationSize: 100,
				DimensionsSize: 30,
				ObjectivesSize:  2,
				FloorLimiter:   0.0,
				CeilLimiter:    1.0,
			},
			wantErr: true,
		},
		{
			name: "invalid generations (too large)",
			config: &api.DEConfig{
				Executions:     10,
				Generations:    10001,
				PopulationSize: 100,
				DimensionsSize: 30,
				ObjectivesSize:  2,
				FloorLimiter:   0.0,
				CeilLimiter:    1.0,
			},
			wantErr: true,
		},
		{
			name: "invalid population size (too small)",
			config: &api.DEConfig{
				Executions:     10,
				Generations:    100,
				PopulationSize: 2, // less than minimum of 3
				DimensionsSize: 30,
				ObjectivesSize:  2,
				FloorLimiter:   0.0,
				CeilLimiter:    1.0,
			},
			wantErr: true,
		},
		{
			name: "invalid population size (too large)",
			config: &api.DEConfig{
				Executions:     10,
				Generations:    100,
				PopulationSize: 10001,
				DimensionsSize: 30,
				ObjectivesSize:  2,
				FloorLimiter:   0.0,
				CeilLimiter:    1.0,
			},
			wantErr: true,
		},
		{
			name: "invalid dimensions size (0)",
			config: &api.DEConfig{
				Executions:     10,
				Generations:    100,
				PopulationSize: 100,
				DimensionsSize: 0,
				ObjectivesSize:  2,
				FloorLimiter:   0.0,
				CeilLimiter:    1.0,
			},
			wantErr: true,
		},
		{
			name: "invalid dimensions size (too large)",
			config: &api.DEConfig{
				Executions:     10,
				Generations:    100,
				PopulationSize: 100,
				DimensionsSize: 1001,
				ObjectivesSize:  2,
				FloorLimiter:   0.0,
				CeilLimiter:    1.0,
			},
			wantErr: true,
		},
		{
			name: "invalid objectives size (0)",
			config: &api.DEConfig{
				Executions:     10,
				Generations:    100,
				PopulationSize: 100,
				DimensionsSize: 30,
				ObjectivesSize:  0,
				FloorLimiter:   0.0,
				CeilLimiter:    1.0,
			},
			wantErr: true,
		},
		{
			name: "invalid objectives size (too large)",
			config: &api.DEConfig{
				Executions:     10,
				Generations:    100,
				PopulationSize: 100,
				DimensionsSize: 30,
				ObjectivesSize:  11,
				FloorLimiter:   0.0,
				CeilLimiter:    1.0,
			},
			wantErr: true,
		},
		{
			name: "invalid floor >= ceil",
			config: &api.DEConfig{
				Executions:     10,
				Generations:    100,
				PopulationSize: 100,
				DimensionsSize: 30,
				ObjectivesSize:  2,
				FloorLimiter:   1.0,
				CeilLimiter:    1.0,
			},
			wantErr: true,
		},
		{
			name: "invalid floor > ceil",
			config: &api.DEConfig{
				Executions:     10,
				Generations:    100,
				PopulationSize: 100,
				DimensionsSize: 30,
				ObjectivesSize:  2,
				FloorLimiter:   2.0,
				CeilLimiter:    1.0,
			},
			wantErr: true,
		},
		{
			name: "invalid memory limit (dimensions × population exceeds max)",
			config: &api.DEConfig{
				Executions:     10,
				Generations:    100,
				PopulationSize: 10000,
				DimensionsSize: 1001, // 10000 × 1001 = 10,010,000 > 10,000,000
				ObjectivesSize:  2,
				FloorLimiter:   0.0,
				CeilLimiter:    1.0,
			},
			wantErr: true,
		},
		{
			name: "valid at memory limit boundary",
			config: &api.DEConfig{
				Executions:     10,
				Generations:    100,
				PopulationSize: 10000,
				DimensionsSize: 1000, // 10000 × 1000 = 10,000,000 (exactly at limit)
				ObjectivesSize:  2,
				FloorLimiter:   0.0,
				CeilLimiter:    1.0,
			},
			wantErr: false,
		},
		{
			name: "invalid memory limit (extreme case)",
			config: &api.DEConfig{
				Executions:     10,
				Generations:    100,
				PopulationSize: 10000,
				DimensionsSize: 5000, // 10000 × 5000 = 50,000,000 >> 10,000,000
				ObjectivesSize:  2,
				FloorLimiter:   0.0,
				CeilLimiter:    1.0,
			},
			wantErr: true,
		},
		{
			name: "valid with GDE3 config",
			config: &api.DEConfig{
				Executions:     10,
				Generations:    100,
				PopulationSize: 100,
				DimensionsSize: 30,
				ObjectivesSize:  2,
				FloorLimiter:   0.0,
				CeilLimiter:    1.0,
				AlgorithmConfig: &api.DEConfig_Gde3{
					Gde3: &api.GDE3Config{
						Cr: 0.5,
						F:  0.5,
						P:  0.5,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid with bad GDE3 config",
			config: &api.DEConfig{
				Executions:     10,
				Generations:    100,
				PopulationSize: 100,
				DimensionsSize: 30,
				ObjectivesSize:  2,
				FloorLimiter:   0.0,
				CeilLimiter:    1.0,
				AlgorithmConfig: &api.DEConfig_Gde3{
					Gde3: &api.GDE3Config{
						Cr: 1.5,
						F:  0.5,
						P:  0.5,
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDEConfig(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateGDE3Config(t *testing.T) {
	tests := []struct {
		config  *api.GDE3Config
		name    string
		wantErr bool
	}{
		{
			name: "valid config",
			config: &api.GDE3Config{
				Cr: 0.5,
				F:  0.5,
				P:  0.5,
			},
			wantErr: false,
		},
		{
			name: "valid minimum values",
			config: &api.GDE3Config{
				Cr: 0.0,
				F:  0.0,
				P:  0.0,
			},
			wantErr: false,
		},
		{
			name: "valid maximum CR and P",
			config: &api.GDE3Config{
				Cr: 1.0,
				F:  2.0,
				P:  1.0,
			},
			wantErr: false,
		},
		{
			name:    "nil config (allowed)",
			config:  nil,
			wantErr: false,
		},
		{
			name: "invalid CR > 1",
			config: &api.GDE3Config{
				Cr: 1.5,
				F:  0.5,
				P:  0.5,
			},
			wantErr: true,
		},
		{
			name: "invalid CR < 0",
			config: &api.GDE3Config{
				Cr: -0.1,
				F:  0.5,
				P:  0.5,
			},
			wantErr: true,
		},
		{
			name: "invalid F > 2",
			config: &api.GDE3Config{
				Cr: 0.5,
				F:  2.5,
				P:  0.5,
			},
			wantErr: true,
		},
		{
			name: "invalid F < 0",
			config: &api.GDE3Config{
				Cr: 0.5,
				F:  -0.1,
				P:  0.5,
			},
			wantErr: true,
		},
		{
			name: "invalid P > 1",
			config: &api.GDE3Config{
				Cr: 0.5,
				F:  0.5,
				P:  1.5,
			},
			wantErr: true,
		},
		{
			name: "invalid P < 0",
			config: &api.GDE3Config{
				Cr: 0.5,
				F:  0.5,
				P:  -0.1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGDE3Config(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateRange(t *testing.T) {
	tests := []struct {
		name    string
		field   string
		value   int64
		min     int64
		max     int64
		wantErr bool
	}{
		{"within range", "test", 5, 1, 10, false},
		{"exact min", "test", 1, 1, 10, false},
		{"exact max", "test", 10, 1, 10, false},
		{"below min", "test", 0, 1, 10, true},
		{"above max", "test", 11, 1, 10, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRange(tt.value, tt.min, tt.max, tt.field)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateRangeFloat(t *testing.T) {
	tests := []struct {
		name    string
		field   string
		value   float32
		min     float32
		max     float32
		wantErr bool
	}{
		{"within range", "test", 0.5, 0.0, 1.0, false},
		{"exact min", "test", 0.0, 0.0, 1.0, false},
		{"exact max", "test", 1.0, 0.0, 1.0, false},
		{"below min", "test", -0.1, 0.0, 1.0, true},
		{"above max", "test", 1.1, 0.0, 1.0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRange(tt.value, tt.min, tt.max, tt.field)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateRunAsyncRequest(t *testing.T) {
	validConfig := &api.DEConfig{
		Executions:     10,
		Generations:    100,
		PopulationSize: 100,
		DimensionsSize: 30,
		ObjectivesSize:  2,
		FloorLimiter:   0.0,
		CeilLimiter:    1.0,
	}

	tests := []struct {
		name      string
		algorithm string
		variant   string
		problem   string
		config    *api.DEConfig
		wantErr   bool
	}{
		{
			name:      "valid request with rand/1",
			algorithm: "gde3",
			variant:   "rand/1",
			problem:   "zdt1",
			config:    validConfig,
			wantErr:   false,
		},
		{
			name:      "valid request with rand/2 and sufficient population",
			algorithm: "gde3",
			variant:   "rand/2",
			problem:   "zdt1",
			config: &api.DEConfig{
				Executions:     10,
				Generations:    100,
				PopulationSize: 6, // minimum for rand/2
				DimensionsSize: 30,
				ObjectivesSize:  2,
				FloorLimiter:   0.0,
				CeilLimiter:    1.0,
			},
			wantErr: false,
		},
		{
			name:      "invalid request - rand/2 requires minimum 6 population",
			algorithm: "gde3",
			variant:   "rand/2",
			problem:   "zdt1",
			config: &api.DEConfig{
				Executions:     10,
				Generations:    100,
				PopulationSize: 4, // too small for rand/2
				DimensionsSize: 30,
				ObjectivesSize:  2,
				FloorLimiter:   0.0,
				CeilLimiter:    1.0,
			},
			wantErr: true,
		},
		{
			name:      "valid request with best/1 and population 3",
			algorithm: "gde3",
			variant:   "best/1",
			problem:   "zdt1",
			config: &api.DEConfig{
				Executions:     10,
				Generations:    100,
				PopulationSize: 3, // minimum for best/1
				DimensionsSize: 30,
				ObjectivesSize:  2,
				FloorLimiter:   0.0,
				CeilLimiter:    1.0,
			},
			wantErr: false,
		},
		{
			name:      "valid request with best/2 and population 5",
			algorithm: "gde3",
			variant:   "best/2",
			problem:   "zdt1",
			config: &api.DEConfig{
				Executions:     10,
				Generations:    100,
				PopulationSize: 5, // minimum for best/2
				DimensionsSize: 30,
				ObjectivesSize:  2,
				FloorLimiter:   0.0,
				CeilLimiter:    1.0,
			},
			wantErr: false,
		},
		{
			name:      "invalid request - best/2 requires minimum 5 population",
			algorithm: "gde3",
			variant:   "best/2",
			problem:   "zdt1",
			config: &api.DEConfig{
				Executions:     10,
				Generations:    100,
				PopulationSize: 4, // too small for best/2
				DimensionsSize: 30,
				ObjectivesSize:  2,
				FloorLimiter:   0.0,
				CeilLimiter:    1.0,
			},
			wantErr: true,
		},
		{
			name:      "valid request with pbest and population 3",
			algorithm: "gde3",
			variant:   "pbest",
			problem:   "zdt1",
			config: &api.DEConfig{
				Executions:     10,
				Generations:    100,
				PopulationSize: 3, // minimum for pbest
				DimensionsSize: 30,
				ObjectivesSize:  2,
				FloorLimiter:   0.0,
				CeilLimiter:    1.0,
			},
			wantErr: false,
		},
		{
			name:      "valid request with current-to-best/1 and population 4",
			algorithm: "gde3",
			variant:   "current-to-best/1",
			problem:   "zdt1",
			config: &api.DEConfig{
				Executions:     10,
				Generations:    100,
				PopulationSize: 4, // minimum for current-to-best/1
				DimensionsSize: 30,
				ObjectivesSize:  2,
				FloorLimiter:   0.0,
				CeilLimiter:    1.0,
			},
			wantErr: false,
		},
		{
			name:      "invalid request - DE config validation fails",
			algorithm: "gde3",
			variant:   "rand/1",
			problem:   "zdt1",
			config: &api.DEConfig{
				Executions:     -1, // invalid
				Generations:    100,
				PopulationSize: 100,
				DimensionsSize: 30,
				ObjectivesSize:  2,
				FloorLimiter:   0.0,
				CeilLimiter:    1.0,
			},
			wantErr: true,
		},
		{
			name:      "unknown variant defaults to minimum 4",
			algorithm: "gde3",
			variant:   "unknown-variant",
			problem:   "zdt1",
			config: &api.DEConfig{
				Executions:     10,
				Generations:    100,
				PopulationSize: 4,
				DimensionsSize: 30,
				ObjectivesSize:  2,
				FloorLimiter:   0.0,
				CeilLimiter:    1.0,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRunAsyncRequest(tt.algorithm, tt.variant, tt.problem, tt.config)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetMinPopulationForVariant(t *testing.T) {
	tests := []struct {
		variant  string
		expected int
	}{
		{"rand/1", 4},
		{"rand/2", 6},
		{"best/1", 3},
		{"best/2", 5},
		{"pbest", 3},
		{"current-to-best/1", 4},
		{"unknown-variant", 4}, // default fallback
		{"", 4},                // empty string fallback
	}

	for _, tt := range tests {
		t.Run(tt.variant, func(t *testing.T) {
			result := getMinPopulationForVariant(tt.variant)
			assert.Equal(t, tt.expected, result)
		})
	}
}
