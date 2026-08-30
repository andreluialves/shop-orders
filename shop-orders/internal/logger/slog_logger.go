package logger

import (
	"log/slog"
	"os"
)

type SlogLogger struct {
	logger *slog.Logger
}

func NewSlogLogger() *SlogLogger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	return &SlogLogger{logger: slog.New(handler)}
}

func (s *SlogLogger) Info(msg string, args ...any)  { s.logger.Info(msg, args...) }
func (s *SlogLogger) Warn(msg string, args ...any)  { s.logger.Warn(msg, args...) }
func (s *SlogLogger) Error(msg string, args ...any) { s.logger.Error(msg, args...) }
func (s *SlogLogger) Debug(msg string, args ...any) { s.logger.Debug(msg, args...) }
