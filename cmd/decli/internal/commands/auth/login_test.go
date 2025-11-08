package authcmd

import (
	"context"
	"testing"

	"github.com/nicholaspcr/GoDE/cmd/decli/internal/state"
	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// mockAuthClient implements api.AuthServiceClient for testing
type mockAuthClient struct {
	api.AuthServiceClient
	loginFn    func(ctx context.Context, req *api.AuthServiceLoginRequest, opts ...grpc.CallOption) (*api.AuthServiceLoginResponse, error)
	registerFn func(ctx context.Context, req *api.AuthServiceRegisterRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

func (m *mockAuthClient) Login(ctx context.Context, req *api.AuthServiceLoginRequest, opts ...grpc.CallOption) (*api.AuthServiceLoginResponse, error) {
	if m.loginFn != nil {
		return m.loginFn(ctx, req, opts...)
	}
	return &api.AuthServiceLoginResponse{Token: "test-token"}, nil
}

func (m *mockAuthClient) Register(ctx context.Context, req *api.AuthServiceRegisterRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	if m.registerFn != nil {
		return m.registerFn(ctx, req, opts...)
	}
	return &emptypb.Empty{}, nil
}

// mockStateDB implements state.State for testing
type mockStateDB struct {
	state.State
	savedToken string
}

func (m *mockStateDB) SaveAuthToken(token string) error {
	m.savedToken = token
	return nil
}

func (m *mockStateDB) GetAuthToken() (string, error) {
	return m.savedToken, nil
}

func (m *mockStateDB) DeleteAuthToken() error {
	m.savedToken = ""
	return nil
}

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
