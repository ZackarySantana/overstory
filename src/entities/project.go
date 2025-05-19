package entities

import "go.mongodb.org/mongo-driver/v2/bson"

type Project struct {
	ID bson.ObjectID `bson:"_id,omitempty"`

	OrganizationID bson.ObjectID `bson:"organization_id,omitempty"`

	Name string `bson:"name,omitempty"`
}
