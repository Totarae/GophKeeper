package repository

import (
	"database/sql"
	"errors"
	"github.com/Totarae/GophKeeper/internal/server/model"
)

// Самая примитивная репа, пока только C/R

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(login, passwordHash, masterPasswordHash string) error {
	_, err := r.db.Exec(
		queryCreateUser,
		login, passwordHash, masterPasswordHash,
	)
	return err
}

func (r *UserRepository) GetUserByLogin(login string) (*model.User, error) {
	var u model.User
	err := r.db.QueryRow(
		queryGetUserByLogin,
		login,
	).Scan(&u.ID, &u.Login, &u.PasswordHash, &u.MasterPasswordHash)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}
