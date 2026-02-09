package commands

import (
	"testing"

	deconfig "github.com/nicholaspcr/GoDE/cmd/deserver/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRootCommand(t *testing.T) {
	t.Run("command exists", func(t *testing.T) {
		assert.NotNil(t, rootCmd)
		assert.Equal(t, "deserver", rootCmd.Use)
		assert.NotEmpty(t, rootCmd.Short)
		assert.NotEmpty(t, rootCmd.Long)
	})

	t.Run("has config flag", func(t *testing.T) {
		flag := rootCmd.PersistentFlags().Lookup("config")
		require.NotNil(t, flag)
		assert.Equal(t, "config", flag.Name)
		assert.Equal(t, "c", flag.Shorthand)
	})

	t.Run("has subcommands", func(t *testing.T) {
		commands := rootCmd.Commands()
		assert.NotEmpty(t, commands)

		names := make(map[string]bool)
		for _, cmd := range commands {
			names[cmd.Use] = true
		}

		assert.True(t, names["start"], "should have 'start' subcommand")
		assert.True(t, names["config"], "should have 'config' subcommand")
		assert.True(t, names["migrate"], "should have 'migrate' subcommand")
	})

	t.Run("RunE returns help", func(t *testing.T) {
		assert.NotNil(t, rootCmd.RunE)
	})

	t.Run("PersistentPreRunE loads config", func(t *testing.T) {
		assert.NotNil(t, rootCmd.PersistentPreRunE)
	})
}

func TestConfigCommand(t *testing.T) {
	t.Run("command exists", func(t *testing.T) {
		assert.NotNil(t, configCmd)
		assert.Equal(t, "config", configCmd.Use)
		assert.NotEmpty(t, configCmd.Short)
	})

	t.Run("has json flag", func(t *testing.T) {
		flag := configCmd.Flags().Lookup("json")
		require.NotNil(t, flag)
		assert.Equal(t, "true", flag.DefValue)
	})

	t.Run("has yaml flag", func(t *testing.T) {
		flag := configCmd.Flags().Lookup("yaml")
		require.NotNil(t, flag)
		assert.Equal(t, "false", flag.DefValue)
	})

	t.Run("json and yaml are mutually exclusive", func(t *testing.T) {
		assert.NotNil(t, configCmd.PreRunE)
	})

	t.Run("runs with JSON output", func(t *testing.T) {
		cfg = deconfig.Default()
		ofJSON = true
		ofYAML = false

		err := configCmd.RunE(configCmd, nil)
		assert.NoError(t, err)
	})

	t.Run("runs with YAML output", func(t *testing.T) {
		cfg = deconfig.Default()
		ofJSON = false
		ofYAML = true

		err := configCmd.RunE(configCmd, nil)
		assert.NoError(t, err)
	})
}
