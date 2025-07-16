package manager

import (
	"GophKeeper/internal/server/jwt"
	"GophKeeper/internal/server/model"
	"GophKeeper/internal/server/repository"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserExists = errors.New("user already exists")
)

type UserRepository interface {
	GetUserByLogin(login string) (*model.User, error)
	CreateUser(login, passwordHash, masterPasswordHash string) error
}

type UserManager struct {
	userRepo UserRepository
	jwt      *jwt.Container
}

func NewUserManager(
	userRepo *repository.UserRepository,
	jwt *jwt.Container,
) *UserManager {
	return &UserManager{
		userRepo: userRepo,
		jwt:      jwt,
	}
}

func (m *UserManager) Register(login, password, masterPassword string) (string, error) {
	existing, err := m.userRepo.GetUserByLogin(login)
	if err != nil {
		return "", err
	}
	if existing != nil {
		return "", ErrUserExists
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	masterPasswordHash, err := bcrypt.GenerateFromPassword([]byte(masterPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	err = m.userRepo.CreateUser(login, string(passwordHash), string(masterPasswordHash))
	if err != nil {
		return "", err
	}

	user, err := m.userRepo.GetUserByLogin(login)
	if err != nil || user == nil {
		return "", err
	}

	return m.jwt.Encode(user.ID, login)
}

func (m *UserManager) DecodeToken(token string) (*jwt.Claims, error) {
	return m.jwt.Decode(token)
}
