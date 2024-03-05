package logger

import (
	"io"
	"log/slog"
)

type logger struct {
	slog *slog.Logger
}

var _ ILogger = (*logger)(nil)

var Log = logger{slog.New(slog.NewJSONHandler(io.Discard, nil))}

func Setup(lvl string, w io.Writer) {
	opts := &slog.HandlerOptions{
		Level: getLoggerLevel(lvl),
	}

	Log = logger{slog.New(slog.NewJSONHandler(w, opts))}
}

func (l *logger) Debug(msg string, args ...any) {
	l.slog.Debug(msg, args...)
}
func (l *logger) Info(msg string, args ...any) {
	l.slog.Info(msg, args...)
}
func (l *logger) Warn(msg string, args ...any) {
	l.slog.Warn(msg, args...)
}
func (l *logger) Error(msg string, args ...any) {
	l.slog.Error(msg, args...)
}

var loggerLevelMap = map[string]slog.Level{
	"debug": slog.LevelDebug,
	"info":  slog.LevelInfo,
	"warn":  slog.LevelWarn,
	"error": slog.LevelError,
}

func getLoggerLevel(c string) slog.Level {
	level, exist := loggerLevelMap[c]
	if !exist {
		return slog.LevelError
	}

	return level
}
