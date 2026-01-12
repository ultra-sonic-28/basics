package logger

import "log/slog"

// Niveaux personnalisés
const (
	LevelInfo     slog.Level = 0
	LevelDebug    slog.Level = 1
	LevelWarning  slog.Level = 2
	LevelCritical slog.Level = 3
	LevelFatal    slog.Level = 4
)

func levelToString(lvl slog.Level) string {
	switch lvl {
	case LevelInfo:
		return "INFO"
	case LevelDebug:
		return "DEBUG"
	case LevelWarning:
		return "WARNING"
	case LevelCritical:
		return "CRITICAL"
	case LevelFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}
