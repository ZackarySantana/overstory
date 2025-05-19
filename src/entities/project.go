package entities

import "go.mongodb.org/mongo-driver/v2/bson"

type Project struct {
	ID bson.ObjectID `bson:"_id"`

	Name string
}
