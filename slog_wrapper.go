package commonutils

import (
	"errors"
	"io"
	"log/slog"
)

type LogLevel string

var (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

var ErrInvalidLogLevel = errors.New("invalid log level")

func ConvertLevelStringToSlogLevel(levelStr string) (slog.Level, error) {
	var slogLevel slog.Level
	switch levelStr {
	case string(LogLevelDebug):
		slogLevel = slog.LevelDebug
	case string(LogLevelInfo):
		slogLevel = slog.LevelInfo
	case string(LogLevelWarn):
		slogLevel = slog.LevelWarn
	case string(LogLevelError):
		slogLevel = slog.LevelError
	default:
		return 0, ErrInvalidLogLevel
	}
	return slogLevel, nil
}

func UseTextLogger(level slog.Level, w io.Writer) *slog.Logger {
	handler := slog.NewTextHandler(w, &slog.HandlerOptions{
		Level: level,
	})
	logger := slog.New(handler)
	slog.SetLogLoggerLevel(level)
	slog.SetDefault(logger)
	return logger
}

func UseJSONLogger(level slog.Level, w io.Writer) *slog.Logger {
	handler := slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: level,
	})
	logger := slog.New(handler)
	slog.SetLogLoggerLevel(level)
	slog.SetDefault(logger)
	return logger
}
