package controllers

import (
	"context"
	"errors"
	"fmt"
	"ginapp/internal/domain/qrcode"
	"ginapp/internal/infrastructure/s3"
	"ginapp/internal/interfaces/api/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	TicketCollection     *mongo.Collection
	UserCollection       *mongo.Collection
	EventCollection      *mongo.Collection
	TicketUnitCollection *mongo.Collection
}

func NewApplication(ticketCollection, userCollection, eventCollection, ticketUnitCollection *mongo.Collection) *Application {
	return &Application{
		TicketCollection:     ticketCollection,
		UserCollection:       userCollection,
		EventCollection:      eventCollection,
		TicketUnitCollection: ticketUnitCollection,
	}
}

func (app *Application) InstantBuy() gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		s3BucketName := os.Getenv("S3_BUCKET_NAME")
		s3Client, err := s3.NewS3Client()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create S3 client"})
		}

		userID, ok := c.Get("userId")
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to assert userID"})
			return
		}
		ticketID := c.Param("ticketId")
		eventID := c.Param("eventId")

		if userID == "" || ticketID == "" || eventID == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "userId, productId, and eventId are required"})
			return
		}
		qrCodeFileName := fmt.Sprintf("%s/%s/qrcode_%s_%s.png", eventID, ticketID, userID, time.Now().Format("2006-01-02 15:04:05"))
		price, err := updateTicketCount(ctx, app.TicketCollection, app.UserCollection, ticketID, userID.(string), qrCodeFileName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		qrCodeData := "ProductID: " + ticketID
		bucketBasics := s3.BucketBasics{S3Client: s3Client}
		err = qrcode.GenerateAndUploadQRCode(ctx, qrCodeData, s3BucketName, qrCodeFileName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate and upload QR code"})
			return
		}

		err = app.updateUserPurchases(ctx, userID.(string), eventID, price, qrCodeFileName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Successfully placed the order"})
	}
}

func (app *Application) updateUserPurchases(ctx context.Context, userID string, eventID string, price uint64, qrCodeFileName string) error {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	eventObjectID, err := primitive.ObjectIDFromHex(eventID)
	if err != nil {
		return err
	}

	isScanned := false

	ticketUnit := models.TicketUnit{
		EventID:   eventObjectID,
		OwnerID:   userObjectID,
		Price:     &price,
		IsScanned: &isScanned,
	}

	result, err := app.TicketUnitCollection.InsertOne(ctx, ticketUnit)
	if err != nil {
		return err
	}

	ticketUnitID := result.InsertedID
	update := bson.M{
		"$push": bson.M{
			"purchases": ticketUnitID,
		},
	}

	_, err = app.UserCollection.UpdateOne(ctx, bson.M{"_id": userID}, update)
	return err
}

func updateTicketCount(ctx context.Context, ticketCollection *mongo.Collection, userCollection *mongo.Collection, ticketID string, userID string, qrCodeFileName string) (uint64, error) {
	var ticket models.Ticket
	ticketObjectID, err := primitive.ObjectIDFromHex(ticketID)
	if err != nil {
		return 0, err
	}
	err = ticketCollection.FindOne(ctx, bson.M{"_id": ticketObjectID}).Decode(&ticket)
	if err != nil {
		return 0, err
	}
	if *ticket.Stock == 0 {
		return 0, errors.New("ticket is out of stock")
	}
	// PREUVE D ACHAT ICI

	update := bson.M{
		"$inc": bson.M{
			"stock": -1,
		},
	}
	_, err = ticketCollection.UpdateOne(ctx, bson.M{"_id": ticketObjectID}, update)
	if err != nil {
		return 0, err
	}
	return *ticket.Price, nil
}

// func (app *Application) AddToCart() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 			ticketQueryID := c.Query("id")
// 			if ticketQueryID == "" {
// 					log.Println("ticket id is empty")
// 					_ = c.AbortWithError(http.StatusBadRequest, errors.New("ticket id is empty"))
// 					return
// 			}
// 			userQueryID := c.MustGet("userId").(string)

// 			ticketID, err := primitive.ObjectIDFromHex(ticketQueryID)
// 			if err != nil {
// 					log.Println(err)
// 					c.AbortWithStatus(http.StatusInternalServerError)
// 					return
// 			}
// 			var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
// 			defer cancel()

// 			err = AddTicketToCart(ctx, app.TicketCollection, app.UserCollection, ticketID, userQueryID)
// 			if err != nil {
// 					c.IndentedJSON(http.StatusInternalServerError, err)
// 			}
// 			c.IndentedJSON(201, "Successfully Added to the cart")
// 	}
// }

// func (app *Application) RemoveItem() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 			ticketQueryID := c.Query("id")

// 			if ticketQueryID == "" {
// 					log.Println("ticket id is invalid")
// 					_ = c.AbortWithError(http.StatusBadRequest, errors.New("ticket id is empty"))
// 					return
// 			}

// 			userQueryID := c.MustGet("userId").(string)
// 			if userQueryID == "" {
// 					log.Println("user id is empty")
// 					_ = c.AbortWithError(http.StatusBadRequest, errors.New("UserID is empty"))
// 			}

// 			ticketID, err := primitive.ObjectIDFromHex(ticketQueryID)
// 			if err != nil {
// 					log.Println(err)
// 					c.AbortWithStatus(http.StatusInternalServerError)
// 					return
// 			}

// 			var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
// 			defer cancel()
// 			err = RemoveTicketFromCart(ctx, app.UserCollection, ticketID, userQueryID)
// 			if err != nil {
// 					c.IndentedJSON(http.StatusInternalServerError, err)
// 					return
// 			}
// 			c.IndentedJSON(200, "Successfully removed from cart")
// 	}
// }

// func GetItemFromCart() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		user_id := c.Query("id")
// 		if user_id == "" {
// 			c.Header("Content-Type", "application/json")
// 			c.JSON(http.StatusNotFound, gin.H{"error": "invalid id"})
// 			c.Abort()
// 			return
// 		}

// 		usert_id, _ := primitive.ObjectIDFromHex(user_id)

// 		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
// 		defer cancel()

// 		var filledcart models.User
// 		err := UserCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: usert_id}}).Decode(&filledcart)
// 		if err != nil {
// 			log.Println(err)
// 			c.IndentedJSON(500, "not id found")
// 			return
// 		}

// 		filter_match := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: usert_id}}}}
// 		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$usercart"}}}}
// 		grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$usercart.price"}}}}}}
// 		pointcursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{filter_match, unwind, grouping})
// 		if err != nil {
// 			log.Println(err)
// 		}
// 		var listing []bson.M
// 		if err = pointcursor.All(ctx, &listing); err != nil {
// 			log.Println(err)
// 			c.AbortWithStatus(http.StatusInternalServerError)
// 		}
// 		for _, json := range listing {
// 			c.IndentedJSON(200, json["total"])
// 			c.IndentedJSON(200, filledcart.UserCart)
// 		}
// 		ctx.Done()
// 	}
// }

// func (app *Application) BuyFromCart() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		userQueryID := c.Query("id")
// 		productIDStr := c.Query("productID")
// 		productID, err := primitive.ObjectIDFromHex(productIDStr)
// 		if err != nil {
// 			log.Println(err)
// 			c.AbortWithStatus(http.StatusInternalServerError)
// 			return
// 		}
// 		if userQueryID == "" {
// 			log.Panicln("user id is empty")
// 			_ = c.AbortWithError(http.StatusBadRequest, errors.New("UserID is empty"))
// 		}
// 		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
// 		defer cancel()

// 		err = BuyProductFromCart(ctx, app.TicketCollection, productID, userQueryID)
// 		if err != nil {
// 			c.IndentedJSON(http.StatusInternalServerError, err)
// 		}
// 		c.IndentedJSON(200, "Successfully Placed the order")
// 	}
// }

// func RemoveProductofCart(ctx context.Context, userCollection *mongo.Collection, ticketID  primitive.ObjectID, userID string) error {

// 	filter := bson.D{{"_id", userID}}
// 	update := bson.D{{"$pull", bson.D{{"usercart", bson.D{{"_id", ticketID}}}}}}
// 	_, err := userCollection.UpdateOne(ctx, filter, update)
// 	return err
// }

// BuyProductFromCart purchases a product from the user's cart in the database.
// func BuyProductFromCart(ctx context.Context,  userCollection *mongo.Collection, ticketID  primitive.ObjectID, userID string) error {
// 	// This function depends on your application logic.
// 	// You might need to remove the product from the user's cart and decrease the product's stock.
// 	// Here is a simple example:
// 	err := RemoveProductofCart(ctx, userCollection, ticketID, userID)
// 	if err != nil {
// 			return err
// 	}
// 	// You might need to update the product's stock in another collection.
// 	// You would need to write that code here.
// 	return nil
// }

// func RemoveTicketFromCart(ctx context.Context, userCollection *mongo.Collection, ticketID  primitive.ObjectID, userID string) error {
// 	filter := bson.D{{"_id", userID}}
// 	update := bson.D{{"$pull", bson.D{{"usercart", bson.D{{"_id", ticketID}}}}}}
// 	_, err := userCollection.UpdateOne(ctx, filter, update)
// 	return err
// }

// func AddTicketToCart(ctx context.Context, ticketCollection *mongo.Collection, userCollection *mongo.Collection, ticketID primitive.ObjectID, userID string) error {
// 	filter := bson.D{{Key: "email", Value: userID}}
// 	update := bson.D{{Key: "$push", Value: bson.D{{"usercart", bson.D{{"_id", ticketID}}}}}}
// 	_, err := userCollection.UpdateOne(ctx, filter, update)
// 	return err
// }
