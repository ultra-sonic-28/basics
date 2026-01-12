package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"sync"
)

type TextHandler struct {
	mu      sync.Mutex
	w       io.Writer
	level   slog.Level
	appName string
}

func NewTextHandler(w io.Writer, level slog.Level, appName string) *TextHandler {
	return &TextHandler{
		w:       w,
		level:   level,
		appName: appName,
	}
}

func (h *TextHandler) Close() error {
	if c, ok := h.w.(io.Closer); ok {
		return c.Close()
	}
	return nil
}

func (h *TextHandler) Enabled(_ context.Context, lvl slog.Level) bool {
	return lvl >= h.level
}

func (h *TextHandler) Handle(_ context.Context, r slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	timestamp := r.Time.Format("2006-01-02 15:04:05,000")
	level := levelToString(r.Level)

	line := fmt.Sprintf(
		"%s [%s] %s: %s\n",
		timestamp,
		level,
		h.appName,
		r.Message,
	)

	_, err := h.w.Write([]byte(line))
	return err
}

func (h *TextHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}

func (h *TextHandler) WithGroup(_ string) slog.Handler {
	return h
}
