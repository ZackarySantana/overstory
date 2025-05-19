package main

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/zackarysantana/overstory/src/service"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	ctx := context.Background()

	mongodbContainer, err := mongodb.Run(ctx, "mongo:8.0.9")
	if err != nil {
		panic(err)
	}
	defer mongodbContainer.Terminate(ctx)
	endpoint, err := mongodbContainer.ConnectionString(ctx)
	if err != nil {
		panic(err)
	}
	mongoClient, err := mongo.Connect(options.Client().ApplyURI(endpoint))
	if err != nil {
		panic(err)
	}

	s := service.New(mongoClient.Database("test"))
	if err := s.EnsureIndexes(ctx); err != nil {
		panic(err)
	}

	fmt.Println("Indexes ensured")
}
