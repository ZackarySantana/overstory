package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/zackarysantana/overstory/src/logging/slogctx"
	"github.com/zackarysantana/overstory/src/logging/slogerr"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (s *Service) EnsureIndexes(ctx context.Context) error {
	errs := []error{
		s.ensureIndexesForProject(ctx),
	}

	err := errors.Join(errs...)

	slogctx.Logger(ctx).ErrorContext(ctx, "failed to ensure indexes", slogerr.ErrorKey, err)

	return err
}

func (s *Service) ensureIndexesForProject(ctx context.Context) error {
	slogctx.Logger(ctx).DebugContext(ctx, "ensuring indexes for projects")
	index, err := s.db.Collection("projects").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "organization_id", Value: 1},
			{Key: "name", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("failed to create index for projects: %w", err)
	}

	slogctx.Logger(ctx).DebugContext(ctx, "created index for projects", "index", index)

	return err
}
