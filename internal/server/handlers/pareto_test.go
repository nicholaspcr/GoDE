package handlers

import (
	"context"
	"testing"

	"github.com/nicholaspcr/GoDE/internal/store/mock"
	storerrors "github.com/nicholaspcr/GoDE/internal/store/errors"
	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestParetoHandler_Get(t *testing.T) {
	mockStore := &mock.MockStore{}
	handler := NewParetoHandler(mockStore)

	ctx := context.Background()

	t.Run("successful get", func(t *testing.T) {
		testPareto := &api.Pareto{
			Ids:     &api.ParetoIDs{Id: 1, UserId: "testuser"},
			MaxObjs: []float64{1.0, 2.0},
			Vectors: []*api.Vector{
				{
					Elements:         []float64{0.5, 0.5},
					Objectives:       []float64{1.0, 1.0},
					CrowdingDistance: 0.5,
				},
			},
		}

		mockStore.GetParetoFn = func(ctx context.Context, ids *api.ParetoIDs) (*api.Pareto, error) {
			return testPareto, nil
		}

		req := &api.ParetoServiceGetRequest{
			ParetoIds: &api.ParetoIDs{Id: 1},
		}

		resp, err := handler.(*paretoHandler).Get(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Pareto)
		assert.Equal(t, uint64(1), resp.Pareto.Ids.Id)
		assert.Len(t, resp.Pareto.Vectors, 1)
	})

	t.Run("missing pareto_ids", func(t *testing.T) {
		req := &api.ParetoServiceGetRequest{}

		resp, err := handler.(*paretoHandler).Get(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, resp)

		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
	})

	t.Run("pareto not found", func(t *testing.T) {
		mockStore.GetParetoFn = func(ctx context.Context, ids *api.ParetoIDs) (*api.Pareto, error) {
			return nil, storerrors.ErrParetoSetNotFound
		}

		req := &api.ParetoServiceGetRequest{
			ParetoIds: &api.ParetoIDs{Id: 9999},
		}

		resp, err := handler.(*paretoHandler).Get(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, resp)

		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
	})
}

func TestParetoHandler_Delete(t *testing.T) {
	mockStore := &mock.MockStore{}
	handler := NewParetoHandler(mockStore)

	ctx := context.Background()

	t.Run("successful delete", func(t *testing.T) {
		mockStore.DeleteParetoFn = func(ctx context.Context, ids *api.ParetoIDs) error {
			return nil
		}

		req := &api.ParetoServiceDeleteRequest{
			ParetoIds: &api.ParetoIDs{Id: 1},
		}

		resp, err := handler.(*paretoHandler).Delete(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("missing pareto_ids", func(t *testing.T) {
		req := &api.ParetoServiceDeleteRequest{}

		resp, err := handler.(*paretoHandler).Delete(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, resp)

		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
	})
}

func TestParetoHandler_ListByUser(t *testing.T) {
	mockStore := &mock.MockStore{}
	handler := NewParetoHandler(mockStore)

	ctx := context.Background()

	t.Run("missing user_ids", func(t *testing.T) {
		req := &api.ParetoServiceListByUserRequest{}

		// Create a mock stream
		mockStream := &mockParetoStream{ctx: ctx}

		err := handler.(*paretoHandler).ListByUser(req, mockStream)
		assert.Error(t, err)

		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
	})

	t.Run("successful list", func(t *testing.T) {
		testParetos := []*api.Pareto{
			{
				Ids:     &api.ParetoIDs{Id: 1, UserId: "testuser"},
				MaxObjs: []float64{1.0, 2.0},
				Vectors: []*api.Vector{
					{
						Elements:         []float64{0.5, 0.5},
						Objectives:       []float64{1.0, 1.0},
						CrowdingDistance: 0.5,
					},
				},
			},
			{
				Ids:     &api.ParetoIDs{Id: 2, UserId: "testuser"},
				MaxObjs: []float64{2.0, 3.0},
				Vectors: []*api.Vector{
					{
						Elements:         []float64{0.6, 0.6},
						Objectives:       []float64{2.0, 2.0},
						CrowdingDistance: 0.6,
					},
				},
			},
		}

		mockStore.ListParetosFn = func(ctx context.Context, userIds *api.UserIDs, limit, offset int) ([]*api.Pareto, int, error) {
			return testParetos, len(testParetos), nil
		}

		req := &api.ParetoServiceListByUserRequest{
			UserIds: &api.UserIDs{Username: "testuser"},
		}

		mockStream := &mockParetoStream{ctx: ctx}

		err := handler.(*paretoHandler).ListByUser(req, mockStream)
		assert.NoError(t, err)
		assert.Len(t, mockStream.sentMsgs, 2)
	})
}

// mockParetoStream is a mock implementation of api.ParetoService_ListByUserServer for testing
type mockParetoStream struct {
	ctx      context.Context
	sentMsgs []*api.ParetoServiceListByUserResponse
}

func (m *mockParetoStream) Send(resp *api.ParetoServiceListByUserResponse) error {
	m.sentMsgs = append(m.sentMsgs, resp)
	return nil
}

func (m *mockParetoStream) Context() context.Context {
	return m.ctx
}

func (m *mockParetoStream) SendMsg(msg any) error {
	return nil
}

func (m *mockParetoStream) RecvMsg(msg any) error {
	return nil
}

func (m *mockParetoStream) SetHeader(md metadata.MD) error {
	return nil
}

func (m *mockParetoStream) SendHeader(md metadata.MD) error {
	return nil
}

func (m *mockParetoStream) SetTrailer(md metadata.MD) {
}
