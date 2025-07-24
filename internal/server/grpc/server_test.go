package grpc

import (
	__ "GophKeeper/internal/pkg/proto_gen"
	"GophKeeper/internal/server/jwt"
	"GophKeeper/internal/server/manager"
	"GophKeeper/internal/server/model"
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

type mockServerUserManager struct {
	registerFunc func(login, password, masterPassword string) (string, error)
	loginFunc    func(login, password, masterPassword string) (string, error)
}

func (m *mockServerUserManager) Register(login, password, masterPassword string) (string, error) {
	return m.registerFunc(login, password, masterPassword)
}

func (m *mockServerUserManager) Login(login, password, masterPassword string) (string, error) {
	return m.loginFunc(login, password, masterPassword)
}

func (m *mockServerUserManager) DecodeToken(token string) (*jwt.Claims, error) {
	return nil, errors.New("not implemented")
}

type mockUserDataManager struct {
	mergeFunc      func(ctx context.Context, data *model.UserData) error
	getUpdatesFunc func(ctx context.Context, userID uint32, updatedAfter time.Time) ([]*model.UserData, error)
}

func (m *mockUserDataManager) Merge(ctx context.Context, data *model.UserData) error {
	return m.mergeFunc(ctx, data)
}

func (m *mockUserDataManager) GetUpdates(ctx context.Context, userID uint32, updatedAfter time.Time) ([]*model.UserData, error) {
	return m.getUpdatesFunc(ctx, userID, updatedAfter)
}

func TestNewServer(t *testing.T) {
	um := &manager.UserManager{}
	dm := &manager.UserDataManager{}
	s := NewServer(um, dm)

	if s.userManager != um {
		t.Errorf("NewServer() userManager = %v, want %v", s.userManager, um)
	}

	if s.dataManager != dm {
		t.Errorf("NewServer() dataManager = %v, want %v", s.dataManager, dm)
	}
}

func TestServer_Register(t *testing.T) {
	tests := []struct {
		name          string
		req           *__.RegisterRequest
		setupMock     func() UserManagerInterface
		want          *__.TokenResponse
		wantErrCode   codes.Code
		wantErrString string
	}{
		{
			name: "successful registration",
			req: &__.RegisterRequest{
				Login:          "user1",
				Password:       "pass1",
				MasterPassword: "master1",
			},
			setupMock: func() UserManagerInterface {
				return &mockServerUserManager{
					registerFunc: func(login, password, masterPassword string) (string, error) {
						return "token123", nil
					},
				}
			},
			want: &__.TokenResponse{Token: "token123"},
		},
		{
			name: "user already exists",
			req: &__.RegisterRequest{
				Login:          "user1",
				Password:       "pass1",
				MasterPassword: "master1",
			},
			setupMock: func() UserManagerInterface {
				return &mockServerUserManager{
					registerFunc: func(login, password, masterPassword string) (string, error) {
						return "", manager.ErrUserExists
					},
				}
			},
			wantErrCode: codes.AlreadyExists,
		},
		{
			name: "internal error",
			req: &__.RegisterRequest{
				Login:          "user1",
				Password:       "pass1",
				MasterPassword: "master1",
			},
			setupMock: func() UserManagerInterface {
				return &mockServerUserManager{
					registerFunc: func(login, password, masterPassword string) (string, error) {
						return "", errors.New("some error")
					},
				}
			},
			wantErrCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUM := tt.setupMock()
			s := &Server{
				userManager: mockUM,
			}

			got, err := s.Register(context.Background(), tt.req)

			if tt.wantErrCode != 0 {
				if err == nil {
					t.Fatalf("Register() error = nil, want error with code %v", tt.wantErrCode)
				}

				statusErr, ok := status.FromError(err)
				if !ok {
					t.Fatalf("Register() error is not a status error")
				}

				if statusErr.Code() != tt.wantErrCode {
					t.Errorf("Register() error code = %v, want %v", statusErr.Code(), tt.wantErrCode)
				}
			} else {
				if err != nil {
					t.Fatalf("Register() error = %v, want nil", err)
				}

				if got.Token != tt.want.Token {
					t.Errorf("Register() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestServer_Login(t *testing.T) {
	tests := []struct {
		name          string
		req           *__.LoginRequest
		setupMock     func() UserManagerInterface
		want          *__.TokenResponse
		wantErrCode   codes.Code
		wantErrString string
	}{
		{
			name: "successful login",
			req: &__.LoginRequest{
				Login:          "user1",
				Password:       "pass1",
				MasterPassword: "master1",
			},
			setupMock: func() UserManagerInterface {
				return &mockServerUserManager{
					loginFunc: func(login, password, masterPassword string) (string, error) {
						return "token123", nil
					},
				}
			},
			want: &__.TokenResponse{Token: "token123"},
		},
		{
			name: "invalid credentials",
			req: &__.LoginRequest{
				Login:          "user1",
				Password:       "wrong",
				MasterPassword: "master1",
			},
			setupMock: func() UserManagerInterface {
				return &mockServerUserManager{
					loginFunc: func(login, password, masterPassword string) (string, error) {
						return "", manager.ErrInvalidCredentials
					},
				}
			},
			wantErrCode: codes.Unauthenticated,
		},
		{
			name: "internal error",
			req: &__.LoginRequest{
				Login:          "user1",
				Password:       "pass1",
				MasterPassword: "master1",
			},
			setupMock: func() UserManagerInterface {
				return &mockServerUserManager{
					loginFunc: func(login, password, masterPassword string) (string, error) {
						return "", errors.New("some error")
					},
				}
			},
			wantErrCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUM := tt.setupMock()
			s := &Server{
				userManager: mockUM,
			}

			got, err := s.Login(context.Background(), tt.req)

			if tt.wantErrCode != 0 {
				if err == nil {
					t.Fatalf("Login() error = nil, want error with code %v", tt.wantErrCode)
				}

				statusErr, ok := status.FromError(err)
				if !ok {
					t.Fatalf("Login() error is not a status error")
				}

				if statusErr.Code() != tt.wantErrCode {
					t.Errorf("Login() error code = %v, want %v", statusErr.Code(), tt.wantErrCode)
				}
			} else {
				if err != nil {
					t.Fatalf("Login() error = %v, want nil", err)
				}

				if got.Token != tt.want.Token {
					t.Errorf("Login() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestServer_Merge(t *testing.T) {
	testTime := time.Now().UTC()
	testTimePb := timestamppb.New(testTime)

	tests := []struct {
		name          string
		req           *__.MergeRequest
		ctx           context.Context
		setupMock     func() UserDataManagerInterface
		want          *__.DataResponse
		wantErrCode   codes.Code
		wantErrString string
	}{
		{
			name: "successful merge",
			req: &__.MergeRequest{
				DataKey:   "key1",
				DataValue: []byte("value1"),
				UpdatedAt: testTimePb,
				DeletedAt: testTimePb,
			},
			ctx: context.WithValue(
				context.Background(),
				userClaimsKey{},
				&jwt.Claims{SubjectID: uint32(123)},
			),
			setupMock: func() UserDataManagerInterface {
				return &mockUserDataManager{
					mergeFunc: func(ctx context.Context, data *model.UserData) error {
						return nil
					},
				}
			},
			want: &__.DataResponse{
				DataKey:   "key1",
				DataValue: []byte("value1"),
				UpdatedAt: testTimePb,
				DeletedAt: testTimePb,
			},
		},
		{
			name: "no auth in context",
			req: &__.MergeRequest{
				DataKey:   "key1",
				DataValue: []byte("value1"),
				UpdatedAt: testTimePb,
				DeletedAt: testTimePb,
			},
			ctx: context.Background(),
			setupMock: func() UserDataManagerInterface {
				return &mockUserDataManager{}
			},
			wantErrCode: codes.Unauthenticated,
		},
		{
			name: "internal error",
			req: &__.MergeRequest{
				DataKey:   "key1",
				DataValue: []byte("value1"),
				UpdatedAt: testTimePb,
				DeletedAt: testTimePb,
			},
			ctx: context.WithValue(
				context.Background(),
				userClaimsKey{},
				&jwt.Claims{SubjectID: uint32(123)},
			),
			setupMock: func() UserDataManagerInterface {
				return &mockUserDataManager{
					mergeFunc: func(ctx context.Context, data *model.UserData) error {
						return errors.New("some error")
					},
				}
			},
			wantErrCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDM := tt.setupMock()
			s := &Server{
				dataManager: mockDM,
			}

			got, err := s.Merge(tt.ctx, tt.req)

			if tt.wantErrCode != 0 {
				if err == nil {
					t.Fatalf("Merge() error = nil, want error with code %v", tt.wantErrCode)
				}

				statusErr, ok := status.FromError(err)
				if !ok {
					t.Fatalf("Merge() error is not a status error")
				}

				if statusErr.Code() != tt.wantErrCode {
					t.Errorf("Merge() error code = %v, want %v", statusErr.Code(), tt.wantErrCode)
				}
			} else {
				if err != nil {
					t.Fatalf("Merge() error = %v, want nil", err)
				}

				if got.DataKey != tt.want.DataKey {
					t.Errorf("Merge() DataKey = %v, want %v", got.DataKey, tt.want.DataKey)
				}
			}
		})
	}
}

func TestServer_GetUpdates(t *testing.T) {
	testTime := time.Now().UTC()
	testTimePb := timestamppb.New(testTime)

	tests := []struct {
		name          string
		req           *__.GetUpdatesRequest
		ctx           context.Context
		setupMock     func() UserDataManagerInterface
		want          *__.DataListResponse
		wantErrCode   codes.Code
		wantErrString string
	}{
		{
			name: "successful get updates",
			req: &__.GetUpdatesRequest{
				UpdatedAfter: testTimePb,
			},
			ctx: context.WithValue(
				context.Background(),
				userClaimsKey{},
				&jwt.Claims{SubjectID: uint32(123)},
			),
			setupMock: func() UserDataManagerInterface {
				return &mockUserDataManager{
					getUpdatesFunc: func(ctx context.Context, userID uint32, updatedAfter time.Time) ([]*model.UserData, error) {
						return []*model.UserData{
							{
								UserID:    uint32(123),
								DataKey:   "key1",
								DataValue: []byte("value1"),
								UpdatedAt: testTime,
								DeletedAt: testTime,
							},
						}, nil
					},
				}
			},
			want: &__.DataListResponse{
				Items: []*__.DataResponse{
					{
						DataKey:   "key1",
						DataValue: []byte("value1"),
						UpdatedAt: testTimePb,
						DeletedAt: testTimePb,
					},
				},
			},
		},
		{
			name: "no auth in context",
			req: &__.GetUpdatesRequest{
				UpdatedAfter: testTimePb,
			},
			ctx: context.Background(),
			setupMock: func() UserDataManagerInterface {
				return &mockUserDataManager{}
			},
			wantErrCode: codes.Unauthenticated,
		},
		{
			name: "internal error",
			req: &__.GetUpdatesRequest{
				UpdatedAfter: testTimePb,
			},
			ctx: context.WithValue(
				context.Background(),
				userClaimsKey{},
				&jwt.Claims{SubjectID: uint32(123)},
			),
			setupMock: func() UserDataManagerInterface {
				return &mockUserDataManager{
					getUpdatesFunc: func(ctx context.Context, userID uint32, updatedAfter time.Time) ([]*model.UserData, error) {
						return nil, errors.New("some error")
					},
				}
			},
			wantErrCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDM := tt.setupMock()
			s := &Server{
				dataManager: mockDM,
			}

			got, err := s.GetUpdates(tt.ctx, tt.req)

			if tt.wantErrCode != 0 {
				if err == nil {
					t.Fatalf("GetUpdates() error = nil, want error with code %v", tt.wantErrCode)
				}

				statusErr, ok := status.FromError(err)
				if !ok {
					t.Fatalf("GetUpdates() error is not a status error")
				}

				if statusErr.Code() != tt.wantErrCode {
					t.Errorf("GetUpdates() error code = %v, want %v", statusErr.Code(), tt.wantErrCode)
				}
			} else {
				if err != nil {
					t.Fatalf("GetUpdates() error = %v, want nil", err)
				}

				if len(got.Items) != len(tt.want.Items) {
					t.Errorf("GetUpdates() Items len = %v, want %v", len(got.Items), len(tt.want.Items))
				}

				if len(got.Items) > 0 && got.Items[0].DataKey != tt.want.Items[0].DataKey {
					t.Errorf("GetUpdates() DataKey = %v, want %v", got.Items[0].DataKey, tt.want.Items[0].DataKey)
				}
			}
		})
	}
}

func TestConvertError(t *testing.T) {
	tests := []struct {
		name        string
		err         error
		wantCode    codes.Code
		wantMessage string
	}{
		{
			name:        "user exists error",
			err:         manager.ErrUserExists,
			wantCode:    codes.AlreadyExists,
			wantMessage: manager.ErrUserExists.Error(),
		},
		{
			name:        "invalid credentials error",
			err:         manager.ErrInvalidCredentials,
			wantCode:    codes.Unauthenticated,
			wantMessage: manager.ErrInvalidCredentials.Error(),
		},
		{
			name:        "other error",
			err:         errors.New("some error"),
			wantCode:    codes.Internal,
			wantMessage: "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertError(tt.err)

			statusErr, ok := status.FromError(got)
			if !ok {
				t.Fatalf("convertError() result is not a status error")
			}

			if statusErr.Code() != tt.wantCode {
				t.Errorf("convertError() code = %v, want %v", statusErr.Code(), tt.wantCode)
			}

			if statusErr.Message() != tt.wantMessage {
				t.Errorf("convertError() message = %v, want %v", statusErr.Message(), tt.wantMessage)
			}
		})
	}
}
