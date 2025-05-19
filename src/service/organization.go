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
	if !actor.OrganizationRole.CanCreateProject() {
		return ErrUnauthorized
	}

	session, err := s.db.Client().StartSession()
	if err != nil {
		return fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	newProject.ID = bson.ObjectID{}
	newProject.OrganizationID = actor.OrganizationID
	_, err = session.WithTransaction(ctx, func(ctx context.Context) (interface{}, error) {
		inserted, err := s.db.Collection("projects").InsertOne(ctx, newProject)
		if err != nil {
			return nil, fmt.Errorf("database error: %w", err)
		}
		newProject.ID = inserted.InsertedID.(bson.ObjectID)

		if _, err := s.db.Collection("users").UpdateByID(ctx, actor.ID,
			bson.M{"$set": bson.M{
				fmt.Sprintf("projecttorole.%s", newProject.ID.Hex()): entities.ProjectRole{
					Role: entities.Role{Role: "admin"},
				},
			}},
		); err != nil {
			return nil, fmt.Errorf("failed to update user with new project: %w", err)
		}
		actor.SetProjectRole(newProject.ID, entities.ProjectRole{
			Role: entities.Role{Role: "admin"},
		})

		return nil, nil
	})

	return err
}

// CreateOrganizationAndUser creates a new organization and a new user
// belonging to that organization. The user is assigned the "admin" role
// in the organization.
// The organization's ID, the user's OrganizationID, and the user's ID
// are set within this function.
func (s *Service) CreateOrganizationAndUser(ctx context.Context, newOrg *entities.Organization, newUser *entities.User) error {
	session, err := s.db.Client().StartSession()
	if err != nil {
		return fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	newOrg.ID = bson.ObjectID{}
	newUser.ID = bson.ObjectID{}
	_, err = session.WithTransaction(ctx, func(ctx context.Context) (interface{}, error) {
		insertedOrg, err := s.db.Collection("organizations").InsertOne(ctx, newOrg)
		if err != nil {
			return nil, fmt.Errorf("failed to create organization: %w", err)
		}
		newOrg.ID = insertedOrg.InsertedID.(bson.ObjectID)

		newUser.OrganizationID = insertedOrg.InsertedID.(bson.ObjectID)
		newUser.OrganizationRole = entities.OrganizationRole{
			Role: entities.Role{Role: "admin"},
		}

		insertedUser, err := s.db.Collection("users").InsertOne(ctx, newUser)
		if err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
		newUser.ID = insertedUser.InsertedID.(bson.ObjectID)

		return nil, nil
	})

	if err != nil {
		newOrg.ID = bson.ObjectID{}
		newUser.OrganizationID = bson.ObjectID{}
		newUser.OrganizationRole = entities.OrganizationRole{}
		newUser.ID = bson.ObjectID{}
	}

	return err
}
