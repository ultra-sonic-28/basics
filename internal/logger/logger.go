package logger

import (
	"log/slog"
	"os"
)

func NewFileLogger(path, appName string, level slog.Level) (*slog.Logger, func() error, error) {
	file, err := os.OpenFile(
		path,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		return nil, nil, err
	}

	handler := NewTextHandler(file, level, appName)

	closeFn := func() error {
		return file.Close()
	}

	return slog.New(handler), closeFn, nil
}

// InitLogger initialise le logger pour l'application.
// path : chemin du fichier de log
// appName : nom de l'application pour le log
// level : niveau minimum à logger
func InitLogger(path, appName string, level slog.Level) (func() error, error) {
	log, closeFn, err := NewFileLogger(path, appName, level)
	if err != nil {
		return nil, err
	}

	// Définit ce logger comme logger par défaut
	slog.SetDefault(log)
	return closeFn, nil
}
