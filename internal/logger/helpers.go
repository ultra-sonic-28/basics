package logger

import (
	"context"
	"log/slog"
)

// Helpers pour tous les niveaux personnalisés

func Info(msg string, args ...any) {
	slog.Log(context.Background(), LevelInfo, msg, args...)
}

func Debug(msg string, args ...any) {
	slog.Log(context.Background(), LevelDebug, msg, args...)
}

func Warning(msg string, args ...any) {
	slog.Log(context.Background(), LevelWarning, msg, args...)
}

func Critical(msg string, args ...any) {
	slog.Log(context.Background(), LevelCritical, msg, args...)
}

func Fatal(msg string, args ...any) {
	slog.Log(context.Background(), LevelFatal, msg, args...)
}
