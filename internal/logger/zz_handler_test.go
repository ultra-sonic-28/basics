package logger

import (
	"basics/testutils"
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"testing"
	"time"
)

func TestTextHandler_Handle_Format(t *testing.T) {
	tests := []struct {
		level    slog.Level
		message  string
		appName  string
		expected string
		name     string
	}{
		{
			level:   LevelInfo,
			message: "Logging initialized",
			appName: "basics",
			name:    "Info level",
		},
		{
			level:   LevelDebug,
			message: "Debug mode activated",
			appName: "basics",
			name:    "Debug level",
		},
		{
			level:   LevelWarning,
			message: "Low disk space",
			appName: "basics",
			name:    "Warning level",
		},
		{
			level:   LevelCritical,
			message: "Database connection lost",
			appName: "basics",
			name:    "Critical level",
		},
		{
			level:   LevelFatal,
			message: "Application crash",
			appName: "basics",
			name:    "Fatal level",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			handler := NewTextHandler(&buf, LevelDebug, tt.appName)

			// Construire un Record
			record := slog.Record{
				Time:    time.Now(),
				Level:   tt.level,
				Message: tt.message,
			}

			err := handler.Handle(context.Background(), record)
			msg := fmt.Sprintf("Handle() returned error: %v", err)
			testutils.True(t, msg, err == nil)

			got := buf.String()

			// Regex pour vérifier le format
			// Ex: 2025-12-15 13:26:27,821 [INFO] rpg_companion: Logging initialized
			pattern := `^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2},\d{3} \[` +
				levelToString(tt.level) +
				`\] ` + tt.appName + `: ` + tt.message + `\n$`

			matched, err := regexp.MatchString(pattern, got)
			msg = fmt.Sprintf("regexp.MatchString failed: %v", err)
			testutils.True(t, msg, err == nil)

			msg = fmt.Sprintf("Log line mismatch.\nGot: %q\nExpected pattern: %q", got, pattern)
			testutils.True(t, msg, matched)
		})
	}
}
