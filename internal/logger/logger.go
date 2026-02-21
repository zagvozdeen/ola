package logger

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"sync"

	"github.com/zagvozdeen/ola/internal/config"
)

type Logger struct {
	log  *slog.Logger
	file *os.File
	once sync.Once
}

func New(cfg *config.Config) *Logger {
	if cfg.IsProduction {
		file, err := os.OpenFile("ola.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			slog.Error("Failed to open log file", slog.Any("error", err))
			os.Exit(1)
		}
		return &Logger{
			log: slog.New(slog.NewTextHandler(file, &slog.HandlerOptions{
				AddSource:   false,
				Level:       slog.LevelDebug,
				ReplaceAttr: nil,
			})),
			file: file,
		}
	}
	return &Logger{
		log: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource:   false,
			Level:       slog.LevelDebug,
			ReplaceAttr: nil,
		})),
	}
}

func (l *Logger) Close() {
	l.once.Do(func() {
		if l.file != nil {
			if err := l.file.Close(); err != nil {
				slog.Error("Failed to close log file", slog.Any("error", err))
			}
		}
	})
}

func (l *Logger) GetLog() *log.Logger {
	return slog.NewLogLogger(l.log.Handler(), slog.LevelDebug)
}

func (l *Logger) Debug(msg string, args ...slog.Attr) {
	l.log.Debug(msg, toAnySlice(args)...)
}

func (l *Logger) Debugf(format string, a ...any) {
	l.log.Debug(fmt.Sprintf(format, a...))
}

func (l *Logger) Info(msg string, args ...slog.Attr) {
	l.log.Info(msg, toAnySlice(args)...)
}

func (l *Logger) Infof(format string, a ...any) {
	l.log.Info(fmt.Sprintf(format, a...))
}

func (l *Logger) Warn(msg string, args ...slog.Attr) {
	l.log.Warn(msg, toAnySlice(args)...)
}

func (l *Logger) Warnf(format string, a ...any) {
	l.log.Warn(fmt.Sprintf(format, a...))
}

func (l *Logger) Error(msg string, err error, args ...slog.Attr) {
	args = append(args, slog.Any("error", err))
	l.log.Error(msg, toAnySlice(args)...)
}

func (l *Logger) Errorf(format string, a ...any) {
	l.log.Error(fmt.Sprintf(format, a...))
}

func (l *Logger) Fatalf(format string, a ...any) {
	l.Errorf(format, a...)
}

func (l *Logger) Printf(format string, a ...any) {
	l.Infof(format, a...)
}

func toAnySlice[T any](s []T) []any {
	r := make([]any, 0, len(s))
	for _, e := range s {
		r = append(r, e)
	}
	return r
}
