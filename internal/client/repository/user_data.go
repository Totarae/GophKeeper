package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Totarae/GophKeeper/internal/client/model"
	"time"

	_ "github.com/lib/pq" // PostgreSQL драйвер
)

type UserDataRepository struct {
	db *sql.DB
}

func NewUserDataRepository(db *sql.DB) (*UserDataRepository, error) {
	repo := &UserDataRepository{db: db}
	if err := repo.init(); err != nil {
		return nil, err
	}
	return repo, nil
}

func (r *UserDataRepository) Merge(ctx context.Context, data *model.UserData) error {
	_, err := r.db.ExecContext(
		ctx, queryMergeUserData,
		data.DataKey,
		data.DataValue,
		data.UpdatedAt.Unix(),
		data.DeletedAt.Unix(),
	)
	return err
}

func (r *UserDataRepository) Get(ctx context.Context, key string) (*model.UserData, error) {
	row := r.db.QueryRowContext(ctx, queryGetUserData, key)
	d := &model.UserData{}

	var updatedAt, deletedAt int64
	err := row.Scan(&d.ID, &d.DataKey, &d.DataValue, &updatedAt, &deletedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	d.UpdatedAt = time.Unix(updatedAt, 0)
	d.DeletedAt = time.Unix(deletedAt, 0)

	return d, nil
}

func (r *UserDataRepository) GetUpdates(ctx context.Context, after time.Time) ([]*model.UserData, error) {
	rows, err := r.db.QueryContext(ctx, queryGetUserDataUpdates, after.Unix())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*model.UserData
	for rows.Next() {
		var ud model.UserData
		var updatedAt, deletedAt int64
		if err := rows.Scan(&ud.ID, &ud.DataKey, &ud.DataValue, &updatedAt, &deletedAt); err != nil {
			return nil, err
		}
		ud.UpdatedAt = time.Unix(updatedAt, 0)
		ud.DeletedAt = time.Unix(deletedAt, 0)
		result = append(result, &ud)
	}
	return result, rows.Err()
}

func (r *UserDataRepository) init() error {
	queries := []string{
		queryCreateUserDataTable,
		queryCreateIndexUpdatedAt,
		queryCreateIndexDeletedAt,
	}

	for _, query := range queries {
		if _, err := r.db.Exec(query); err != nil {
			return err
		}
	}
	return nil
}
