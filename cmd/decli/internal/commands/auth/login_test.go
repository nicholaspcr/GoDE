package authcmd

import (
	"testing"

	"github.com/nicholaspcr/GoDE/cmd/decli/internal/config"
	"github.com/nicholaspcr/GoDE/cmd/decli/internal/state"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockStateOps implements state.Operations for testing.
type mockStateOps struct {
	token string
	err   error
}

func (m *mockStateOps) GetAuthToken() (string, error) { return m.token, m.err }
func (m *mockStateOps) InvalidateAuthToken() error     { return m.err }
func (m *mockStateOps) SaveAuthToken(t string) error {
	m.token = t
	return m.err
}

func TestLoginCommand(t *testing.T) {
	t.Run("command exists", func(t *testing.T) {
		assert.NotNil(t, loginCmd)
		assert.Equal(t, "login", loginCmd.Use)
		assert.NotEmpty(t, loginCmd.Short)
	})

	t.Run("has RunE function", func(t *testing.T) {
		assert.NotNil(t, loginCmd.RunE)
	})

	t.Run("required flags", func(t *testing.T) {
		flag := loginCmd.Flags().Lookup("username")
		require.NotNil(t, flag)
		assert.Equal(t, "username", flag.Name)

		annotations := flag.Annotations
		if annotations != nil {
			_, exists := annotations[cobra.BashCompOneRequiredFlag]
			assert.True(t, exists || loginCmd.Flag("username").Changed)
		}
	})

	t.Run("has password flag", func(t *testing.T) {
		flag := loginCmd.Flags().Lookup("password")
		require.NotNil(t, flag)
		assert.Equal(t, "password", flag.Name)
		assert.Equal(t, "", flag.DefValue)
	})
}

func TestRegisterCommand(t *testing.T) {
	t.Run("command exists", func(t *testing.T) {
		assert.NotNil(t, registerCmd)
		assert.Equal(t, "register", registerCmd.Use)
		assert.NotEmpty(t, registerCmd.Short)
	})

	t.Run("has RunE function", func(t *testing.T) {
		assert.NotNil(t, registerCmd.RunE)
	})

	t.Run("required flags", func(t *testing.T) {
		usernameFlag := registerCmd.Flags().Lookup("username")
		require.NotNil(t, usernameFlag)
		assert.Equal(t, "username", usernameFlag.Name)

		emailFlag := registerCmd.Flags().Lookup("email")
		require.NotNil(t, emailFlag)
		assert.Equal(t, "email", emailFlag.Name)
	})

	t.Run("has password flag", func(t *testing.T) {
		flag := registerCmd.Flags().Lookup("password")
		require.NotNil(t, flag)
		assert.Equal(t, "password", flag.Name)
		assert.Equal(t, "", flag.DefValue)
	})
}

func TestLogoutCommand(t *testing.T) {
	t.Run("command exists", func(t *testing.T) {
		assert.NotNil(t, logoutCmd)
		assert.Equal(t, "logout", logoutCmd.Use)
		assert.NotEmpty(t, logoutCmd.Short)
	})

	t.Run("has RunE function", func(t *testing.T) {
		assert.NotNil(t, logoutCmd.RunE)
	})

	t.Run("executes with mock state", func(t *testing.T) {
		db = &mockStateOps{}
		err := logoutCmd.RunE(logoutCmd, nil)
		assert.NoError(t, err)
	})
}

func TestAuthCommand(t *testing.T) {
	t.Run("auth command exists", func(t *testing.T) {
		assert.NotNil(t, authCmd)
		assert.Equal(t, "auth", authCmd.Use)
		assert.NotEmpty(t, authCmd.Short)
	})

	t.Run("has RunE that returns help", func(t *testing.T) {
		assert.NotNil(t, authCmd.RunE)
	})

	t.Run("has subcommands", func(t *testing.T) {
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

func TestRegisterCommands(t *testing.T) {
	root := &cobra.Command{Use: "test"}
	RegisterCommands(root)

	commands := root.Commands()
	assert.Len(t, commands, 1)
	assert.Equal(t, "auth", commands[0].Use)
}

func TestSetupConfig(t *testing.T) {
	testCfg := config.Default()
	SetupConfig(testCfg)
	assert.Equal(t, testCfg, cfg)
}

func TestSetupStateHandler(t *testing.T) {
	mock := &mockStateOps{}
	SetupStateHandler(mock)
	assert.Equal(t, state.Operations(mock), db)
}
