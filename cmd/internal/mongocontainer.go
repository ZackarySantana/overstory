package internal

import (
	"context"
	"os"

	"github.com/zackarysantana/overstory/internal/testcontainer/mongodb"
	"github.com/zackarysantana/overstory/src/logging/slogctx"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	MongoDBContainerImage = "mongo:8.0.9"
)

func CreateMongoClient(ctx context.Context) (*mongo.Client, error) {
	if os.Getenv("MONGODB_URI") != "" {
		slogctx.Logger(ctx).DebugContext(ctx, "using MongoDB URI from environment variable")
		return mongo.Connect(options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	}

	slogctx.Logger(ctx).DebugContext(ctx, "using MongoDB testcontainer")
	client, _, err := mongodb.CreateContainer(ctx)
	return client, err
}
