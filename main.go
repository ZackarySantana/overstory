package main

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/zackarysantana/overstory/src/entities"
	"github.com/zackarysantana/overstory/src/service"
	"github.com/zackarysantana/overstory/src/urlquery"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	ctx := context.Background()

	mongodbContainer, err := mongodb.Run(ctx, "mongo:8.0.9", mongodb.WithReplicaSet("rs0"))
	if err != nil {
		panic(err)
	}
	defer mongodbContainer.Terminate(ctx)
	endpoint, err := mongodbContainer.ConnectionString(ctx)
	if err != nil {
		panic(err)
	}
	endpoint, err = urlquery.AddQueryParam(endpoint, "directConnection", "true")
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

	newOrganization := &entities.Organization{
		Name: "Test Organization",
	}

	newUser := &entities.User{
		Username: "testuser",
	}

	if err := s.CreateOrganizationAndUser(ctx, newOrganization, newUser); err != nil {
		panic(fmt.Errorf("failed to create organization and user: %w", err))
	}
}
