package entities

import "go.mongodb.org/mongo-driver/v2/bson"

type Role string

const (
	RoleAdmin  Role = "admin"
	RoleMember Role = "member"
	RoleGuest  Role = "guest"
)

type User struct {
	ID bson.ObjectID `bson:"_id"`

	OrganizationID   bson.ObjectID
	OrganizationRole Role

	ProjectToRole map[bson.ObjectID]Role
}
