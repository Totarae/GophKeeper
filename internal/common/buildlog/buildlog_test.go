package buildlog_test

import (
	"github.com/Totarae/GophKeeper/internal/common/buildlog"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"strings"
	"testing"
)

func TestPrint(t *testing.T) {
	// Перехватываем stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	buildlog.Print("v1.0.0", "2025-07-29", "abc123")

	// Завершаем захват вывода
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = old

	output := string(out)

	require.True(t, strings.Contains(output, "Build version: v1.0.0"))
	require.True(t, strings.Contains(output, "Build date: 2025-07-29"))
	require.True(t, strings.Contains(output, "Build commit: abc123"))
}
