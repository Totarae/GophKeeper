package syncro

import (
	"context"
	"errors"
	"github.com/Totarae/GophKeeper/internal/client/model"
	__ "github.com/Totarae/GophKeeper/internal/pkg/proto_gen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type MockGRPCClient struct {
	mock.Mock
}

func (m *MockGRPCClient) Merge(ctx context.Context, data *model.UserData) (*__.DataResponse, error) {
	args := m.Called(ctx, data)
	return args.Get(0).(*__.DataResponse), args.Error(1)
}

func (m *MockGRPCClient) GetUpdates(ctx context.Context, lastSync time.Time) (*__.DataListResponse, error) {
	args := m.Called(ctx, lastSync)
	return args.Get(0).(*__.DataListResponse), args.Error(1)
}

type MockUserDataManager struct {
	mock.Mock
}

func (m *MockUserDataManager) GetUpdates(ctx context.Context, lastSync time.Time) ([]*model.UserData, error) {
	args := m.Called(ctx, lastSync)
	return args.Get(0).([]*model.UserData), args.Error(1)
}

func (m *MockUserDataManager) Merge(ctx context.Context, data *model.UserData) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

type MockMetaManager struct {
	mock.Mock
}

func (m *MockMetaManager) GetLastSync(ctx context.Context) (time.Time, error) {
	args := m.Called(ctx)
	return args.Get(0).(time.Time), args.Error(1)
}

func (m *MockMetaManager) SetLastSync(ctx context.Context, lastSync time.Time) error {
	args := m.Called(ctx, lastSync)
	return args.Error(0)
}

func TestNew(t *testing.T) {
	client := &MockGRPCClient{}
	userDataMgr := &MockUserDataManager{}
	metaManager := &MockMetaManager{}
	interval := time.Second * 30

	s := New(client, userDataMgr, metaManager, interval)

	assert.Equal(t, client, s.client)
	assert.Equal(t, userDataMgr, s.userDataMgr)
	assert.Equal(t, metaManager, s.metaManager)
	assert.Equal(t, interval, s.interval)
	assert.NotNil(t, s.stopCh)
}

func TestSynchronizer_StartStop(t *testing.T) {
	client := &MockGRPCClient{}
	userDataMgr := &MockUserDataManager{}
	metaManager := &MockMetaManager{}
	interval := time.Millisecond * 50

	lastSyncTime := time.Now().Add(-time.Hour).UTC()
	metaManager.On("GetLastSync", mock.Anything).Return(lastSyncTime, nil)
	userDataMgr.On("GetUpdates", mock.Anything, lastSyncTime).Return([]*model.UserData{}, nil)
	client.On("GetUpdates", mock.Anything, lastSyncTime).Return(&__.DataListResponse{Items: []*__.DataResponse{}}, nil)
	metaManager.On("SetLastSync", mock.Anything, mock.AnythingOfType("time.Time")).Return(nil)

	s := New(client, userDataMgr, metaManager, interval)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()

	s.Start(ctx)
	time.Sleep(time.Millisecond * 75)
	s.Stop()

	metaManager.AssertCalled(t, "GetLastSync", mock.Anything)
	userDataMgr.AssertCalled(t, "GetUpdates", mock.Anything, lastSyncTime)
	client.AssertCalled(t, "GetUpdates", mock.Anything, lastSyncTime)
	metaManager.AssertCalled(t, "SetLastSync", mock.Anything, mock.AnythingOfType("time.Time"))
}

func TestSynchronizer_syncOnce(t *testing.T) {
	client := &MockGRPCClient{}
	userDataMgr := &MockUserDataManager{}
	metaManager := &MockMetaManager{}
	interval := time.Second * 30

	lastSyncTime := time.Now().Add(-time.Hour).UTC()
	localUpdate := &model.UserData{
		DataKey:   "test-key",
		DataValue: []byte("test-value"),
		UpdatedAt: time.Now(),
	}

	remoteItem := &__.DataResponse{
		DataKey:   "remote-key",
		DataValue: []byte("remote-value"),
		UpdatedAt: nil,
		DeletedAt: nil,
	}

	metaManager.On("GetLastSync", mock.Anything).Return(lastSyncTime, nil)
	userDataMgr.On("GetUpdates", mock.Anything, lastSyncTime).Return([]*model.UserData{localUpdate}, nil)
	client.On("Merge", mock.Anything, localUpdate).Return(&__.DataResponse{}, nil).Once()
	client.On("GetUpdates", mock.Anything, lastSyncTime).Return(&__.DataListResponse{
		Items: []*__.DataResponse{remoteItem},
	}, nil)
	userDataMgr.On("Merge", mock.Anything, mock.AnythingOfType("*model.UserData")).Return(nil)
	metaManager.On("SetLastSync", mock.Anything, mock.AnythingOfType("time.Time")).Return(nil)

	s := New(client, userDataMgr, metaManager, interval)

	ctx := context.Background()
	s.syncOnce(ctx)

	metaManager.AssertCalled(t, "GetLastSync", ctx)
	userDataMgr.AssertCalled(t, "GetUpdates", ctx, lastSyncTime)
	client.AssertCalled(t, "Merge", ctx, localUpdate)
	client.AssertCalled(t, "GetUpdates", ctx, lastSyncTime)
	userDataMgr.AssertCalled(t, "Merge", ctx, mock.AnythingOfType("*model.UserData"))
	metaManager.AssertCalled(t, "SetLastSync", ctx, mock.AnythingOfType("time.Time"))
}

func TestSynchronizer_pushLocalUpdates(t *testing.T) {
	tests := []struct {
		name           string
		setupMocks     func(*MockUserDataManager, *MockGRPCClient)
		expectedResult bool
	}{
		{
			name: "success_no_updates",
			setupMocks: func(userDataMgr *MockUserDataManager, client *MockGRPCClient) {
				userDataMgr.On("GetUpdates", mock.Anything, mock.AnythingOfType("time.Time")).
					Return([]*model.UserData{}, nil)
			},
			expectedResult: true,
		},
		{
			name: "success_with_updates",
			setupMocks: func(userDataMgr *MockUserDataManager, client *MockGRPCClient) {
				updates := []*model.UserData{
					{DataKey: "key1", DataValue: []byte("value1")},
					{DataKey: "key2", DataValue: []byte("value2")},
				}
				userDataMgr.On("GetUpdates", mock.Anything, mock.AnythingOfType("time.Time")).
					Return(updates, nil)
				for _, update := range updates {
					client.On("Merge", mock.Anything, update).Return(&__.DataResponse{}, nil)
				}
			},
			expectedResult: true,
		},
		{
			name: "failure_upsert_error",
			setupMocks: func(userDataMgr *MockUserDataManager, client *MockGRPCClient) {
				updates := []*model.UserData{
					{DataKey: "key1", DataValue: []byte("value1")},
				}
				userDataMgr.On("GetUpdates", mock.Anything, mock.AnythingOfType("time.Time")).
					Return(updates, nil)
				client.On("Merge", mock.Anything, updates[0]).Return(&__.DataResponse{}, errors.New("upsert error"))
			},
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &MockGRPCClient{}
			userDataMgr := &MockUserDataManager{}
			metaManager := &MockMetaManager{}
			interval := time.Second * 30

			tt.setupMocks(userDataMgr, client)

			s := New(client, userDataMgr, metaManager, interval)
			result := s.pushLocalUpdates(context.Background(), time.Now())

			assert.Equal(t, tt.expectedResult, result)
			userDataMgr.AssertExpectations(t)
			client.AssertExpectations(t)
		})
	}
}

func TestSynchronizer_fetchRemoteUpdates(t *testing.T) {
	tests := []struct {
		name           string
		setupMocks     func(*MockGRPCClient, *MockUserDataManager)
		expectedResult bool
	}{
		{
			name: "success_no_updates",
			setupMocks: func(client *MockGRPCClient, userDataMgr *MockUserDataManager) {
				client.On("GetUpdates", mock.Anything, mock.AnythingOfType("time.Time")).
					Return(&__.DataListResponse{Items: []*__.DataResponse{}}, nil)
			},
			expectedResult: true,
		},
		{
			name: "success_with_updates",
			setupMocks: func(client *MockGRPCClient, userDataMgr *MockUserDataManager) {
				updates := []*__.DataResponse{
					{DataKey: "key1", DataValue: []byte("value1")},
					{DataKey: "key2", DataValue: []byte("value2")},
				}
				client.On("GetUpdates", mock.Anything, mock.AnythingOfType("time.Time")).
					Return(&__.DataListResponse{Items: updates}, nil)
				for range updates {
					userDataMgr.On("Merge", mock.Anything, mock.AnythingOfType("*model.UserData")).Return(nil).Once()
				}
			},
			expectedResult: true,
		},
		{
			name: "failure_get_updates_error",
			setupMocks: func(client *MockGRPCClient, userDataMgr *MockUserDataManager) {
				client.On("GetUpdates", mock.Anything, mock.AnythingOfType("time.Time")).
					Return((*__.DataListResponse)(nil), errors.New("get updates error"))
			},
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &MockGRPCClient{}
			userDataMgr := &MockUserDataManager{}
			metaManager := &MockMetaManager{}
			interval := time.Second * 30

			tt.setupMocks(client, userDataMgr)

			s := New(client, userDataMgr, metaManager, interval)
			result := s.fetchRemoteUpdates(context.Background(), time.Now())

			assert.Equal(t, tt.expectedResult, result)
			client.AssertExpectations(t)
			userDataMgr.AssertExpectations(t)
		})
	}
}
