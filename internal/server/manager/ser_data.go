package manager

import (
	"GophKeeper/internal/server/model"
	"GophKeeper/internal/server/repository"
	"context"
	"time"
)

type UserDataRepository interface {
	Merge(ctx context.Context, data *model.UserData) error
	GetUpdates(ctx context.Context, userID uint32, since time.Time) ([]*model.UserData, error)
}

type UserDataManager struct {
	dataRepo UserDataRepository
}

func NewUserDataManager(dataRepo *repository.UserDataRepository) *UserDataManager {
	return &UserDataManager{
		dataRepo: dataRepo,
	}
}

func (m *UserDataManager) Merge(ctx context.Context, data *model.UserData) error {
	return m.dataRepo.Merge(ctx, data)
}

func (m *UserDataManager) GetUpdates(ctx context.Context, userID uint32, since time.Time) ([]*model.UserData, error) {
	return m.dataRepo.GetUpdates(ctx, userID, since)
}
