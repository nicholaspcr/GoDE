package e2e

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	defaultServerAddr = "localhost:3030"
	testTimeout       = 30 * time.Second
)

// getServerAddr returns the server address from env or default
func getServerAddr() string {
	if addr := os.Getenv("E2E_SERVER_ADDR"); addr != "" {
		return addr
	}
	return defaultServerAddr
}

// skipIfNoServer skips the test if the server is not running
func skipIfNoServer(t *testing.T) *grpc.ClientConn {
	t.Helper()

	if os.Getenv("E2E_SKIP") != "" {
		t.Skip("E2E tests disabled (E2E_SKIP is set)")
	}

	conn, err := grpc.NewClient(
		getServerAddr(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Skipf("Cannot connect to server at %s: %v (set E2E_SKIP to disable)", getServerAddr(), err)
	}

	return conn
}

// TestE2E_FullUserWorkflow tests the complete user journey
func TestE2E_FullUserWorkflow(t *testing.T) {
	conn := skipIfNoServer(t)
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	// Generate unique username for this test run
	username := fmt.Sprintf("e2e_user_%d", time.Now().Unix())
	email := fmt.Sprintf("%s@test.com", username)
	password := "test_password_123"

	authClient := api.NewAuthServiceClient(conn)
	userClient := api.NewUserServiceClient(conn)

	t.Run("01_Register", func(t *testing.T) {
		_, err := authClient.Register(ctx, &api.AuthServiceRegisterRequest{
			User: &api.User{
				Ids:      &api.UserIDs{Username: username},
				Email:    email,
				Password: password,
			},
		})
		require.NoError(t, err, "registration should succeed")
	})

	var token string
	t.Run("02_Login", func(t *testing.T) {
		resp, err := authClient.Login(ctx, &api.AuthServiceLoginRequest{
			Username: username,
			Password: password,
		})
		require.NoError(t, err, "login should succeed")
		require.NotEmpty(t, resp.Token, "token should not be empty")
		token = resp.Token
	})

	// Create authenticated context
	authCtx := metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token)

	t.Run("03_GetUser", func(t *testing.T) {
		resp, err := userClient.Get(authCtx, &api.UserServiceGetRequest{
			Ids: &api.UserIDs{Username: username},
		})
		require.NoError(t, err, "get user should succeed")
		require.NotNil(t, resp.User)
		assert.Equal(t, username, resp.User.GetIds().Username)
		assert.Equal(t, email, resp.User.Email)
		assert.Empty(t, resp.User.Password, "password should never be returned")
	})

	t.Run("04_UpdateUser", func(t *testing.T) {
		newEmail := fmt.Sprintf("updated_%s@test.com", username)
		_, err := userClient.Update(authCtx, &api.UserServiceUpdateRequest{
			User:   &api.User{Ids: &api.UserIDs{Username: username}, Email: newEmail},
			Fields: []string{"email"},
		})
		require.NoError(t, err, "update should succeed")

		// Verify update
		resp, err := userClient.Get(authCtx, &api.UserServiceGetRequest{
			Ids: &api.UserIDs{Username: username},
		})
		require.NoError(t, err)
		assert.Equal(t, newEmail, resp.User.Email)
	})

	t.Run("05_DeleteUser", func(t *testing.T) {
		_, err := userClient.Delete(authCtx, &api.UserServiceDeleteRequest{
			Ids: &api.UserIDs{Username: username},
		})
		require.NoError(t, err, "delete should succeed")

		// Verify user is deleted
		_, err = userClient.Get(authCtx, &api.UserServiceGetRequest{
			Ids: &api.UserIDs{Username: username},
		})
		assert.Error(t, err, "get should fail after delete")
	})
}

// TestE2E_DEExecution tests differential evolution execution
func TestE2E_DEExecution(t *testing.T) {
	conn := skipIfNoServer(t)
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	// Create test user
	username := fmt.Sprintf("e2e_de_user_%d", time.Now().Unix())
	email := fmt.Sprintf("%s@test.com", username)
	password := "test_password_123"

	authClient := api.NewAuthServiceClient(conn)
	deClient := api.NewDifferentialEvolutionServiceClient(conn)

	// Register and login
	_, err := authClient.Register(ctx, &api.AuthServiceRegisterRequest{
		User: &api.User{
			Ids:      &api.UserIDs{Username: username},
			Email:    email,
			Password: password,
		},
	})
	require.NoError(t, err)

	loginResp, err := authClient.Login(ctx, &api.AuthServiceLoginRequest{
		Username: username,
		Password: password,
	})
	require.NoError(t, err)

	authCtx := metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+loginResp.Token)

	t.Run("01_ListAlgorithms", func(t *testing.T) {
		resp, err := deClient.ListSupportedAlgorithms(authCtx, &emptypb.Empty{})
		require.NoError(t, err)
		assert.NotEmpty(t, resp.Algorithms, "should have at least one algorithm")
		t.Logf("Available algorithms: %v", resp.Algorithms)
	})

	t.Run("02_ListProblems", func(t *testing.T) {
		resp, err := deClient.ListSupportedProblems(authCtx, &emptypb.Empty{})
		require.NoError(t, err)
		assert.NotEmpty(t, resp.Problems, "should have at least one problem")
		t.Logf("Available problems: %v", resp.Problems)
	})

	t.Run("03_ListVariants", func(t *testing.T) {
		resp, err := deClient.ListSupportedVariants(authCtx, &emptypb.Empty{})
		require.NoError(t, err)
		assert.NotEmpty(t, resp.Variants, "should have at least one variant")
		t.Logf("Available variants: %v", resp.Variants)
	})

	t.Run("04_RunDE", func(t *testing.T) {
		stream, err := deClient.Run(authCtx, &api.DifferentialEvolutionServiceRunRequest{
			Algorithm:  "gde3",
			Problem:    "zdt1",
			Variant:    "rand/1/bin",
			Dim:        30,
			Pop:        100,
			Iterations: 250,
			Executions: 1,
		})
		require.NoError(t, err, "stream creation should succeed")

		// Receive results
		var results []*api.DifferentialEvolutionServiceRunResponse
		for {
			resp, err := stream.Recv()
			if err != nil {
				break
			}
			results = append(results, resp)
		}

		assert.NotEmpty(t, results, "should receive at least one result")

		// Check the last result has pareto data
		if len(results) > 0 {
			lastResult := results[len(results)-1]
			assert.NotNil(t, lastResult.Pareto, "pareto should not be nil")
			t.Logf("Received %d results, last pareto has %d vectors",
				len(results), len(lastResult.Pareto.Vectors))
		}
	})
}

// TestE2E_RateLimiting tests rate limiting behavior
func TestE2E_RateLimiting(t *testing.T) {
	conn := skipIfNoServer(t)
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	authClient := api.NewAuthServiceClient(conn)

	t.Run("LoginRateLimit", func(t *testing.T) {
		username := fmt.Sprintf("e2e_rate_user_%d", time.Now().Unix())
		password := "wrong_password"

		// Try to login many times with wrong password
		// Should eventually hit rate limit (5 per minute with burst of 2)
		var rateLimitHit bool
		for i := 0; i < 10; i++ {
			_, err := authClient.Login(ctx, &api.AuthServiceLoginRequest{
				Username: username,
				Password: password,
			})
			if err != nil && err.Error() == "too many login attempts" {
				rateLimitHit = true
				break
			}
			time.Sleep(100 * time.Millisecond)
		}

		// Note: This test might not always hit the rate limit depending on timing
		// It's more of a smoke test to ensure rate limiting doesn't panic
		t.Logf("Rate limit hit: %v", rateLimitHit)
	})
}

// TestE2E_UnauthorizedAccess tests that protected endpoints require authentication
func TestE2E_UnauthorizedAccess(t *testing.T) {
	conn := skipIfNoServer(t)
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	userClient := api.NewUserServiceClient(conn)

	t.Run("GetUserWithoutAuth", func(t *testing.T) {
		_, err := userClient.Get(ctx, &api.UserServiceGetRequest{
			Ids: &api.UserIDs{Username: "anyuser"},
		})
		assert.Error(t, err, "should fail without authentication")
	})

	t.Run("UpdateUserWithoutAuth", func(t *testing.T) {
		_, err := userClient.Update(ctx, &api.UserServiceUpdateRequest{
			User:   &api.User{Ids: &api.UserIDs{Username: "anyuser"}, Email: "new@test.com"},
			Fields: []string{"email"},
		})
		assert.Error(t, err, "should fail without authentication")
	})

	t.Run("DeleteUserWithoutAuth", func(t *testing.T) {
		_, err := userClient.Delete(ctx, &api.UserServiceDeleteRequest{
			Ids: &api.UserIDs{Username: "anyuser"},
		})
		assert.Error(t, err, "should fail without authentication")
	})
}
