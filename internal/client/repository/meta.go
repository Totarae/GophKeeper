package repository

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq" // PostgreSQL драйвер
)

type MetaRepository struct {
	db *sql.DB
}

func NewMetaRepository(db *sql.DB) (*MetaRepository, error) {
	repo := &MetaRepository{db: db}
	if err := repo.init(); err != nil {
		return nil, err
	}
	return repo, nil
}

func (r *MetaRepository) GetLastSync(ctx context.Context) (time.Time, error) {
	var tsInt int64
	err := r.db.QueryRowContext(ctx, queryGetLastSync).Scan(&tsInt)
	return time.Unix(tsInt, 0).UTC(), err
}

func (r *MetaRepository) SetLastSync(ctx context.Context, t time.Time) error {
	_, err := r.db.ExecContext(ctx, querySetLastSync, t.Unix())
	return err
}

func (r *MetaRepository) GetMasterPasswordHash(ctx context.Context) (string, error) {
	var h string
	err := r.db.QueryRowContext(ctx, queryGetMasterPasswordHash).Scan(&h)
	return h, err
}

func (r *MetaRepository) SetMasterPasswordHash(ctx context.Context, h string) error {
	_, err := r.db.ExecContext(ctx, querySetMasterPasswordHash, h)
	return err
}

func (r *MetaRepository) init() error {
	for _, q := range []string{queryCreateMetaTable, queryInsertInitialMeta} {
		if _, err := r.db.Exec(q); err != nil {
			return err
		}
	}
	return nil
}
