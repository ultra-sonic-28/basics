package logger

import (
	"basics/testutils"
	"bytes"
	"fmt"
	"log/slog"
	"regexp"
	"testing"
)

func TestHelpers_LogOutput(t *testing.T) {
	tests := []struct {
		name    string
		logFunc func(string, ...any)
		level   slog.Level
		message string
		appName string
	}{
		{
			name:    "Info helper",
			logFunc: Info,
			level:   LevelInfo,
			message: "info message",
			appName: "test_app",
		},
		{
			name:    "Debug helper",
			logFunc: Debug,
			level:   LevelDebug,
			message: "debug message",
			appName: "test_app",
		},
		{
			name:    "Warning helper",
			logFunc: Warning,
			level:   LevelWarning,
			message: "warning message",
			appName: "test_app",
		},
		{
			name:    "Critical helper",
			logFunc: Critical,
			level:   LevelCritical,
			message: "critical message",
			appName: "test_app",
		},
		{
			name:    "Fatal helper",
			logFunc: Fatal,
			level:   LevelFatal,
			message: "fatal message",
			appName: "test_app",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			handler := NewTextHandler(&buf, LevelInfo, tt.appName)
			log := slog.New(handler)

			// Remplacer le logger par défaut pour le test
			slog.SetDefault(log)

			// Appel du helper
			tt.logFunc(tt.message)

			got := buf.String()
			testutils.True(t, "expected log output, got empty string", got != "")

			// Vérification du format et du niveau
			pattern := `^\d{4}-\d{2}-\d{2} ` +
				`\d{2}:\d{2}:\d{2},\d{3} ` +
				`\[` + levelToString(tt.level) + `\] ` +
				tt.appName + `: ` + tt.message + `\n$`

			matched, err := regexp.MatchString(pattern, got)

			msg := fmt.Sprintf("regexp.MatchString failed: %v", err)
			testutils.True(t, msg, err == nil)

			msg = fmt.Sprintf(
				"log output does not match expected format\nGot: %q\nExpected pattern: %q",
				got,
				pattern)
			testutils.True(t, msg, matched)
		})
	}
}
