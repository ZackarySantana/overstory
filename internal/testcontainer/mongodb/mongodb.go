package mongodb

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/zackarysantana/overstory/src/urlquery"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	MongoDBContainerImage = "mongo:8.0.9"
)

func CreateContainer(ctx context.Context) (*mongo.Client, func() error, error) {
	mongodbContainer, err := mongodb.Run(ctx, MongoDBContainerImage, mongodb.WithReplicaSet("rs0"))
	if err != nil {
		return nil, func() error {
			return nil
		}, fmt.Errorf("failed to start MongoDB container: %w", err)
	}
	cleanup := func() error {
		return mongodbContainer.Terminate(ctx)
	}
	endpoint, err := mongodbContainer.ConnectionString(ctx)
	if err != nil {
		return nil, cleanup, fmt.Errorf("failed to get MongoDB connection string: %w", err)
	}
	endpoint, err = urlquery.AddQueryParam(endpoint, "directConnection", "true")
	if err != nil {
		return nil, cleanup, fmt.Errorf("failed to add query param to MongoDB connection string: %w", err)
	}
	mongoClient, err := mongo.Connect(options.Client().ApplyURI(endpoint))
	if err != nil {
		return nil, cleanup, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}
	return mongoClient, cleanup, nil
}
