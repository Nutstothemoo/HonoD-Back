package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderItem struct {
	ID         primitive.ObjectID `bson:"_id"`
	OrderID    primitive.ObjectID `json:"order_id"`
	ProductID  primitive.ObjectID `json:"product_id"`
	Quantity   int                `json:"quantity"`
	UnitPrice  float64            `json:"unit_price"`
	TotalPrice float64            `json:"total_price"`
	CreatedAt  time.Time          `json:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at"`
}
