package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID           primitive.ObjectID `bson:"_id"`
	OrderID      string             `json:"order_id"`
	UserID       primitive.ObjectID `json:"user_id"`
	TotalAmount  float64            `json:"total_amount"`
	OrderDate    time.Time          `json:"order_date"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
}