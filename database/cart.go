package database

import (
	"context"
	"errors"
	"ginapp/models"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)


var (
	ErrorCantFindProduct = errors.New("Can't find product")
	ErrorCantDecodeProduct = errors.New("Can't decode product")
	UserIdisNotValid = errors.New("This userId is not valid")
	ErrorCantFindUser = errors.New("Can't find user, this userId is not valid")
	ErrorCantUpdateUser = errors.New("Can't update user")
	ErrorCantGetItemFromCart = errors.New("Can't get item from cart")
	ErrorCantAddItemToCart = errors.New("Can't add item to cart")
	ErrorCantRemoveItemFromCart = errors.New("Can't remove item from cart")
)

func AddProductToCart(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string ) error{

		searchFromdb, err := prodCollection.Find(ctx, bson.M{"_id": productID})
		if err != nil {
			log.Fatal(err)
			return 	ErrorCantFindProduct
		}
		var productCart []models.Product
		err = searchFromdb.All(ctx, &productCart)
		if err != nil {
			log.Fatal(err)
			return ErrorCantDecodeProduct
		}
		id, err:= primitive.ObjectIDFromHex(userID)
		if err != nil {
			log.Fatal(err)
			return UserIdisNotValid
		}
		filter := bson.D{primitive.E{Key: "_id", Value: id}}
		update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "usercart", Value: bson.D{{Key: "$each", Value: productCart}}}}}}
		_, err = userCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			return ErrorCantUpdateUser
		}
		return nil
	}
func RemoveItemFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}
func BuyItemFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}
func InstantBuy() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}
