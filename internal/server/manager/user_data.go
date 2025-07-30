package manager

import (
	"context"
	"github.com/Totarae/GophKeeper/internal/server/model"
	"github.com/Totarae/GophKeeper/internal/server/repository"
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
