
package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type Ticket struct {
	TicketName *string           	`json:"ticket_name" bson:"ticket_name"`
	Price      *uint64           	`json:"price" bson:"price"`
	Rating     *uint8           	`json:"rating" bson:"rating"`
	Image      *string           	`json:"image" bson:"image"`
	Stock  			*uint64           `json:"stock" bson:"stock"`
	EventID    primitive.ObjectID `bson:"event_id" json:"event_id"`
	DealerID   primitive.ObjectID `bson:"dealer_id" json:"dealer_id"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}

type TicketUnit struct {
	EventID   primitive.ObjectID 		`bson:"event_id" json:"event_id"`
	OwnerID   primitive.ObjectID 		`bson:"owner_id" json:"owner_id"`
	Price      *uint64           		`json:"price" bson:"price"`
	IsScanned   *bool             		`json:"is_scanned" bson:"is_scanned"`
}

type TicketUser struct {
Ticket_name *string                      `json:"ticket_name" bson:"ticket_name"`
Price *uint64                            `json:"price" bson:"price"`
Rating *uint8                            `json:"rating" bson:"rating"`
Image *string                            `json:"image" bson:"image"`
}

