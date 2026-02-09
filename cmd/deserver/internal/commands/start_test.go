package commands

import (
	"testing"

	deconfig "github.com/nicholaspcr/GoDE/cmd/deserver/internal/config"
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

	t.Run("fails with invalid server config", func(t *testing.T) {
		cfg = deconfig.Default()
		// Default config has empty JWTSecret which fails validation
		cfg.Server.JWTSecret = ""
		err := StartCmd.RunE(StartCmd, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid server configuration")
	})
}
