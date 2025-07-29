package app

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Totarae/GophKeeper/internal/server/config"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func Test_initDB(t *testing.T) {
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	}

	cfg := &config.Config{
		DatabaseDSN: dsn,
	}

	db, err := initDB(cfg)
	require.NoError(t, err)
	require.NotNil(t, db)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	require.NoError(t, err)

	_ = db.Close()
}

func TestNew(t *testing.T) {
	// Задаем окружение для config.NewConfig()
	t.Setenv("APP_SECRET", "mysecret")
	t.Setenv("DATABASE_DSN", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	t.Setenv("PORT", ":50051")

	app, err := New()
	require.NoError(t, err)
	require.NotNil(t, app)

	require.NotNil(t, app.db)
	require.NotNil(t, app.services)
	require.NotNil(t, app.services.userManager)
	require.NotNil(t, app.services.dataManager)

	err = app.db.Ping()
	require.NoError(t, err)

	_ = app.db.Close()
}
