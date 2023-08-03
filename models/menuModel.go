package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Menu struct {
	ID  primitive.ObjectID     `bson:"_id"`
	Name string                `json:"name" validate:"required,min=2,max=100"`
	Category string            `json:"category"`
	Start_Date *time.Time      `json:"start_date"`
	End_Date *time.Time        `json:"end_date"`
	CreatedAt time.Time        `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time        `bson:"updated_at" json:"updated_at"`
	Menu_id string             `json:"menu_id"`
}