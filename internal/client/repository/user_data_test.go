package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Totarae/GophKeeper/internal/client/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserDataRepository_Init(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectExec(`CREATE TABLE IF NOT EXISTS user_data`).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(`CREATE INDEX IF NOT EXISTS idx_updated_at ON user_data\(updated_at\)`).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(`CREATE INDEX IF NOT EXISTS idx_deleted_at ON user_data\(deleted_at\)`).WillReturnResult(sqlmock.NewResult(0, 1))

	repo := &UserDataRepository{db: db}
	err = repo.init()
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUserDataRepository_Upsert_Merge(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &UserDataRepository{db: db}
	ctx := context.Background()
	now := time.Now().UTC()

	mock.ExpectExec(`INSERT INTO user_data .* ON CONFLICT .* DO UPDATE SET .*`).
		WithArgs("test-key", []byte("test-value"), now.Unix(), int64(0)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Merge(ctx, &model.UserData{
		DataKey:   "test-key",
		DataValue: []byte("test-value"),
		UpdatedAt: now,
		DeletedAt: time.Unix(0, 0),
	})
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUserDataRepository_Upsert_UpdateExisting(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &UserDataRepository{db: db}
	ctx := context.Background()
	now := time.Now().UTC()

	// Первое добавление
	mock.ExpectExec(`INSERT INTO user_data .* ON CONFLICT .* DO UPDATE SET .*`).
		WithArgs("test-key", []byte("initial-value"), now.Unix(), int64(0)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Merge(ctx, &model.UserData{
		DataKey:   "test-key",
		DataValue: []byte("initial-value"),
		UpdatedAt: now,
		DeletedAt: time.Unix(0, 0),
	})
	require.NoError(t, err)

	// Обновление записи
	later := now.Add(time.Hour)
	mock.ExpectExec(`INSERT INTO user_data .* ON CONFLICT .* DO UPDATE SET .*`).
		WithArgs("test-key", []byte("updated-value"), later.Unix(), later.Unix()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Merge(ctx, &model.UserData{
		DataKey:   "test-key",
		DataValue: []byte("updated-value"),
		UpdatedAt: later,
		DeletedAt: later,
	})
	require.NoError(t, err)

	// Проверка, что все ожидания выполнены
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUserDataRepository_Get(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &UserDataRepository{db: db}
	ctx := context.Background()

	now := time.Now().Unix()

	mock.ExpectQuery(`SELECT id, data_key, data_value, updated_at, deleted_at FROM user_data WHERE data_key = \$1`).
		WithArgs("test-key").
		WillReturnRows(sqlmock.NewRows([]string{"id", "data_key", "data_value", "updated_at", "deleted_at"}).
			AddRow(1, "test-key", []byte("value"), now, int64(0)))

	result, err := repo.Get(ctx, "test-key")
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "test-key", result.DataKey)
	assert.Equal(t, []byte("value"), result.DataValue)
	assert.Equal(t, now, result.UpdatedAt.Unix())
	assert.Equal(t, int64(0), result.DeletedAt.Unix())
}

func TestUserDataRepository_Get_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &UserDataRepository{db: db}
	ctx := context.Background()

	mock.ExpectQuery(`SELECT id, data_key, data_value, updated_at, deleted_at FROM user_data WHERE data_key = \$1`).
		WithArgs("no-key").
		WillReturnError(sql.ErrNoRows)

	result, err := repo.Get(ctx, "no-key")
	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestUserDataRepository_GetUpdates(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &UserDataRepository{db: db}
	ctx := context.Background()

	now := time.Now().Unix()

	mock.ExpectQuery(`SELECT id, data_key, data_value, updated_at, deleted_at FROM user_data WHERE updated_at > \$1`).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "data_key", "data_value", "updated_at", "deleted_at"}).
			AddRow(1, "key1", []byte("val1"), now, int64(0)).
			AddRow(2, "key2", []byte("val2"), now, int64(0)))

	results, err := repo.GetUpdates(ctx, time.Now().Add(-1*time.Hour))
	require.NoError(t, err)
	require.Len(t, results, 2)
	assert.Equal(t, "key1", results[0].DataKey)
	assert.Equal(t, "key2", results[1].DataKey)
}
