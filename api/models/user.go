package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID primitive.ObjectID                      `json:"_id" bson:"_id"`
	First_Name *string                         `json:"first_name" validate:"required,min=2,max=100"`
	Last_Name *string                          `json:"last_name" validate:"required,min=2,max=100"`
	Avatar    string                           `json:"avatar"`
	FacebookID *string                         `json:"facebookId"`
	GoogleID *string													 `json:"googleId"`
	User_ID string                             `json:"user_id"`
	Username *string                           `json:"username" validate:"required,min=2,max=100"`
	Password *string                           `json:"password"`
	Email *string                              `json:"email" validate:"required"`
	Phone *string                              `json:"phone" validate:"required"`
	Role *string                               `json:"role"`
	Token *string                              `json:"token"`
	Refresh_Token *string                      `json:"refresh_token"`
	Created_At time.Time                       `json:"created_at"`
	Updated_At time.Time                       `json:"updated_at"`
	UserCart []TicketUser                     `json:"usercart" bson:"usercart"`
	Address_Details []Address                  `json:"adress" bson:"adress"`
	Order_History []Order                      `json:"order_history" bson:"order_history"`
	Order_Canceled []Order                     `json:"order_canceled" bson:"order_canceled"`
	Order_Refunded []Order                     `json:"order_refunded" bson:"order_refunded"`
	}

	type Ticket struct {
    ID        primitive.ObjectID `bson:"_id" json:"_id"`
    TicketName *string           `json:"ticket_name" bson:"ticket_name"`
    Price      *uint64           `json:"price" bson:"price"`
    Rating     *uint8            `json:"rating" bson:"rating"`
    Image      *string           `json:"image" bson:"image"`
    EventID    primitive.ObjectID `bson:"event_id" json:"event_id"`
    DealerID   primitive.ObjectID `bson:"dealer_id" json:"dealer_id"`
    CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
    UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}

type TicketUser struct {
	ID primitive.ObjectID                    `json:"_id" bson:"_id"`
	Ticket_name *string                      `json:"ticket_name" bson:"ticket_name"`
	Price *uint64                            `json:"price" bson:"price"`
	Rating *uint8                            `json:"rating" bson:"rating"`
	Image *string                            `json:"image" bson:"image"`
}

type Address struct {
	ID primitive.ObjectID                      `json:"_id" bson:"_id"`
	House *string                              `json:"house" bson:"house"`
	Street *string                             `json:"street" bson:"street"`
	City *string                               `json:"city" bson:"city"`
	Zipcode *uint                              `json:"zipcode" bson:"zipcode"`
	State *string                              `json:"state" bson:"state"`
	Country *string                            `json:"country" bson:"country"`
}

type Order struct {
	Order_ID primitive.ObjectID                  `json:"_id" bson:"_id"`
	Order_Cart []TicketUser                      `bson:"order_cart" json:"order_cart" `
	Orderered_At time.Time                       `bson:"ordered_at" json:"ordered_at"`
	Ticket_Count int                             `bson:"ticket_count" json:"ticket_count"`
	Price int                                    `bson:"price" json:"price"`
	Discount	float64                            `bson:"discount" json:"discount"`
	Payment_Method Payement                      `bson:"payment_method" json:"payment_method"`
}

type Payement struct {
	Digital bool                                  `bson:"digital"`
	COD bool                                      `bson:"cod"`
}

type Contact struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone     string `json:"phone"`
	Birthday  string `json:"birthday"`
	ZipCode   string `json:"zipCode"`
	FederalTaxNumber string `json:"federalTaxNumber"`
}

type MetaUserData struct {
	ExternalID string `json:"external_id"`
	Em         string `json:"em"`
	Ph         string `json:"ph"`
	Fn         string `json:"fn"`
	Ln         string `json:"ln"`
}

