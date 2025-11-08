package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMigrateCommand(t *testing.T) {
	t.Run("parent command exists", func(t *testing.T) {
		assert.NotNil(t, migrateCmd)
		assert.Equal(t, "migrate", migrateCmd.Use)
		assert.NotEmpty(t, migrateCmd.Short)
		assert.NotEmpty(t, migrateCmd.Long)
	})

	t.Run("has subcommands", func(t *testing.T) {
		commands := migrateCmd.Commands()
		assert.NotEmpty(t, commands)

		commandNames := make(map[string]bool)
		for _, cmd := range commands {
			commandNames[cmd.Use] = true
		}

		assert.True(t, commandNames["up"], "up command should be registered")
		assert.True(t, commandNames["down"], "down command should be registered")
		assert.True(t, commandNames["version"], "version command should be registered")
	})
}

func TestMigrateUpCommand(t *testing.T) {
	t.Run("command exists", func(t *testing.T) {
		assert.NotNil(t, migrateUpCmd)
		assert.Equal(t, "up", migrateUpCmd.Use)
		assert.NotEmpty(t, migrateUpCmd.Short)
		assert.NotEmpty(t, migrateUpCmd.Long)
	})

	t.Run("has RunE function", func(t *testing.T) {
		assert.NotNil(t, migrateUpCmd.RunE)
	})
}

func TestMigrateDownCommand(t *testing.T) {
	t.Run("command exists", func(t *testing.T) {
		assert.NotNil(t, migrateDownCmd)
		assert.Equal(t, "down", migrateDownCmd.Use)
		assert.NotEmpty(t, migrateDownCmd.Short)
		assert.NotEmpty(t, migrateDownCmd.Long)
	})

	t.Run("has RunE function", func(t *testing.T) {
		assert.NotNil(t, migrateDownCmd.RunE)
	})

	t.Run("has steps flag", func(t *testing.T) {
		flag := migrateDownCmd.Flags().Lookup("steps")
		require.NotNil(t, flag)
		assert.Equal(t, "steps", flag.Name)
		assert.Equal(t, "n", flag.Shorthand)
		assert.Equal(t, "1", flag.DefValue)
	})
}

func TestMigrateVersionCommand(t *testing.T) {
	t.Run("command exists", func(t *testing.T) {
		assert.NotNil(t, migrateVersionCmd)
		assert.Equal(t, "version", migrateVersionCmd.Use)
		assert.NotEmpty(t, migrateVersionCmd.Short)
		assert.NotEmpty(t, migrateVersionCmd.Long)
	})

	t.Run("has RunE function", func(t *testing.T) {
		assert.NotNil(t, migrateVersionCmd.RunE)
	})
}
