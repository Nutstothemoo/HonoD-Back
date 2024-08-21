package helpers

import (
	"context"
	"errors"
	"ginapp/internal/interfaces/api/models"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrorCantFindTicket         = errors.New("can't find ticket")
	ErrorCantDecodeTicket       = errors.New("can't decode ticket")
	ErrorUserIdisNotValid       = errors.New("this userId is not valid")
	ErrorCantFindUser           = errors.New("can't find user, this userId is not valid")
	ErrorCantUpdateUser         = errors.New("can't update user")
	ErrorCantGetItemFromCart    = errors.New("can't get item from cart")
	ErrorCantAddItemToCart      = errors.New("can't add item to cart")
	ErrorCantRemoveItemFromCart = errors.New("can't remove item from cart")
)

func AddTicketToCart(ctx context.Context, ticketCollection *mongo.Collection, userCollection *mongo.Collection, ticketID primitive.ObjectID, userID string) error {

	searchFromdb, err := ticketCollection.Find(ctx, bson.M{"_id": ticketID})
	if err != nil {
		log.Fatal(err)
		return ErrorCantFindTicket
	}
	var ticketCart []models.Ticket
	err = searchFromdb.All(ctx, &ticketCart)
	if err != nil {
		log.Fatal(err)
		return ErrorCantDecodeTicket
	}
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Fatal(err)
		return ErrorUserIdisNotValid
	}
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "usercart", Value: bson.D{{Key: "$each", Value: ticketCart}}}}}}
	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return ErrorCantUpdateUser
	}
	return nil
}

func RemoveTicketfromCart(ctx context.Context, ticketCollection *mongo.Collection, userCollection *mongo.Collection, ticketID primitive.ObjectID, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrorUserIdisNotValid
	}
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.M{"$pull": bson.M{"usercart": bson.M{"_id": ticketID}}}
	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		return ErrorCantRemoveItemFromCart
	}
	return nil

}
func BuyTicketFromCart(ctx context.Context, userCollection *mongo.Collection, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrorUserIdisNotValid
	}
	var getcartitems models.User
	var ordercart models.Order
	ordercart.Order_ID = primitive.NewObjectID()
	ordercart.Orderered_At = time.Now()
	ordercart.Order_Cart = make([]models.TicketUser, 0)
	ordercart.Payment_Method.COD = true
	unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$usercart"}}}}
	grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$usercart.price"}}}}}}
	currentresults, err := userCollection.Aggregate(ctx, mongo.Pipeline{unwind, grouping})
	ctx.Done()
	if err != nil {
		panic(err)
	}
	var getusercart []bson.M
	if err = currentresults.All(ctx, &getusercart); err != nil {
		panic(err)
	}
	var total_price int32
	for _, user_item := range getusercart {
		price := user_item["total"]
		total_price = price.(int32)
	}
	ordercart.Price = int(total_price)
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: ordercart}}}}
	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Println(err)
	}
	err = userCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: id}}).Decode(&getcartitems)
	if err != nil {
		log.Println(err)
	}
	filter2 := bson.D{primitive.E{Key: "_id", Value: id}}
	update2 := bson.M{"$push": bson.M{"orders.$[].order_list": bson.M{"$each": getcartitems.UserCart}}}
	_, err = userCollection.UpdateOne(ctx, filter2, update2)
	if err != nil {
		return ErrorCantUpdateUser
	}
	usercart_empty := make([]models.TicketUser, 0)
	filtered := bson.D{primitive.E{Key: "_id", Value: id}}
	updated := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "usercart", Value: usercart_empty}}}}
	_, err = userCollection.UpdateOne(ctx, filtered, updated)
	if err != nil {
		return ErrorCantUpdateUser

	}
	return nil
}

func InstantBuyer(ctx context.Context, ticketCollection *mongo.Collection, userCollection *mongo.Collection, ticketID primitive.ObjectID, UserID string) error {
	id, err := primitive.ObjectIDFromHex(UserID)
	if err != nil {
		log.Println(err)
		return ErrorUserIdisNotValid
	}
	var product_details models.TicketUser
	var orders_detail models.Order
	orders_detail.Order_ID = primitive.NewObjectID()
	orders_detail.Orderered_At = time.Now()
	orders_detail.Order_Cart = make([]models.TicketUser, 0)
	orders_detail.Payment_Method.COD = true
	err = ticketCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: ticketID}}).Decode(&product_details)
	if err != nil {
		log.Println(err)
	}
	if product_details.Price != nil {
		orders_detail.Price = int(*product_details.Price)
	} else {
		orders_detail.Price = 0
	}
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: orders_detail}}}}
	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println(err)
	}
	return nil
}
