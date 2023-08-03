package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Table struct {
	ID          primitive.ObjectID `bson:"_id"`
	Number      int                `json:"number"`
	Capacity    int                `json:"capacity"`
	Description string             `json:"description"`
	Status      string             `json:"status" validate:"eq=AVAILABLE|eq=OCCUPIED|eq=RESERVED|eq=BLOCKED"`
}