package app_test

import (
	"github.com/Totarae/GophKeeper/internal/client/app"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewApp_Success(t *testing.T) {

	t.Setenv("APP_ENV", "test")

	t.Cleanup(func() { os.Args = os.Args[:1] })
	os.Args = []string{
		"test",
		"-addr=localhost:50501",
		"-interval=1",
		"postgres://postgres:postgres@localhost:5432/test?sslmode=disable",
		"master123",
	}

	a, err := app.New()
	require.NoError(t, err)
	require.NotNil(t, a)
}

func TestApp_Run_CancelImmediately(t *testing.T) {

	t.Cleanup(func() { os.Args = os.Args[:1] })
	os.Args = []string{
		"test",
		"-addr=localhost:50501",
		"-interval=1",
		"postgres://postgres:postgres@localhost:5432/test?sslmode=disable",
		"master123",
	}

	a, err := app.New()
	require.NoError(t, err)
	require.NotNil(t, a)

	done := make(chan struct{})
	go func() {
		a.Run()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Run did not return in time")
	}
}
