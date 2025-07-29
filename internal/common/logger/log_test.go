package logger_test

import (
	"bytes"
	"github.com/Totarae/GophKeeper/internal/common/logger"
	"go.uber.org/zap"
	"io"
	"os"
	"strings"
	"testing"
)

func captureOutput(f func()) string {
	old := os.Stdout

	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	_ = w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}

func TestPrintBuildInfo(t *testing.T) {
	output := captureOutput(func() {
		logger.PrintBuildInfo("v1.2.3", "2025-07-29", "abcdef123")
	})

	if !strings.Contains(output, "v1.2.3") ||
		!strings.Contains(output, "2025-07-29") ||
		!strings.Contains(output, "abcdef123") {
		t.Errorf("PrintBuildInfo output mismatch:\n%s", output)
	}
}

func TestInit(t *testing.T) {
	t.Cleanup(func() {
		logger.Logger = zap.NewNop() // сброс
	})

	logger.Init("test-name", "info")
	if logger.Logger == nil {
		t.Fatal("logger was not initialized")
	}
}
