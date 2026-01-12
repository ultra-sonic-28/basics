package logger

import (
	"basics/testutils"
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

func TestNewFileLogger(t *testing.T) {
	tests := []struct {
		name    string
		level   slog.Level
		message string
		appName string
	}{
		{
			name:    "Info level logging",
			level:   LevelInfo,
			message: "file logger initialized",
			appName: "test_app",
		},
		{
			name:    "Debug level logging",
			level:   LevelDebug,
			message: "debug message",
			appName: "test_app",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			logPath := filepath.Join(dir, "test.log")

			log, closeFn, err := NewFileLogger(logPath, tt.appName, tt.level)
			msg := fmt.Sprintf("NewFileLogger returned error: %v", err)
			testutils.True(t, msg, err == nil)

			t.Cleanup(func() {
				_ = closeFn()
			})

			log.Log(context.Background(), tt.level, tt.message)

			data, err := os.ReadFile(logPath)
			msg = fmt.Sprintf("failed to read log file: %v", err)
			testutils.True(t, msg, err == nil)

			got := string(data)
			testutils.True(t, "expected log output, got empty file", got != "")

			pattern := `^\d{4}-\d{2}-\d{2} ` +
				`\d{2}:\d{2}:\d{2},\d{3} ` +
				`\[` + levelToString(tt.level) + `\] ` +
				tt.appName + `: ` + tt.message + `\n$`

			matched, err := regexp.MatchString(pattern, got)
			msg = fmt.Sprintf("regexp.MatchString failed: %v", err)
			testutils.True(t, msg, err == nil)

			msg = fmt.Sprintf(
				"log file content mismatch\nGot: %q\nExpected pattern: %q",
				got,
				pattern)
			testutils.True(t, msg, matched)
		})
	}
}

func TestInitLogger(t *testing.T) {
	tests := []struct {
		name    string
		level   slog.Level
		message string
		appName string
	}{
		{
			name:    "InitLogger with Info",
			level:   LevelInfo,
			message: "application started",
			appName: "test_app",
		},
		{
			name:    "InitLogger with Debug",
			level:   LevelDebug,
			message: "debug enabled",
			appName: "test_app",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			logPath := filepath.Join(dir, "app.log")

			closeLogger, err := InitLogger(logPath, tt.appName, tt.level)
			msg := fmt.Sprintf("InitLogger returned error: %v", err)
			testutils.True(t, msg, err == nil)

			t.Cleanup(func() {
				_ = closeLogger()
			})

			// Utilisation du logger par défaut
			slog.Log(context.Background(), tt.level, tt.message)

			data, err := os.ReadFile(logPath)
			msg = fmt.Sprintf("failed to read log file: %v", err)
			testutils.True(t, msg, err == nil)

			got := string(data)
			testutils.True(t, "expected log output, got empty file", got != "")

			pattern := `^\d{4}-\d{2}-\d{2} ` +
				`\d{2}:\d{2}:\d{2},\d{3} ` +
				`\[` + levelToString(tt.level) + `\] ` +
				tt.appName + `: ` + tt.message + `\n$`

			matched, err := regexp.MatchString(pattern, got)
			msg = fmt.Sprintf("regexp.MatchString failed: %v", err)
			testutils.True(t, msg, err == nil)

			msg = fmt.Sprintf(
				"log output mismatch\nGot: %q\nExpected pattern: %q",
				got,
				pattern)
			testutils.True(t, msg, matched)
		})
	}
}
