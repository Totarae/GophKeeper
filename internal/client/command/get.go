package command

import (
	"GophKeeper/internal/client/aes"
	"GophKeeper/internal/client/model"
	"context"
	"errors"
)

type UserDataGetter interface {
	Get(ctx context.Context, key string) (*model.UserData, error)
}

type GetCommand struct {
	dataManager    UserDataGetter
	masterPassword []byte
}

func NewGetCommand(dataManager UserDataGetter, masterPassword []byte) *GetCommand {
	return &GetCommand{
		dataManager:    dataManager,
		masterPassword: masterPassword,
	}
}

func (c *GetCommand) Execute(ctx context.Context, args []string) (string, error) {
	if len(args) < 1 {
		return "", errors.New("args: <key>")
	}

	data, err := c.dataManager.Get(ctx, args[0])
	if err != nil {
		return "", err
	}

	raw, err := aes.Decrypt(c.masterPassword, data.DataValue)
	if err != nil {
		return "", err
	}

	val, err := value.FromBytes(raw)
	if err != nil {
		return "", err
	}

	return val.String(), nil
}
