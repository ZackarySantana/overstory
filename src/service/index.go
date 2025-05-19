package service

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
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
			{Key: "OrganizationID", Value: 1},
			{Key: "Name", Value: 1},
		},
	})

	fmt.Println(index)

	return err
}
