package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStartCommand(t *testing.T) {
	t.Run("command exists", func(t *testing.T) {
		assert.NotNil(t, StartCmd)
		assert.Equal(t, "start", StartCmd.Use)
		assert.NotEmpty(t, StartCmd.Short)
		assert.NotEmpty(t, StartCmd.Long)
	})

	t.Run("has aliases", func(t *testing.T) {
		assert.Contains(t, StartCmd.Aliases, "run")
	})

	t.Run("has RunE function", func(t *testing.T) {
		assert.NotNil(t, StartCmd.RunE)
	})
}
