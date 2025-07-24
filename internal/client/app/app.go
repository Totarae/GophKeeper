package app

import (
	"GophKeeper/internal/client/cli"
	"GophKeeper/internal/client/config"
	"GophKeeper/internal/client/syncro"
	"GophKeeper/internal/common/logger"
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
	db       *sql.DB
	registry cli.CommandRegistry
	syncer   *syncro.Synchronizer
}

func New() (*App, error) {
	conf, err := config.ParseArgs()
	if err != nil {
		return nil, fmt.Errorf("can`t parse arguments: %w", err)
	}

	logger.Init("client", zap.InfoLevel.String())

	dbPath := conf.DBPath
	if err := touchFilepath(dbPath); err != nil {
		return nil, fmt.Errorf("can`t touch filepath: %w", err)
	}

	db, err := sql.Open("postgres", dbPath)
	if err != nil {
		return nil, fmt.Errorf("can't open db: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	client, err := grpc.NewClient(conf.ServerAddr)
	if err != nil {
		return nil, fmt.Errorf("can`t create client: %w", err)
	}

	syncer := syncro.New(client, userDataManager, metaManager, time.Duration(conf.SyncIntervalSec)*time.Second)

	return &App{
		syncer: syncer,
		registry: cli.CommandRegistry{
			"get":      command.NewGetCommand(userDataManager, []byte(conf.MasterPassword)),
			"set":      command.NewSetCommand(userDataManager, []byte(conf.MasterPassword)),
			"login":    command.NewLoginCommand(client, []byte(conf.MasterPassword)),
			"register": command.NewRegisterCommand(client, []byte(conf.MasterPassword)),
		},
		db: db,
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
