package decmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListAlgorithmsCommand(t *testing.T) {
	t.Run("command exists", func(t *testing.T) {
		assert.NotNil(t, listAlgorithmsCmd)
		assert.Equal(t, "list-algorithms", listAlgorithmsCmd.Use)
		assert.NotEmpty(t, listAlgorithmsCmd.Short)
		assert.Contains(t, listAlgorithmsCmd.Short, "algorithm")
	})
}

func TestListProblemsCommand(t *testing.T) {
	t.Run("command exists", func(t *testing.T) {
		assert.NotNil(t, listProblemsCmd)
		assert.Equal(t, "list-problems", listProblemsCmd.Use)
		assert.NotEmpty(t, listProblemsCmd.Short)
		assert.Contains(t, listProblemsCmd.Short, "problem")
	})
}

func TestListVariantsCommand(t *testing.T) {
	t.Run("command exists", func(t *testing.T) {
		assert.NotNil(t, listVariantsCmd)
		assert.Equal(t, "list-variants", listVariantsCmd.Use)
		assert.NotEmpty(t, listVariantsCmd.Short)
		assert.Contains(t, listVariantsCmd.Short, "variant")
	})
}

func TestDECommand(t *testing.T) {
	t.Run("de command exists", func(t *testing.T) {
		assert.NotNil(t, deCmd)
		assert.Equal(t, "de", deCmd.Use)
		assert.NotEmpty(t, deCmd.Short)
	})

	t.Run("has subcommands", func(t *testing.T) {
		commands := deCmd.Commands()
		assert.NotEmpty(t, commands)

		commandNames := make(map[string]bool)
		for _, cmd := range commands {
			commandNames[cmd.Use] = true
		}

		assert.True(t, commandNames["list-algorithms"], "list-algorithms should be registered")
		assert.True(t, commandNames["list-problems"], "list-problems should be registered")
		assert.True(t, commandNames["list-variants"], "list-variants should be registered")
		assert.True(t, commandNames["run"], "run should be registered")
	})
}
