package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id"`
	FirstName    string             `json:"first_name"`
	LastName     string             `json:"last_name"`
	Email        string             `json:"email"`
	PasswordHash string             `json:"-"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
}
