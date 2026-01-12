package logger

import (
	"basics/testutils"
	"log/slog"
	"path/filepath"
	"testing"
)

func TestInitLogger_OpenFileError(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()

	invalidPath := filepath.Join(
		tmpDir,
		"does",
		"not",
		"exist",
		"test.log",
	)

	// Sauvegarde du logger par défaut pour comparaison
	originalLogger := slog.Default()

	closeFn, err := InitLogger(
		invalidPath,
		"test_app",
		slog.LevelInfo,
	)

	testutils.True(t, "expected error, got nil", err != nil)
	testutils.True(t, "expected closeFn to be nil on error", closeFn == nil)
	testutils.True(t, "default logger should not be modified on InitLogger error", slog.Default() == originalLogger)
}
