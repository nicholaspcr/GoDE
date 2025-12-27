package handlers

import (
	"context"
	"errors"
	"testing"

	"github.com/nicholaspcr/GoDE/internal/store/mock"
	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func TestUserHandler_Create(t *testing.T) {
	tests := []struct {
		req       *api.UserServiceCreateRequest
		setupMock func(*mock.MockStore)
		name      string
		wantCode  codes.Code
		wantErr   bool
	}{
		{
			name: "successful creation",
			req: &api.UserServiceCreateRequest{
				User: &api.User{
					Ids:      &api.UserIDs{Username: "testuser"},
					Email:    "test@example.com",
					Password: "validpass123",
				},
			},
			setupMock: func(m *mock.MockStore) {
				m.CreateUserFn = func(ctx context.Context, user *api.User) error {
					return nil
				}
			},
			wantErr: false,
		},
		{
			name: "invalid username",
			req: &api.UserServiceCreateRequest{
				User: &api.User{
					Ids:      &api.UserIDs{Username: "ab"},
					Email:    "test@example.com",
					Password: "validpass123",
				},
			},
			setupMock: func(m *mock.MockStore) {
				// Should not be called
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "invalid email",
			req: &api.UserServiceCreateRequest{
				User: &api.User{
					Ids:      &api.UserIDs{Username: "testuser"},
					Email:    "invalid-email",
					Password: "validpass123",
				},
			},
			setupMock: func(m *mock.MockStore) {
				// Should not be called
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "password too short",
			req: &api.UserServiceCreateRequest{
				User: &api.User{
					Ids:      &api.UserIDs{Username: "testuser"},
					Email:    "test@example.com",
					Password: "short",
				},
			},
			setupMock: func(m *mock.MockStore) {
				// Should not be called
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "database error",
			req: &api.UserServiceCreateRequest{
				User: &api.User{
					Ids:      &api.UserIDs{Username: "testuser"},
					Email:    "test@example.com",
					Password: "validpass123",
				},
			},
			setupMock: func(m *mock.MockStore) {
				m.CreateUserFn = func(ctx context.Context, user *api.User) error {
					return errors.New("database error")
				}
			},
			wantErr:  true,
			wantCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &mock.MockStore{}
			tt.setupMock(mockStore)

			handler := NewUserHandler(mockStore)

			_, err := handler.(*userHandler).Create(context.Background(), tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.wantCode != codes.OK {
					st, ok := status.FromError(err)
					require.True(t, ok)
					assert.Equal(t, tt.wantCode, st.Code())
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserHandler_Get(t *testing.T) {
	tests := []struct {
		req        *api.UserServiceGetRequest
		setupMock  func(*mock.MockStore)
		wantResult *api.User
		name       string
		wantErr    bool
	}{
		{
			name: "successful get",
			req: &api.UserServiceGetRequest{
				UserIds: &api.UserIDs{Username: "testuser"},
			},
			setupMock: func(m *mock.MockStore) {
				m.GetUserFn = func(ctx context.Context, userIDs *api.UserIDs) (*api.User, error) {
					return &api.User{
						Ids:      &api.UserIDs{Username: "testuser"},
						Email:    "test@example.com",
						Password: "hashed",
					}, nil
				}
			},
			wantErr: false,
			wantResult: &api.User{
				Ids:      &api.UserIDs{Username: "testuser"},
				Email:    "test@example.com",
				Password: "hashed",
			},
		},
		{
			name: "user not found",
			req: &api.UserServiceGetRequest{
				UserIds: &api.UserIDs{Username: "nonexistent"},
			},
			setupMock: func(m *mock.MockStore) {
				m.GetUserFn = func(ctx context.Context, userIDs *api.UserIDs) (*api.User, error) {
					return nil, errors.New("user not found")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &mock.MockStore{}
			tt.setupMock(mockStore)

			handler := NewUserHandler(mockStore)

			resp, err := handler.(*userHandler).Get(context.Background(), tt.req)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				assert.Equal(t, tt.wantResult.Ids.Username, resp.User.Ids.Username)
				assert.Equal(t, tt.wantResult.Email, resp.User.Email)
			}
		})
	}
}

func TestUserHandler_Update(t *testing.T) {
	tests := []struct {
		req       *api.UserServiceUpdateRequest
		setupMock func(*mock.MockStore)
		name      string
		wantErr   bool
	}{
		{
			name: "successful update",
			req: &api.UserServiceUpdateRequest{
				User: &api.User{
					Ids:   &api.UserIDs{Username: "testuser"},
					Email: "newemail@example.com",
				},
				FieldMask: &fieldmaskpb.FieldMask{
					Paths: []string{"email"},
				},
			},
			setupMock: func(m *mock.MockStore) {
				m.UpdateUserFn = func(ctx context.Context, user *api.User, fields ...string) error {
					assert.Equal(t, []string{"email"}, fields)
					return nil
				}
			},
			wantErr: false,
		},
		{
			name: "update error",
			req: &api.UserServiceUpdateRequest{
				User: &api.User{
					Ids:   &api.UserIDs{Username: "testuser"},
					Email: "newemail@example.com",
				},
				FieldMask: &fieldmaskpb.FieldMask{
					Paths: []string{"email"},
				},
			},
			setupMock: func(m *mock.MockStore) {
				m.UpdateUserFn = func(ctx context.Context, user *api.User, fields ...string) error {
					return errors.New("update failed")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &mock.MockStore{}
			tt.setupMock(mockStore)

			handler := NewUserHandler(mockStore)

			_, err := handler.(*userHandler).Update(context.Background(), tt.req)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserHandler_Delete(t *testing.T) {
	tests := []struct {
		req       *api.UserServiceDeleteRequest
		setupMock func(*mock.MockStore)
		name      string
		wantErr   bool
	}{
		{
			name: "successful deletion",
			req: &api.UserServiceDeleteRequest{
				UserIds: &api.UserIDs{Username: "testuser"},
			},
			setupMock: func(m *mock.MockStore) {
				m.DeleteUserFn = func(ctx context.Context, userIDs *api.UserIDs) error {
					return nil
				}
			},
			wantErr: false,
		},
		{
			name: "deletion error",
			req: &api.UserServiceDeleteRequest{
				UserIds: &api.UserIDs{Username: "testuser"},
			},
			setupMock: func(m *mock.MockStore) {
				m.DeleteUserFn = func(ctx context.Context, userIDs *api.UserIDs) error {
					return errors.New("delete failed")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &mock.MockStore{}
			tt.setupMock(mockStore)

			handler := NewUserHandler(mockStore)

			_, err := handler.(*userHandler).Delete(context.Background(), tt.req)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
