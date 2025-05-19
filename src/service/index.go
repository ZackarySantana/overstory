package service

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (s *Service) EnsureIndexes(ctx context.Context) error {
	errs := []error{
		s.ensureIndexesForProject(ctx),
	}

	return errors.Join(errs...)
}

func (s *Service) ensureIndexesForProject(ctx context.Context) error {
	index, err := s.db.Collection("projects").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "organization_id", Value: 1},
			{Key: "name", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})

	fmt.Println(index)

	return err
}
