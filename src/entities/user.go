package entities

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Role struct {
	Role string `bson:"role,omitempty"`
}

const (
	RoleAdmin = "admin"
)

func (r *Role) IsAdmin() bool {
	if r == nil {
		return false
	}

	return r.Role == "admin"
}

type OrganizationRole struct {
	Role `bson:"role,omitempty"`

	CreateProject bool `bson:"create_project,omitempty"`
}

func (r *OrganizationRole) CanCreateProject() bool {
	if r == nil {
		return false
	}

	return r.IsAdmin() || r.CreateProject
}

type ProjectRole struct {
	Role `bson:"role,omitempty"`
}

type User struct {
	ID bson.ObjectID `bson:"_id,omitempty"`

	OrganizationID   bson.ObjectID    `bson:"organization_id,omitempty"`
	OrganizationRole OrganizationRole `bson:"organization_role,omitempty"`

	ProjectToRole map[bson.ObjectID]ProjectRole `bson:"project_to_role,omitempty"`

	Username string `bson:"username,omitempty"`
}

func (u *User) SetProjectRole(projectID bson.ObjectID, role ProjectRole) {
	if u.ProjectToRole == nil {
		u.ProjectToRole = make(map[bson.ObjectID]ProjectRole)
	}
	u.ProjectToRole[projectID] = role
}
