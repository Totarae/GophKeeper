package grpc

import (
	"context"
	"github.com/Totarae/GophKeeper/internal/server/jwt"
	"github.com/Totarae/GophKeeper/internal/server/model"
	"time"
)

type UserManagerInterface interface {
	Register(login, password, masterPassword string) (string, error)
	Login(login, password, masterPassword string) (string, error)
	DecodeToken(token string) (*jwt.Claims, error)
}

type UserDataManagerInterface interface {
	Merge(ctx context.Context, data *model.UserData) error
	GetUpdates(ctx context.Context, userID uint32, updatedAfter time.Time) ([]*model.UserData, error)
}
