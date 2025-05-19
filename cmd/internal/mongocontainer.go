package internal

import (
	"context"
	"os"

	"github.com/zackarysantana/overstory/internal/testcontainer/mongodb"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	MongoDBContainerImage = "mongo:8.0.9"
)

func CreateMongoClient(ctx context.Context) (*mongo.Client, error) {
	if os.Getenv("MONGODB_URI") != "" {
		return mongo.Connect(options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	}

	client, _, err := mongodb.CreateContainer(ctx)
	return client, err
}
