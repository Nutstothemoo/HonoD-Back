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
	Username *string                           `json:"username" validate:"required,min=2,max=100"`
	Password *string                           `json:"password"`
	Email *string                              `json:"email" validate:"required"`
	Phone *string                              `json:"phone" validate:"required"`
	Role *string                               `json:"role"`
	Token *string                              `json:"token"`
	Refresh_Token *string                      `json:"refresh_token"`
	Created_At time.Time                       `json:"created_at"`
	Updated_At time.Time                       `json:"updated_at"`
	User_ID string                             `json:"user_id"`
	UserCart []ProductUser                     `json:"usercart" bson:"usercart"`
	Address_Details []Address                  `json:"adress" bson:"adress"`
	Order_History []Order                      `json:"order_history" bson:"order_history"`
	Order_In_Progress []Order                  `json:"order_in_progress" bson:"order_in_progress"`
	Order_In_Delivery []Order                  `json:"order_in_delivery" bson:"order_in_delivery"`
	Order_Delivered []Order                    `json:"order_delivered" bson:"order_delivered"`
	Order_Canceled []Order                     `json:"order_canceled" bson:"order_canceled"`
	Order_Refunded []Order                     `json:"order_refunded" bson:"order_refunded"`
	Order_Returned []Order                     `json:"order_returned" bson:"order_returned"`
	Order_Awaiting_Payment []Order             `json:"order_awaiting_payment" bson:"order_awaiting_payment"`
	Order_Status []Order                       `json:"orders" bson:"orders"`
	}

type Product struct {
	ID primitive.ObjectID                      `bson:"_id" json:"_id"`
	Product_Name *string                       `json:"product_name" bson:"product_name"`
  Price *uint64                              `json:"price" bson:"price"`
	Rating *uint8                              `json:"rating" bson:"rating"`
	Image *string                              `json:"image" bson:"image"`
}	

type ProductUser struct {
	ID primitive.ObjectID                    `json:"_id" bson:"_id"`
	Product_Name *string                     `json:"product_name" bson:"product_name"`
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
	Order_Cart []ProductUser                     `bson:"order_cart" json:"order_cart" `
	Orderered_At time.Time                       `bson:"ordered_at" json:"ordered_at"`
	Price int                               `bson:"price" json:"price"`
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

// type User struct {
// 	ID       string   `json:"id"`
// 	FirstName string `json:"firstName"`
// 	LastName  string `json:"lastName"`
// 	Avatar    string `json:"avatar"`
// 	Email     string `json:"email"`
// 	Contact   Contact `json:"contact"`
// 	FacebookID string `json:"facebookId"`
// 	MetaUserData MetaUserData `json:"metaUserData"`
// }