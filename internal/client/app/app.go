package app

import (
	"GophKeeper/internal/client/config"
	"GophKeeper/internal/common/logger"
	"GophKeeper/internal/server/manager"
	"GophKeeper/internal/server/repository"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

type App struct {
	db *sql.DB
}

func New() (*App, error) {
	conf, err := config.ParseArgs()
	if err != nil {
		return nil, fmt.Errorf("can`t parse arguments: %w", err)
	}

	logger.Init("client", zap.InfoLevel.String())


	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	return &struct{
	}, nil
}

func touchFilepath(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(path)
		if err != nil {
			return err
		}

		return file.Close()
	}

	return nil
}

func (a *App) Run() {
	defer a.db.Close()


	// завершаем всё
	var wg sync.WaitGroup
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

}
