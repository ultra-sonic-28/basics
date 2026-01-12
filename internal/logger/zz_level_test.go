package logger

import (
	"basics/testutils"
	"context"
	"fmt"
	"log/slog"
	"testing"
)

func TestLevelToString(t *testing.T) {
	tests := []struct {
		level    slog.Level
		expected string
	}{
		{LevelInfo, "INFO"},
		{LevelDebug, "DEBUG"},
		{LevelWarning, "WARNING"},
		{LevelCritical, "CRITICAL"},
		{LevelFatal, "FATAL"},
		{slog.Level(99), "UNKNOWN"}, // niveau inconnu
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			got := levelToString(tt.level)
			msg := fmt.Sprintf("levelToString(%v) = %q; want %q", tt.level, got, tt.expected)
			testutils.True(t, msg, got == tt.expected)
		})
	}
}

func TestTextHandler_Enabled(t *testing.T) {
	tests := []struct {
		handlerLevel slog.Level
		checkLevel   slog.Level
		expected     bool
		name         string
	}{
		{LevelInfo, LevelInfo, true, "Info >= Info"},
		{LevelInfo, LevelDebug, true, "Debug >= Info"},
		{LevelWarning, LevelInfo, false, "Info < Warning"},
		{LevelWarning, LevelWarning, true, "Warning >= Warning"},
		{LevelCritical, LevelFatal, true, "Fatal >= Critical"},
		{LevelFatal, LevelDebug, false, "Debug < Fatal"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &TextHandler{level: tt.handlerLevel}
			got := h.Enabled(context.Background(), tt.checkLevel)
			msg := fmt.Sprintf("Enabled(level=%v) = %v; want %v", tt.checkLevel, got, tt.expected)
			testutils.True(t, msg, got == tt.expected)
		})
	}
}
