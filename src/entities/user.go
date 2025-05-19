package entities

type User struct {
	ID any `bson:"_id"`

	OrganizationID any `bson:"project_id"`
}
