package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Totarae/GophKeeper/internal/client/cli"
	"github.com/Totarae/GophKeeper/internal/client/command"
	"github.com/Totarae/GophKeeper/internal/client/config"
	"github.com/Totarae/GophKeeper/internal/client/grpc"
	"github.com/Totarae/GophKeeper/internal/client/manager"
	"github.com/Totarae/GophKeeper/internal/client/repository"
	"github.com/Totarae/GophKeeper/internal/client/syncro"
	"github.com/Totarae/GophKeeper/internal/common/logger"
	"go.uber.org/zap"
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
	conf, err := config.NewConfig()
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

	userDataRepo, err := repository.NewUserDataRepository(db)
	if err != nil {
		return nil, fmt.Errorf("can`t create data repo: %w", err)
	}

	metaRepo, err := repository.NewMetaRepository(db)
	if err != nil {
		return nil, fmt.Errorf("can`t create meta repo: %w", err)
	}

	userDataManager := manager.NewUserDataManager(userDataRepo)
	metaManager := manager.NewMetaManager(metaRepo)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	ok, err := metaManager.MasterPasswordHashDefined(ctx)
	if err != nil {
		return nil, err
	}
	if !ok {
		if err := metaManager.SetMasterPassword(ctx, conf.MasterPassword); err != nil {
			return nil, err
		}
	} else {
		if err := metaManager.ValidateMasterPassword(ctx, conf.MasterPassword); err != nil {
			return nil, fmt.Errorf("invalid master password: %w", err)
		}
	}

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
	defer a.syncer.Stop()

	var wg sync.WaitGroup
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	wg.Add(2)
	go func() {
		defer wg.Done()
		a.syncer.Start(ctx)
		stop()
	}()
	go func() {
		defer wg.Done()
		cli.Run(ctx, a.registry)
		stop()
	}()

	wg.Wait()

}
