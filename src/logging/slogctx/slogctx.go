package slogctx

import (
	"context"
	"log/slog"
)

type key struct{}

func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, key{}, logger)
}

func Logger(ctx context.Context) *slog.Logger {
	if ctx == nil {
		return slog.Default()
	}
	logger, ok := ctx.Value(key{}).(*slog.Logger)
	if !ok || logger == nil {
		return slog.Default()
	}
	return logger
}
