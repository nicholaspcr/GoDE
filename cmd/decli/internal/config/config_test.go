package config

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestDefault(t *testing.T) {
	cfg := Default()
	require.NotNil(t, cfg)

	// Log defaults
	assert.Equal(t, "json", cfg.Log.Config.Type)

	// Server defaults
	assert.Equal(t, "localhost:3030", cfg.Server.GRPCAddr)
	assert.Equal(t, "http://localhost:8081", cfg.Server.HTTPAddr)

	// State defaults
	assert.Equal(t, "file", cfg.State.Provider)
	assert.NotEmpty(t, cfg.State.Filepath)
}

func TestConfig_StringifyJSON(t *testing.T) {
	cfg := Default()

	t.Run("produces valid JSON", func(t *testing.T) {
		result, err := cfg.StringifyJSON()
		require.NoError(t, err)
		assert.NotEmpty(t, result)

		var parsed map[string]any
		err = json.Unmarshal([]byte(result), &parsed)
		require.NoError(t, err)

		assert.Contains(t, parsed, "log")
		assert.Contains(t, parsed, "server")
		assert.Contains(t, parsed, "state")
		assert.Contains(t, parsed, "run")
	})

	t.Run("is indented", func(t *testing.T) {
		result, err := cfg.StringifyJSON()
		require.NoError(t, err)
		assert.Contains(t, result, "    ")
	})
}

func TestConfig_StringifyYAML(t *testing.T) {
	cfg := Default()

	t.Run("produces valid YAML", func(t *testing.T) {
		result, err := cfg.StringifyYAML()
		require.NoError(t, err)
		assert.NotEmpty(t, result)

		var parsed map[string]any
		err = yaml.Unmarshal([]byte(result), &parsed)
		require.NoError(t, err)

		assert.Contains(t, parsed, "log")
		assert.Contains(t, parsed, "server")
		assert.Contains(t, parsed, "state")
		assert.Contains(t, parsed, "run")
	})
}

func TestConfig_RoundTrip_JSON(t *testing.T) {
	cfg := Default()
	cfg.Run = RunConfig{
		Algorithm: "gde3",
		Variant:   "rand1",
		Problem:   "zdt1",
		DeConfig: DEConfig{
			Executions:     5,
			Generations:    100,
			PopulationSize: 50,
			DimensionsSize: 10,
			ObjectivesSize: 2,
			FloorLimiter:   0.0,
			CeilLimiter:    1.0,
			GDE3: GDE3Config{
				CR: 0.9,
				F:  0.5,
				P:  0.1,
			},
		},
	}

	jsonStr, err := cfg.StringifyJSON()
	require.NoError(t, err)

	var roundTripped Config
	err = json.Unmarshal([]byte(jsonStr), &roundTripped)
	require.NoError(t, err)

	assert.Equal(t, cfg.Run.Algorithm, roundTripped.Run.Algorithm)
	assert.Equal(t, cfg.Run.Variant, roundTripped.Run.Variant)
	assert.Equal(t, cfg.Run.Problem, roundTripped.Run.Problem)
	assert.Equal(t, cfg.Run.DeConfig.Executions, roundTripped.Run.DeConfig.Executions)
	assert.Equal(t, cfg.Run.DeConfig.Generations, roundTripped.Run.DeConfig.Generations)
	assert.Equal(t, cfg.Run.DeConfig.PopulationSize, roundTripped.Run.DeConfig.PopulationSize)
	assert.Equal(t, cfg.Run.DeConfig.GDE3.CR, roundTripped.Run.DeConfig.GDE3.CR)
	assert.Equal(t, cfg.Run.DeConfig.GDE3.F, roundTripped.Run.DeConfig.GDE3.F)
	assert.Equal(t, cfg.Server.GRPCAddr, roundTripped.Server.GRPCAddr)
	assert.Equal(t, cfg.Server.HTTPAddr, roundTripped.Server.HTTPAddr)
}

func TestConfig_RoundTrip_YAML(t *testing.T) {
	cfg := Default()

	yamlStr, err := cfg.StringifyYAML()
	require.NoError(t, err)

	var roundTripped Config
	err = yaml.Unmarshal([]byte(yamlStr), &roundTripped)
	require.NoError(t, err)

	assert.Equal(t, cfg.Server.GRPCAddr, roundTripped.Server.GRPCAddr)
	assert.Equal(t, cfg.Server.HTTPAddr, roundTripped.Server.HTTPAddr)
	assert.Equal(t, cfg.State.Provider, roundTripped.State.Provider)
}
