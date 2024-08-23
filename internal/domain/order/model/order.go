package models

import (
	ticketModel "ginapp/internal/domain/ticket/model" // Import the ticket models package
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	Order_ID       primitive.ObjectID       `json:"_id" bson:"_id"`
	Order_Cart     []ticketModel.TicketUser `bson:"order_cart" json:"order_cart" `
	Orderered_At   time.Time                `bson:"ordered_at" json:"ordered_at"`
	Ticket_Count   int                      `bson:"ticket_count" json:"ticket_count"`
	Price          int                      `bson:"price" json:"price"`
	Discount       float64                  `bson:"discount" json:"discount"`
	Payment_Method Payement                 `bson:"payment_method" json:"payment_method"`
}

type Payement struct {
	Digital bool `bson:"digital"`
	COD     bool `bson:"cod"`
}
