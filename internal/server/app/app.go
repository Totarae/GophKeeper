package app

import (
	"GophKeeper/internal/server/config"
	"database/sql"
	"fmt"
	"go.uber.org/zap"
)

type App struct {
	cfg      *config.Config
	db       *sql.DB
	services *services
	logger   *zap.Logger
}

func New() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	var logger *zap.Logger

	//var level string

	if cfg.Debug {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}

	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}
	logger = logger.Named("server")

	db, err := initDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to init db: %w", err)
	}

	logger.Info("initialization successful")

	return &App{
		cfg:      cfg,
		db:       db,
		services: initServices(db, cfg),
	}, nil
}
