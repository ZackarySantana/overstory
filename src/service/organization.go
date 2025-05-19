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
