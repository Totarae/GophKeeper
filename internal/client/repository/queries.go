package repository

const (
	queryCreateMetaTable = `
		CREATE TABLE IF NOT EXISTS meta (
			id INT PRIMARY KEY CHECK (id = 0),
			last_sync BIGINT NOT NULL,
			master_password_hash TEXT NOT NULL
		)
	`

	queryInsertInitialMeta = `
		INSERT INTO meta (id, last_sync, master_password_hash)
		VALUES (0, 0, '')
		ON CONFLICT (id) DO NOTHING
	`

	queryGetLastSync = `SELECT last_sync FROM meta WHERE id = 0`
	querySetLastSync = `UPDATE meta SET last_sync = $1 WHERE id = 0`

	queryGetMasterPasswordHash = `SELECT master_password_hash FROM meta WHERE id = 0`
	querySetMasterPasswordHash = `UPDATE meta SET master_password_hash = $1 WHERE id = 0`

	queryCreateUserDataTable = `
		CREATE TABLE IF NOT EXISTS user_data (
			id SERIAL PRIMARY KEY,
			data_key TEXT NOT NULL UNIQUE,
			data_value BYTEA NOT NULL,
			updated_at BIGINT NOT NULL,
			deleted_at BIGINT NOT NULL
		)
	`

	queryCreateIndexUpdatedAt = `CREATE INDEX IF NOT EXISTS idx_updated_at ON user_data(updated_at)`
	queryCreateIndexDeletedAt = `CREATE INDEX IF NOT EXISTS idx_deleted_at ON user_data(deleted_at)`

	queryMergeUserData = `
		INSERT INTO user_data (data_key, data_value, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (data_key) DO UPDATE SET
			data_value = EXCLUDED.data_value,
			updated_at = EXCLUDED.updated_at,
			deleted_at = EXCLUDED.deleted_at
	`

	queryGetUserData = `
		SELECT id, data_key, data_value, updated_at, deleted_at
		FROM user_data WHERE data_key = $1
	`

	queryGetUserDataUpdates = `
		SELECT id, data_key, data_value, updated_at, deleted_at
		FROM user_data WHERE updated_at > $1
	`
)
