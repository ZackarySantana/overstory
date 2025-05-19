package service

import (
	"errors"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
)

type Service struct {
	db *mongo.Database
}

func New(db *mongo.Database) *Service {
	return &Service{
		db: db,
	}
}
