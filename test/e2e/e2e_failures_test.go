//go:build e2e

package e2e

import (
	"context"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// TestE2E_InvalidAuthentication tests authentication failure scenarios
func TestE2E_InvalidAuthentication(t *testing.T) {
	conn := setupConnection(t)
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	authClient := api.NewAuthServiceClient(conn)
	deClient := api.NewDifferentialEvolutionServiceClient(conn)

	t.Run("EmptyToken", func(t *testing.T) {
		authCtx := metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer ")
		_, err := deClient.ListSupportedAlgorithms(authCtx, &emptypb.Empty{})
		assert.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, st.Code())
	})

	t.Run("MalformedToken", func(t *testing.T) {
		authCtx := metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer invalid.token.here")
		_, err := deClient.ListSupportedAlgorithms(authCtx, &emptypb.Empty{})
		assert.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, st.Code())
	})

	t.Run("MissingBearerPrefix", func(t *testing.T) {
		authCtx := metadata.AppendToOutgoingContext(ctx, "authorization", "sometoken")
		_, err := deClient.ListSupportedAlgorithms(authCtx, &emptypb.Empty{})
		assert.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, st.Code())
	})

	t.Run("LoginWithNonexistentUser", func(t *testing.T) {
		_, err := authClient.Login(ctx, &api.AuthServiceLoginRequest{
			Username: "nonexistent_user_12345",
			Password: "password",
		})
		assert.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
	})

	t.Run("LoginWithWrongPassword", func(t *testing.T) {
		// First create a user
		username := fmt.Sprintf("e2e_wrongpass_%d", time.Now().Unix())
		password := "correct_password_123"

		_, err := authClient.Register(ctx, &api.AuthServiceRegisterRequest{
			User: &api.User{
				Ids:      &api.UserIDs{Username: username},
				Email:    fmt.Sprintf("%s@test.com", username),
				Password: password,
			},
		})
		require.NoError(t, err)

		// Try to login with wrong password
		_, err = authClient.Login(ctx, &api.AuthServiceLoginRequest{
			Username: username,
			Password: "wrong_password",
		})
		assert.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, st.Code())
	})
}

// TestE2E_InvalidDEConfiguration tests DE execution with invalid configurations
func TestE2E_InvalidDEConfiguration(t *testing.T) {
	conn := setupConnection(t)
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	// Setup authenticated user
	username := fmt.Sprintf("e2e_invalid_config_%d", time.Now().Unix())
	token := setupTestUser(t, ctx, conn, username)
	authCtx := metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token)

	deClient := api.NewDifferentialEvolutionServiceClient(conn)

	baseConfig := &api.DEConfig{
		Executions:     1,
		Generations:    250,
		PopulationSize: 100,
		DimensionsSize: 30,
		ObjectivesSize: 2,
		FloorLimiter:   0.0,
		CeilLimiter:    1.0,
		AlgorithmConfig: &api.DEConfig_Gde3{
			Gde3: &api.GDE3Config{
				Cr: 0.5,
				F:  0.5,
				P:  0.1,
			},
		},
	}

	t.Run("UnsupportedAlgorithm", func(t *testing.T) {
		_, err := deClient.RunAsync(authCtx, &api.RunAsyncRequest{
			Algorithm: "unsupported_algorithm",
			Problem:   "zdt1",
			Variant:   "rand1",
			DeConfig:  baseConfig,
		})
		assert.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
	})

	t.Run("UnsupportedProblem", func(t *testing.T) {
		_, err := deClient.RunAsync(authCtx, &api.RunAsyncRequest{
			Algorithm: "gde3",
			Problem:   "unsupported_problem",
			Variant:   "rand1",
			DeConfig:  baseConfig,
		})
		assert.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
	})

	t.Run("UnsupportedVariant", func(t *testing.T) {
		_, err := deClient.RunAsync(authCtx, &api.RunAsyncRequest{
			Algorithm: "gde3",
			Problem:   "zdt1",
			Variant:   "unsupported_variant",
			DeConfig:  baseConfig,
		})
		assert.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
	})

	t.Run("ZeroPopulationSize", func(t *testing.T) {
		invalidConfig := *baseConfig
		invalidConfig.PopulationSize = 0

		_, err := deClient.RunAsync(authCtx, &api.RunAsyncRequest{
			Algorithm: "gde3",
			Problem:   "zdt1",
			Variant:   "rand1",
			DeConfig:  &invalidConfig,
		})
		assert.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
	})

	t.Run("NegativeGenerations", func(t *testing.T) {
		invalidConfig := *baseConfig
		invalidConfig.Generations = -1

		_, err := deClient.RunAsync(authCtx, &api.RunAsyncRequest{
			Algorithm: "gde3",
			Problem:   "zdt1",
			Variant:   "rand1",
			DeConfig:  &invalidConfig,
		})
		assert.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
	})

	t.Run("InvalidFloorCeilRange", func(t *testing.T) {
		invalidConfig := *baseConfig
		invalidConfig.FloorLimiter = 1.0
		invalidConfig.CeilLimiter = 0.0 // Floor > Ceil

		_, err := deClient.RunAsync(authCtx, &api.RunAsyncRequest{
			Algorithm: "gde3",
			Problem:   "zdt1",
			Variant:   "rand1",
			DeConfig:  &invalidConfig,
		})
		assert.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
	})

	t.Run("InvalidCrossoverRate", func(t *testing.T) {
		invalidConfig := *baseConfig
		invalidConfig.AlgorithmConfig = &api.DEConfig_Gde3{
			Gde3: &api.GDE3Config{
				Cr: 1.5, // CR > 1.0
				F:  0.5,
				P:  0.1,
			},
		}

		_, err := deClient.RunAsync(authCtx, &api.RunAsyncRequest{
			Algorithm: "gde3",
			Problem:   "zdt1",
			Variant:   "rand1",
			DeConfig:  &invalidConfig,
		})
		assert.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
	})

	t.Run("PopulationTooSmallForVariant", func(t *testing.T) {
		invalidConfig := *baseConfig
		invalidConfig.PopulationSize = 3 // Too small for most variants

		_, err := deClient.RunAsync(authCtx, &api.RunAsyncRequest{
			Algorithm: "gde3",
			Problem:   "zdt1",
			Variant:   "best2",
			DeConfig:  &invalidConfig,
		})
		assert.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
	})
}

// TestE2E_AsyncExecutionLifecycle tests async execution scenarios
func TestE2E_AsyncExecutionLifecycle(t *testing.T) {
	conn := setupConnection(t)
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	// Setup authenticated user
	username := fmt.Sprintf("e2e_async_%d", time.Now().Unix())
	token := setupTestUser(t, ctx, conn, username)
	authCtx := metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token)

	deClient := api.NewDifferentialEvolutionServiceClient(conn)

	config := &api.DEConfig{
		Executions:     2,
		Generations:    100,
		PopulationSize: 50,
		DimensionsSize: 10,
		ObjectivesSize: 2,
		FloorLimiter:   0.0,
		CeilLimiter:    1.0,
		AlgorithmConfig: &api.DEConfig_Gde3{
			Gde3: &api.GDE3Config{
				Cr: 0.5,
				F:  0.5,
				P:  0.1,
			},
		},
	}

	t.Run("SubmitAndCancel", func(t *testing.T) {
		// Submit execution
		resp, err := deClient.RunAsync(authCtx, &api.RunAsyncRequest{
			Algorithm: "gde3",
			Problem:   "zdt1",
			Variant:   "rand1",
			DeConfig:  config,
		})
		require.NoError(t, err)
		executionID := resp.ExecutionId

		// Wait a bit to ensure execution starts
		time.Sleep(100 * time.Millisecond)

		// Cancel execution
		_, err = deClient.CancelExecution(authCtx, &api.CancelExecutionRequest{
			ExecutionId: executionID,
		})
		require.NoError(t, err)

		// Check status is cancelled
		time.Sleep(100 * time.Millisecond)
		statusResp, err := deClient.GetExecutionStatus(authCtx, &api.GetExecutionStatusRequest{
			ExecutionId: executionID,
		})
		require.NoError(t, err)
		assert.Equal(t, api.ExecutionStatus_EXECUTION_STATUS_CANCELLED, statusResp.Execution.Status)
	})

	t.Run("GetStatusOfNonexistentExecution", func(t *testing.T) {
		_, err := deClient.GetExecutionStatus(authCtx, &api.GetExecutionStatusRequest{
			ExecutionId: "nonexistent-execution-id",
		})
		assert.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
	})

	t.Run("GetResultsBeforeCompletion", func(t *testing.T) {
		// Submit execution
		resp, err := deClient.RunAsync(authCtx, &api.RunAsyncRequest{
			Algorithm: "gde3",
			Problem:   "zdt1",
			Variant:   "rand1",
			DeConfig:  config,
		})
		require.NoError(t, err)
		executionID := resp.ExecutionId

		// Try to get results immediately (execution is still pending/running)
		_, err = deClient.GetExecutionResults(authCtx, &api.GetExecutionResultsRequest{
			ExecutionId: executionID,
		})
		assert.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.FailedPrecondition, st.Code())

		// Cleanup
		_, _ = deClient.CancelExecution(authCtx, &api.CancelExecutionRequest{
			ExecutionId: executionID,
		})
	})

	t.Run("DeleteRunningExecution", func(t *testing.T) {
		// Submit execution
		resp, err := deClient.RunAsync(authCtx, &api.RunAsyncRequest{
			Algorithm: "gde3",
			Problem:   "zdt1",
			Variant:   "rand1",
			DeConfig:  config,
		})
		require.NoError(t, err)
		executionID := resp.ExecutionId

		// Wait a bit to ensure execution starts
		time.Sleep(100 * time.Millisecond)

		// Delete execution
		_, err = deClient.DeleteExecution(authCtx, &api.DeleteExecutionRequest{
			ExecutionId: executionID,
		})
		require.NoError(t, err)

		// Verify deletion
		_, err = deClient.GetExecutionStatus(authCtx, &api.GetExecutionStatusRequest{
			ExecutionId: executionID,
		})
		assert.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
	})

	t.Run("DoubleCancel", func(t *testing.T) {
		// Submit execution
		resp, err := deClient.RunAsync(authCtx, &api.RunAsyncRequest{
			Algorithm: "gde3",
			Problem:   "zdt1",
			Variant:   "rand1",
			DeConfig:  config,
		})
		require.NoError(t, err)
		executionID := resp.ExecutionId

		// Cancel first time
		_, err = deClient.CancelExecution(authCtx, &api.CancelExecutionRequest{
			ExecutionId: executionID,
		})
		require.NoError(t, err)

		// Cancel second time - should not error
		_, err = deClient.CancelExecution(authCtx, &api.CancelExecutionRequest{
			ExecutionId: executionID,
		})
		assert.NoError(t, err, "double cancel should be idempotent")
	})
}

// TestE2E_ProgressStreaming tests progress streaming scenarios
func TestE2E_ProgressStreaming(t *testing.T) {
	conn := setupConnection(t)
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	// Setup authenticated user
	username := fmt.Sprintf("e2e_progress_%d", time.Now().Unix())
	token := setupTestUser(t, ctx, conn, username)
	authCtx := metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token)

	deClient := api.NewDifferentialEvolutionServiceClient(conn)

	config := &api.DEConfig{
		Executions:     2,
		Generations:    50,
		PopulationSize: 30,
		DimensionsSize: 10,
		ObjectivesSize: 2,
		FloorLimiter:   0.0,
		CeilLimiter:    1.0,
		AlgorithmConfig: &api.DEConfig_Gde3{
			Gde3: &api.GDE3Config{
				Cr: 0.5,
				F:  0.5,
				P:  0.1,
			},
		},
	}

	t.Run("StreamProgressForNonexistentExecution", func(t *testing.T) {
		stream, err := deClient.StreamProgress(authCtx, &api.StreamProgressRequest{
			ExecutionId: "nonexistent-execution-id",
		})
		require.NoError(t, err)

		_, err = stream.Recv()
		assert.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
	})

	t.Run("StreamProgressWithClientDisconnect", func(t *testing.T) {
		// Submit execution
		resp, err := deClient.RunAsync(authCtx, &api.RunAsyncRequest{
			Algorithm: "gde3",
			Problem:   "zdt1",
			Variant:   "rand1",
			DeConfig:  config,
		})
		require.NoError(t, err)
		executionID := resp.ExecutionId

		// Start streaming
		streamCtx, streamCancel := context.WithCancel(authCtx)
		stream, err := deClient.StreamProgress(streamCtx, &api.StreamProgressRequest{
			ExecutionId: executionID,
		})
		require.NoError(t, err)

		// Receive a few progress updates
		for i := 0; i < 3; i++ {
			_, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Logf("Progress receive error: %v", err)
				break
			}
		}

		// Cancel stream (simulate client disconnect)
		streamCancel()

		// Verify we get context cancelled error
		_, err = stream.Recv()
		assert.Error(t, err)

		// Cleanup
		_, _ = deClient.CancelExecution(authCtx, &api.CancelExecutionRequest{
			ExecutionId: executionID,
		})
	})

	t.Run("StreamProgressForCompletedExecution", func(t *testing.T) {
		// Submit quick execution
		quickConfig := *config
		quickConfig.Executions = 1
		quickConfig.Generations = 10
		quickConfig.PopulationSize = 20

		resp, err := deClient.RunAsync(authCtx, &api.RunAsyncRequest{
			Algorithm: "gde3",
			Problem:   "zdt1",
			Variant:   "rand1",
			DeConfig:  &quickConfig,
		})
		require.NoError(t, err)
		executionID := resp.ExecutionId

		// Wait for completion
		time.Sleep(2 * time.Second)

		// Try to stream progress for completed execution
		stream, err := deClient.StreamProgress(authCtx, &api.StreamProgressRequest{
			ExecutionId: executionID,
		})
		require.NoError(t, err)

		// Should get EOF or no updates quickly
		timeout := time.After(2 * time.Second)
		updates := 0
	receiveLoop:
		for {
			select {
			case <-timeout:
				break receiveLoop
			default:
				_, err := stream.Recv()
				if err == io.EOF {
					break receiveLoop
				}
				if err != nil {
					break receiveLoop
				}
				updates++
			}
		}

		t.Logf("Received %d progress updates for completed execution", updates)
	})
}

// TestE2E_ConcurrentExecutions tests concurrent execution scenarios
func TestE2E_ConcurrentExecutions(t *testing.T) {
	conn := setupConnection(t)
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	// Setup authenticated user
	username := fmt.Sprintf("e2e_concurrent_%d", time.Now().Unix())
	token := setupTestUser(t, ctx, conn, username)
	authCtx := metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token)

	deClient := api.NewDifferentialEvolutionServiceClient(conn)

	config := &api.DEConfig{
		Executions:     1,
		Generations:    100,
		PopulationSize: 50,
		DimensionsSize: 10,
		ObjectivesSize: 2,
		FloorLimiter:   0.0,
		CeilLimiter:    1.0,
		AlgorithmConfig: &api.DEConfig_Gde3{
			Gde3: &api.GDE3Config{
				Cr: 0.5,
				F:  0.5,
				P:  0.1,
			},
		},
	}

	t.Run("MultipleSimultaneousExecutions", func(t *testing.T) {
		// Submit multiple executions concurrently
		numExecutions := 5
		executionIDs := make([]string, numExecutions)

		for i := 0; i < numExecutions; i++ {
			resp, err := deClient.RunAsync(authCtx, &api.RunAsyncRequest{
				Algorithm: "gde3",
				Problem:   "zdt1",
				Variant:   "rand1",
				DeConfig:  config,
			})
			require.NoError(t, err)
			executionIDs[i] = resp.ExecutionId
		}

		// Verify all executions are listed
		listResp, err := deClient.ListExecutions(authCtx, &api.ListExecutionsRequest{
			Limit:  100,
			Offset: 0,
		})
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(listResp.Executions), numExecutions)

		// Cleanup - cancel all executions
		for _, execID := range executionIDs {
			_, _ = deClient.CancelExecution(authCtx, &api.CancelExecutionRequest{
				ExecutionId: execID,
			})
		}
	})

	t.Run("ListExecutionsWithPagination", func(t *testing.T) {
		// List with small page size
		page1, err := deClient.ListExecutions(authCtx, &api.ListExecutionsRequest{
			Limit:  2,
			Offset: 0,
		})
		require.NoError(t, err)
		assert.LessOrEqual(t, len(page1.Executions), 2)

		// If there are more results, get next page
		if page1.HasMore {
			page2, err := deClient.ListExecutions(authCtx, &api.ListExecutionsRequest{
				Limit:  2,
				Offset: 2,
			})
			require.NoError(t, err)
			assert.LessOrEqual(t, len(page2.Executions), 2)
		}
	})

	t.Run("ListExecutionsWithStatusFilter", func(t *testing.T) {
		// List only completed executions
		listResp, err := deClient.ListExecutions(authCtx, &api.ListExecutionsRequest{
			Status: api.ExecutionStatus_EXECUTION_STATUS_COMPLETED,
			Limit:  100,
			Offset: 0,
		})
		require.NoError(t, err)

		// Verify all returned executions have COMPLETED status
		for _, exec := range listResp.Executions {
			assert.Equal(t, api.ExecutionStatus_EXECUTION_STATUS_COMPLETED, exec.Status)
		}
	})
}

// TestE2E_AccessControl tests cross-user access control
func TestE2E_AccessControl(t *testing.T) {
	conn := setupConnection(t)
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	// Setup two users
	user1 := fmt.Sprintf("e2e_user1_%d", time.Now().Unix())
	token1 := setupTestUser(t, ctx, conn, user1)
	authCtx1 := metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token1)

	user2 := fmt.Sprintf("e2e_user2_%d", time.Now().Unix())
	token2 := setupTestUser(t, ctx, conn, user2)
	authCtx2 := metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token2)

	deClient := api.NewDifferentialEvolutionServiceClient(conn)

	config := &api.DEConfig{
		Executions:     1,
		Generations:    50,
		PopulationSize: 30,
		DimensionsSize: 10,
		ObjectivesSize: 2,
		FloorLimiter:   0.0,
		CeilLimiter:    1.0,
		AlgorithmConfig: &api.DEConfig_Gde3{
			Gde3: &api.GDE3Config{
				Cr: 0.5,
				F:  0.5,
				P:  0.1,
			},
		},
	}

	t.Run("AccessAnotherUsersExecution", func(t *testing.T) {
		// User1 submits execution
		resp, err := deClient.RunAsync(authCtx1, &api.RunAsyncRequest{
			Algorithm: "gde3",
			Problem:   "zdt1",
			Variant:   "rand1",
			DeConfig:  config,
		})
		require.NoError(t, err)
		executionID := resp.ExecutionId

		// User2 tries to access User1's execution
		_, err = deClient.GetExecutionStatus(authCtx2, &api.GetExecutionStatusRequest{
			ExecutionId: executionID,
		})
		assert.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())

		// User2 tries to cancel User1's execution
		_, err = deClient.CancelExecution(authCtx2, &api.CancelExecutionRequest{
			ExecutionId: executionID,
		})
		assert.Error(t, err)
		st, ok = status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())

		// Cleanup
		_, _ = deClient.CancelExecution(authCtx1, &api.CancelExecutionRequest{
			ExecutionId: executionID,
		})
	})

	t.Run("ListExecutionsOnlyOwnExecutions", func(t *testing.T) {
		// User1 submits execution
		_, err := deClient.RunAsync(authCtx1, &api.RunAsyncRequest{
			Algorithm: "gde3",
			Problem:   "zdt1",
			Variant:   "rand1",
			DeConfig:  config,
		})
		require.NoError(t, err)

		// User2 lists executions - should not see User1's execution
		listResp, err := deClient.ListExecutions(authCtx2, &api.ListExecutionsRequest{
			Limit:  100,
			Offset: 0,
		})
		require.NoError(t, err)

		// Verify no execution from User1
		for _, exec := range listResp.Executions {
			assert.NotEqual(t, user1, exec.UserId)
		}
	})
}

// setupTestUser is a helper to create and login a test user
func setupTestUser(t *testing.T, ctx context.Context, conn *grpc.ClientConn, username string) string {
	t.Helper()

	authClient := api.NewAuthServiceClient(conn)

	// Register
	_, err := authClient.Register(ctx, &api.AuthServiceRegisterRequest{
		User: &api.User{
			Ids:      &api.UserIDs{Username: username},
			Email:    fmt.Sprintf("%s@test.com", username),
			Password: "test_password_123",
		},
	})
	require.NoError(t, err)

	// Login
	loginResp, err := authClient.Login(ctx, &api.AuthServiceLoginRequest{
		Username: username,
		Password: "test_password_123",
	})
	require.NoError(t, err)

	return loginResp.Token
}
