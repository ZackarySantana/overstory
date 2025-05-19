package main

import (
	"context"
	"fmt"

	"github.com/zackarysantana/overstory/cmd/internal"
	"github.com/zackarysantana/overstory/src/entities"
	"github.com/zackarysantana/overstory/src/service"
)

func main() {
	ctx := context.Background()

	mongoClient, cleanup, err := internal.UseMongoDBContainer(ctx)
	if err != nil {
		panic(err)
	}
	defer cleanup()

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

	newProject := &entities.Project{
		Name: "Test Project",
	}

	fmt.Println(newOrganization)

	fmt.Println(newUser)

	if err := s.CreateProject(ctx, newUser, newProject); err != nil {
		panic(fmt.Errorf("failed to create project: %w", err))
	}
	fmt.Println("1st project", newProject)

	if err := s.CreateProject(ctx, newUser, newProject); err != nil {
		panic(fmt.Errorf("failed to create project: %w", err))
	}
	fmt.Println("2nd project", newProject)
}
