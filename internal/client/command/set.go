package command

import (
	"context"
	"errors"
	"github.com/Totarae/GophKeeper/internal/client/aes"
	"github.com/Totarae/GophKeeper/internal/client/model"
	value "github.com/Totarae/GophKeeper/internal/client/values"
	"time"
)

type DataMerger interface {
	Merge(ctx context.Context, data *model.UserData) error
}

type SetCommand struct {
	dataManager    DataMerger
	masterPassword []byte
}

func NewSetCommand(dataManager DataMerger, masterPassword []byte) *SetCommand {
	return &SetCommand{
		dataManager:    dataManager,
		masterPassword: masterPassword,
	}
}

func (c *SetCommand) Execute(ctx context.Context, args []string) (string, error) {
	if len(args) < 3 {
		return "", errors.New("args: <key> <type> <args...>")
	}

	val, err := value.FromUserInput(args[1], args[2:])
	if err != nil {
		return "", err
	}

	raw, err := val.ToBytes()
	if err != nil {
		return "", err
	}

	encRaw, err := aes.Encrypt(c.masterPassword, raw)
	if err != nil {
		return "", err
	}

	err = c.dataManager.Merge(ctx, &model.UserData{
		DataKey:   args[0],
		DataValue: encRaw,
		UpdatedAt: time.Now(),
		DeletedAt: time.Unix(0, 0),
	})

	if err != nil {
		return "", err
	}

	return "saved successful", nil
}
