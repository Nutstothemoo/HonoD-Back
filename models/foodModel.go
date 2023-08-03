package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Food struct {
	ID  primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name *string `json:"name" validate:"required,min=2,max=100"`
	Price *float64 `json:"price" validate:"required"`
	Food_image *string `json:"food_image"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
	Food_id *string `json:"food_id"`
	Menu_id *string `json:"menu_id" validate:"required"`
}