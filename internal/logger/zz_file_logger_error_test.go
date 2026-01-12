package logger

import (
	"basics/testutils"
	"log/slog"
	"path/filepath"
	"testing"
)

func TestNewFileLogger_OpenFileError(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()

	// Chemin volontairement invalide
	invalidPath := filepath.Join(
		tmpDir,
		"does",
		"not",
		"exist",
		"test.log",
	)

	logger, closeFn, err := NewFileLogger(
		invalidPath,
		"test_app",
		slog.LevelInfo,
	)

	testutils.True(t, "expected error, got nil", err != nil)
	testutils.True(t, "expected logger to be nil on error", logger == nil)
	testutils.True(t, "expected closeFn to be nil on error", closeFn == nil)
}
