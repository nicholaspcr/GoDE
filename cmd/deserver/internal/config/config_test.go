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
	assert.Equal(t, "json", cfg.Log.Type)

	// Server defaults
	assert.NotEmpty(t, cfg.Server.LisAddr)
	assert.NotEmpty(t, cfg.Server.HTTPPort)

	// Store defaults
	assert.NotEmpty(t, cfg.Store.Redis.Host)
	assert.Greater(t, cfg.Store.Redis.Port, 0)
	assert.NotZero(t, cfg.Store.ExecutionTTL)
	assert.NotZero(t, cfg.Store.ResultTTL)
	assert.NotZero(t, cfg.Store.ProgressTTL)
}

func TestConfig_StringifyJSON(t *testing.T) {
	cfg := Default()

	t.Run("produces valid JSON", func(t *testing.T) {
		result, err := cfg.StringifyJSON()
		require.NoError(t, err)
		assert.NotEmpty(t, result)

		// Verify it's valid JSON by parsing it back
		var parsed map[string]any
		err = json.Unmarshal([]byte(result), &parsed)
		require.NoError(t, err)

		// Should have top-level keys
		assert.Contains(t, parsed, "log")
		assert.Contains(t, parsed, "store")
		assert.Contains(t, parsed, "server")
	})

	t.Run("is indented", func(t *testing.T) {
		result, err := cfg.StringifyJSON()
		require.NoError(t, err)
		assert.Contains(t, result, "    ") // 4-space indent
	})
}

func TestConfig_StringifyYAML(t *testing.T) {
	cfg := Default()

	t.Run("produces valid YAML", func(t *testing.T) {
		result, err := cfg.StringifyYAML()
		require.NoError(t, err)
		assert.NotEmpty(t, result)

		// Verify it's valid YAML by parsing it back
		var parsed map[string]any
		err = yaml.Unmarshal([]byte(result), &parsed)
		require.NoError(t, err)

		// Should have top-level keys
		assert.Contains(t, parsed, "log")
		assert.Contains(t, parsed, "store")
		assert.Contains(t, parsed, "server")
	})
}

func TestConfig_RoundTrip_JSON(t *testing.T) {
	cfg := Default()

	// Serialize to JSON
	jsonStr, err := cfg.StringifyJSON()
	require.NoError(t, err)

	// Deserialize back
	var roundTripped Config
	err = json.Unmarshal([]byte(jsonStr), &roundTripped)
	require.NoError(t, err)

	// Verify key fields survived the round trip
	assert.Equal(t, cfg.Log.Type, roundTripped.Log.Type)
	assert.Equal(t, cfg.Log.Level, roundTripped.Log.Level)
	assert.Equal(t, cfg.Store.Redis.Host, roundTripped.Store.Redis.Host)
	assert.Equal(t, cfg.Store.Redis.Port, roundTripped.Store.Redis.Port)
}

func TestConfig_RoundTrip_YAML(t *testing.T) {
	cfg := Default()

	// Serialize to YAML
	yamlStr, err := cfg.StringifyYAML()
	require.NoError(t, err)

	// Deserialize back
	var roundTripped Config
	err = yaml.Unmarshal([]byte(yamlStr), &roundTripped)
	require.NoError(t, err)

	// Verify key fields survived the round trip
	assert.Equal(t, cfg.Log.Type, roundTripped.Log.Type)
	assert.Equal(t, cfg.Log.Level, roundTripped.Log.Level)
	assert.Equal(t, cfg.Store.Redis.Host, roundTripped.Store.Redis.Host)
	assert.Equal(t, cfg.Store.Redis.Port, roundTripped.Store.Redis.Port)
}
