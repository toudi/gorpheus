package commands

import (
	"context"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"

	"github.com/toudi/gorpheus/v2/interfaces"
)

type Logger struct {
	Level  slog.Level
	logger *slog.Logger
	ctx    context.Context
}

func loggingLevel(verbosity int) slog.Level {
	if verbosity > 1 {
		return slog.LevelDebug
	}
	if verbosity > 0 {
		return slog.LevelWarn
	}
	return slog.LevelInfo
}

func NewLogger(ctx context.Context, verbosity int) interfaces.Logger {
	var logOutput = os.Stderr

	logger := &Logger{
		logger: slog.New(tint.NewHandler(
			logOutput, &tint.Options{
				Level: loggingLevel(verbosity),
			},
		)),
		ctx: ctx,
	}
	return logger
}

func (l Logger) Log(m interfaces.LogMessage) {
	l.logger.Log(l.ctx, m.Level, m.Message, m.Payload...)
}
