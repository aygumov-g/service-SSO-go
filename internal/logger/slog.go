package logger

import (
	"log/slog"
	"os"
)

type SlogLogger struct {
	log *slog.Logger
}

func New() Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	return &SlogLogger{
		log: slog.New(handler),
	}
}

func (l *SlogLogger) Info(msg string, args ...any) {
	l.log.Info(msg, args...)
}

func (l *SlogLogger) Error(msg string, args ...any) {
	l.log.Error(msg, args...)
}
