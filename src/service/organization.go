package service

import (
	"context"
	"fmt"

	"github.com/zackarysantana/overstory/src/entities"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// CreateProject creates a new project belonging to the actor's organization.
// The actor must have permission to create projects in the organization.
// The project must not already exist in the organization.
// The project's ID and OrganizationID are set within this function.
func (s *Service) CreateProject(ctx context.Context, actor *entities.User, newProject *entities.Project) error {
	if actor.OrganizationRole.CanCreateProject() {
		return ErrUnauthorized
	}

	newProject.OrganizationID = actor.OrganizationID
	inserted, err := s.db.Collection("projects").InsertOne(ctx, newProject)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	newProject.ID = inserted.InsertedID.(bson.ObjectID)

	return nil
}

func (s *Service) CreateOrganizationAndUser(ctx context.Context, newOrg *entities.Organization, newUser *entities.User) error {
	session, err := s.db.Client().StartSession()
	if err != nil {
		return fmt.Errorf("failed to start session: %w", err)
	}

	_, err = session.WithTransaction(ctx, func(ctx context.Context) (interface{}, error) {
		insertedOrg, err := s.db.Collection("organizations").InsertOne(ctx, newOrg)
		if err != nil {
			return nil, fmt.Errorf("failed to create organization: %w", err)
		}

		newUser.OrganizationID = insertedOrg.InsertedID.(bson.ObjectID)
		newUser.OrganizationRole = entities.OrganizationRole{
			Role: entities.Role{Role: "admin"},
		}
		if _, err := s.db.Collection("users").InsertOne(ctx, newUser); err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}

		return nil, nil
	})

	if err != nil {
		newUser.OrganizationID = bson.ObjectID{}
		newUser.OrganizationRole = entities.OrganizationRole{}
	}

	return err
}
