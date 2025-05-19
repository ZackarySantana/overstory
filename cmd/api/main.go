package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/zackarysantana/overstory/cmd/internal"
	"github.com/zackarysantana/overstory/src/entities"
	"github.com/zackarysantana/overstory/src/service"
	"github.com/zackarysantana/overstory/src/slogctx"
)

func main() {
	ctx := context.Background()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ctx = slogctx.WithLogger(ctx, logger)

	mongoClient, err := internal.CreateMongoClient(ctx)
	if err != nil {
		logger.ErrorContext(ctx, "failed to create MongoDB client", "error", err)
		os.Exit(1)
	}

	s := service.New(mongoClient.Database("test"))
	if err := s.EnsureIndexes(ctx); err != nil {
		logger.ErrorContext(ctx, "failed to ensure indexes", "error", err)
		os.Exit(1)
	}

	newOrganization := &entities.Organization{
		Name: "Test Organization",
	}

	newUser := &entities.User{
		Username: "testuser",
	}

	if err := s.CreateOrganizationAndUser(ctx, newOrganization, newUser); err != nil {
		logger.ErrorContext(ctx, "failed to create organization and user", "error", err)
		os.Exit(1)
	}

	newProject := &entities.Project{
		Name: "Test Project",
	}

	logger.InfoContext(ctx, "created organization", "organization", newOrganization)

	logger.InfoContext(ctx, "created user", "user", newUser)

	if err := s.CreateProject(ctx, newUser, newProject); err != nil {
		logger.ErrorContext(ctx, "failed to create project", "error", err)
		os.Exit(1)
	}
	logger.InfoContext(ctx, "created project 1st", "project", newProject)

	if err := s.CreateProject(ctx, newUser, newProject); err != nil {
		logger.ErrorContext(ctx, "failed to create project", "error", err)
		os.Exit(1)
	}
	logger.InfoContext(ctx, "created project 2nd", "project", newProject)
}
