package repository

const (
	queryCreateUser = `
		INSERT INTO "user" (login, password_hash, master_password_hash)
		VALUES ($1, $2, $3)
	`

	queryGetUserByLogin = `
		SELECT id, login, password_hash, master_password_hash
		FROM "user" WHERE login = $1
	`

	querySelectUserDataUpdatedAt = `
		SELECT updated_at FROM user_data WHERE user_id = $1 AND data_key = $2
	`

	queryInsertUserData = `
		INSERT INTO user_data (user_id, data_key, data_value, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	queryUpdateUserData = `
		UPDATE user_data
		SET data_value = $1, updated_at = $2, deleted_at = $3
		WHERE user_id = $4 AND data_key = $5
	`

	queryGetUserUpdates = `
		SELECT id, user_id, data_key, data_value, updated_at, deleted_at
		FROM user_data
		WHERE user_id = $1 AND srv_updated_at > $2
	`
)
