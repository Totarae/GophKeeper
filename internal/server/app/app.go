package app

import (
	"GophKeeper/internal/server/config"
	"context"
	"database/sql"
	"fmt"
	"go.uber.org/zap"
	"net"
	"os/signal"
	"syscall"
	"time"
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

func (a *App) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	//grpc server

	listener, err := net.Listen("tcp", a.cfg.Port)
	if err != nil {
		return fmt.Errorf("failed to create listener: %w", err)
	}

	go func() {
		logger.Logger.Info("Starting gRPC server", zap.String("addr", a.cfg.Listen))

		stop()
	}()

	<-ctx.Done()
	stop()

	logger.Logger.Info("Shutting down server...")

	logger.Logger.Info("Server stopped gracefully")

	return nil
}

func initDB(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.DatabaseDSN+"?parseTime=true")
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	return db, nil
}

type services struct {
	userManager *manager.UserManager
	dataManager *manager.UserDataManager
}

func initServices(db *sql.DB, cfg *config.Config) *services {
	userRepo := repository.NewUserRepository(db)
	dataRepo := repository.NewUserDataRepository(db)

	return &services{
		userManager: manager.NewUserManager(userRepo, jwt.New(cfg.AppSecret)),
		dataManager: manager.NewUserDataManager(dataRepo),
	}
}
