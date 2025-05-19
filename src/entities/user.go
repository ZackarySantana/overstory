package entities

import "go.mongodb.org/mongo-driver/v2/bson"

type Role struct {
	Role string
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
	Role

	CreateProject bool
}

func (r *OrganizationRole) CanCreateProject() bool {
	if r == nil {
		return false
	}

	return r.IsAdmin()
}

type ProjectRole struct {
	Role
}

type User struct {
	ID bson.ObjectID `bson:"_id"`

	OrganizationID   bson.ObjectID
	OrganizationRole OrganizationRole

	ProjectToRole map[bson.ObjectID]ProjectRole

	Username string
}
