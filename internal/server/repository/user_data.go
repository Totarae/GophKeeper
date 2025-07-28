package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Totarae/GophKeeper/internal/server/model"
	"time"
)

type UserDataRepository struct {
	db *sql.DB
}

func NewUserDataRepository(db *sql.DB) *UserDataRepository {
	return &UserDataRepository{db: db}
}

func (r *UserDataRepository) Merge(ctx context.Context, data *model.UserData) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var currentVersion time.Time
	err = tx.QueryRowContext(ctx,
		querySelectUserDataUpdatedAt,
		data.UserID, data.DataKey,
	).Scan(&currentVersion)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if errors.Is(err, sql.ErrNoRows) {
		_, err = tx.ExecContext(ctx, queryInsertUserData, data.UserID, data.DataKey, data.DataValue, data.UpdatedAt, data.DeletedAt)

		if err != nil {
			return err
		}
	} else {
		if data.UpdatedAt.Before(currentVersion) {
			return nil
		}

		_, err := tx.ExecContext(ctx, queryUpdateUserData, data.DataValue, data.UpdatedAt, data.DeletedAt, data.UserID, data.DataKey)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *UserDataRepository) GetUpdates(ctx context.Context, userID uint32, since time.Time) ([]*model.UserData, error) {
	rows, err := r.db.QueryContext(ctx,
		queryGetUserUpdates,
		userID, since,
	)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	defer rows.Close()

	var result []*model.UserData
	for rows.Next() {
		d := &model.UserData{}
		err = rows.Scan(
			&d.ID,
			&d.UserID,
			&d.DataKey,
			&d.DataValue,
			&d.UpdatedAt,
			&d.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, d)
	}
	return result, nil
}
