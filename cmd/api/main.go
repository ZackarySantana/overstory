package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/zackarysantana/overstory/cmd/internal"
	"github.com/zackarysantana/overstory/src/entities"
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

	newOrganization := &entities.Organization{
		Name: "Test Organization",
	}

	newUser := &entities.User{
		Username: "testuser",
	}

	if err := s.CreateOrganizationAndUser(ctx, newOrganization, newUser); err != nil {
		logger.ErrorContext(ctx, "failed to create organization and user", slogerr.ErrorKey, err)
		os.Exit(1)
	}
	logger.InfoContext(ctx, "created organization", "organization", newOrganization)
	logger.InfoContext(ctx, "created user", "user", newUser)

	newProject := &entities.Project{
		Name: "Test Project",
	}

	if err := s.CreateProject(ctx, newUser, newProject); err != nil {
		logger.ErrorContext(ctx, "failed to create project", slogerr.ErrorKey, err)
		os.Exit(1)
	}
	logger.InfoContext(ctx, "created project 1st", "project", newProject)

	if err := s.CreateProject(ctx, newUser, newProject); err != nil {
		logger.ErrorContext(ctx, "failed to create project", slogerr.ErrorKey, err)
		os.Exit(1)
	}
	logger.InfoContext(ctx, "created project 2nd", "project", newProject)
}
