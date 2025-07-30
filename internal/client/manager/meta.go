package manager

import (
	"context"
	"github.com/Totarae/GophKeeper/internal/client/repository"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type MetaReader interface {
	GetLastSync(ctx context.Context) (time.Time, error)
	GetMasterPasswordHash(ctx context.Context) (string, error)
}

type MetaWriter interface {
	SetLastSync(ctx context.Context, t time.Time) error
	SetMasterPasswordHash(ctx context.Context, h string) error
}

type MetaRepository interface {
	MetaReader
	MetaWriter
}

type MetaManager struct {
	repo MetaRepository
}

func NewMetaManager(repo *repository.MetaRepository) *MetaManager {
	return &MetaManager{repo: repo}
}

func (m *MetaManager) GetLastSync(ctx context.Context) (time.Time, error) {
	return m.repo.GetLastSync(ctx)
}

func (m *MetaManager) SetLastSync(ctx context.Context, t time.Time) error {
	return m.repo.SetLastSync(ctx, t)
}

func (m *MetaManager) MasterPasswordHashDefined(ctx context.Context) (bool, error) {
	h, err := m.repo.GetMasterPasswordHash(ctx)
	if err != nil {
		return false, err
	}

	return h != "", nil
}

func (m *MetaManager) ValidateMasterPassword(ctx context.Context, password string) error {
	h, err := m.repo.GetMasterPasswordHash(ctx)
	if err != nil {
		return err
	}

	return bcrypt.CompareHashAndPassword([]byte(h), []byte(password))
}

func (m *MetaManager) SetMasterPassword(ctx context.Context, password string) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return m.repo.SetMasterPasswordHash(ctx, string(passwordHash))
}
