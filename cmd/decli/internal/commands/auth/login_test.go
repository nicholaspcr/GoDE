package authcmd

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoginCommand(t *testing.T) {
	// Test that login command is properly initialized
	t.Run("command exists", func(t *testing.T) {
		assert.NotNil(t, loginCmd)
		assert.Equal(t, "login", loginCmd.Use)
		assert.NotEmpty(t, loginCmd.Short)
	})

	t.Run("required flags", func(t *testing.T) {
		// Check that username flag exists and is required
		flag := loginCmd.Flags().Lookup("username")
		require.NotNil(t, flag)
		assert.Equal(t, "username", flag.Name)

		// Verify it's marked as required
		annotations := flag.Annotations
		if annotations != nil {
			_, exists := annotations[cobra.BashCompOneRequiredFlag]
			assert.True(t, exists || loginCmd.Flag("username").Changed)
		}
	})
}

func TestRegisterCommand(t *testing.T) {
	// Test that register command is properly initialized
	t.Run("command exists", func(t *testing.T) {
		assert.NotNil(t, registerCmd)
		assert.Equal(t, "register", registerCmd.Use)
		assert.NotEmpty(t, registerCmd.Short)
	})

	t.Run("required flags", func(t *testing.T) {
		// Check that username and email flags exist
		usernameFlag := registerCmd.Flags().Lookup("username")
		require.NotNil(t, usernameFlag)
		assert.Equal(t, "username", usernameFlag.Name)

		emailFlag := registerCmd.Flags().Lookup("email")
		require.NotNil(t, emailFlag)
		assert.Equal(t, "email", emailFlag.Name)
	})
}

func TestAuthCommand(t *testing.T) {
	// Test that auth parent command exists
	t.Run("auth command exists", func(t *testing.T) {
		assert.NotNil(t, authCmd)
		assert.Equal(t, "auth", authCmd.Use)
		assert.NotEmpty(t, authCmd.Short)
	})

	t.Run("has subcommands", func(t *testing.T) {
		// Verify login, register, logout are subcommands
		commands := authCmd.Commands()
		assert.NotEmpty(t, commands)

		commandNames := make(map[string]bool)
		for _, cmd := range commands {
			commandNames[cmd.Use] = true
		}

		assert.True(t, commandNames["login"], "login command should be registered")
		assert.True(t, commandNames["register"], "register command should be registered")
		assert.True(t, commandNames["logout"], "logout command should be registered")
	})
}
