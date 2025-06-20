package main

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/zackarysantana/overstory/cmd/internal"
	"github.com/zackarysantana/overstory/src/api"
	"github.com/zackarysantana/overstory/src/logging/slogctx"
	"github.com/zackarysantana/overstory/src/logging/slogerr"
	"github.com/zackarysantana/overstory/src/service"
)

func main() {
	ctx := context.Background()

	var handler slog.Handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	handler = slogerr.NewErrorNilDropperHandler(handler)
	logger := slog.New(handler)
	ctx = slogctx.WithLogger(ctx, logger)

	mongoClient, err := internal.CreateMongoClient(ctx)
	if err != nil {
		logger.ErrorContext(ctx, "failed to create MongoDB client", slogerr.ErrorKey, err)
		os.Exit(1)
	}

	s := service.New(mongoClient.Database("test"))
	if err := s.EnsureIndexes(ctx); err != nil {
		logger.ErrorContext(ctx, "failed to ensure indexes", slogerr.ErrorKey, err)
		os.Exit(1)
	}

	server := &http.Server{
		Addr:    ":8080",
		Handler: api.New(ctx, s),
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
		ErrorLog: slog.NewLogLogger(handler, slog.LevelError),
	}

	if err := server.ListenAndServe(); err != nil {
		logger.ErrorContext(ctx, "failed to start server", slogerr.ErrorKey, err)
	}
}
