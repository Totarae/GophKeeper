package repository

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetaRepository_Init(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectExec(`CREATE TABLE IF NOT EXISTS meta`).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(`INSERT INTO meta \(id, last_sync, master_password_hash\) VALUES \(0, 0, ''\) ON CONFLICT \(id\) DO NOTHING`).
		WillReturnResult(sqlmock.NewResult(0, 1))

	repo := &MetaRepository{db: db}
	err = repo.init()

	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMetaRepository_GetLastSync(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	expectedTime := time.Now().Unix()
	mock.ExpectQuery(`SELECT last_sync FROM meta WHERE id = 0`).
		WillReturnRows(sqlmock.NewRows([]string{"last_sync"}).AddRow(expectedTime))

	repo := &MetaRepository{db: db}
	ctx := context.Background()

	ts, err := repo.GetLastSync(ctx)
	require.NoError(t, err)
	assert.Equal(t, expectedTime, ts.Unix())
}

func TestMetaRepository_SetLastSync(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	timestamp := time.Now().Unix()
	mock.ExpectExec(`UPDATE meta SET last_sync = \$1 WHERE id = 0`).
		WithArgs(timestamp).
		WillReturnResult(sqlmock.NewResult(0, 1))

	repo := &MetaRepository{db: db}
	ctx := context.Background()

	err = repo.SetLastSync(ctx, time.Unix(timestamp, 0))
	require.NoError(t, err)
}

func TestMetaRepository_GetMasterPasswordHash(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	expectedHash := "hash123"
	mock.ExpectQuery(`SELECT master_password_hash FROM meta WHERE id = 0`).
		WillReturnRows(sqlmock.NewRows([]string{"master_password_hash"}).AddRow(expectedHash))

	repo := &MetaRepository{db: db}
	ctx := context.Background()

	hash, err := repo.GetMasterPasswordHash(ctx)
	require.NoError(t, err)
	assert.Equal(t, expectedHash, hash)
}

func TestMetaRepository_SetMasterPasswordHash(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	hash := "hash456"
	mock.ExpectExec(`UPDATE meta SET master_password_hash = \$1 WHERE id = 0`).
		WithArgs(hash).
		WillReturnResult(sqlmock.NewResult(0, 1))

	repo := &MetaRepository{db: db}
	ctx := context.Background()

	err = repo.SetMasterPasswordHash(ctx, hash)
	require.NoError(t, err)
}
