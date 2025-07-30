-- +goose Up
CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       login VARCHAR(255) UNIQUE NOT NULL,
                       password_hash VARCHAR(255) NOT NULL,
                       master_password_hash VARCHAR(255) NOT NULL
);

CREATE TABLE user_data (
                           id SERIAL PRIMARY KEY,
                           user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                           data_key VARCHAR(255) NOT NULL,
                           data_value BYTEA NOT NULL,
                           version INTEGER NOT NULL DEFAULT 1,
                           srv_updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                           updated_at TIMESTAMP NOT NULL,
                           deleted_at TIMESTAMP,
                           UNIQUE (user_id, data_key)
);

CREATE INDEX idx_updated_at ON user_data(updated_at);
CREATE INDEX idx_deleted_at ON user_data(deleted_at);

-- +goose Down
DROP INDEX IF EXISTS idx_deleted_at;
DROP INDEX IF EXISTS idx_updated_at;
DROP TABLE IF EXISTS user_data;
DROP TABLE IF EXISTS users;
