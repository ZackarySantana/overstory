package entities

import "go.mongodb.org/mongo-driver/v2/bson"

type Organization struct {
	ID bson.ObjectID `bson:"_id,omitempty"`

	Name string `bson:"name,omitempty"`
}
