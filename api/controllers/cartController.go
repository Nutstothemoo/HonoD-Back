package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"time"
	"ginapp/api/models"
	"ginapp/api/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)



type Application struct {
	TicketCollection *mongo.Collection
	UserCollection   *mongo.Collection
	EventCollection  *mongo.Collection
}

func NewApplication(ticketCollection, userCollection, eventCollection *mongo.Collection) *Application {
	return &Application{
			TicketCollection: ticketCollection,
			UserCollection:   userCollection,
			EventCollection:  eventCollection,
	}
}


func (app *Application) AddToCart() gin.HandlerFunc {
	return func(c *gin.Context) {
			ticketQueryID := c.Query("id")
			if ticketQueryID == "" {
					log.Println("ticket id is empty")
					_ = c.AbortWithError(http.StatusBadRequest, errors.New("ticket id is empty"))
					return
			}
			userQueryID := c.MustGet("userId").(string)

			ticketID, err := primitive.ObjectIDFromHex(ticketQueryID)
			if err != nil {
					log.Println(err)
					c.AbortWithStatus(http.StatusInternalServerError)
					return
			}
			var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			err = AddTicketToCart(ctx, app.TicketCollection, app.UserCollection, ticketID, userQueryID)
			if err != nil {
					c.IndentedJSON(http.StatusInternalServerError, err)
			}
			c.IndentedJSON(201, "Successfully Added to the cart")
	}
}

func (app *Application) RemoveItem() gin.HandlerFunc {
	return func(c *gin.Context) {
			ticketQueryID := c.Query("id")
			
			if ticketQueryID == "" {
					log.Println("ticket id is invalid")
					_ = c.AbortWithError(http.StatusBadRequest, errors.New("ticket id is empty"))
					return
			}

			userQueryID := c.MustGet("userId").(string)
			if userQueryID == "" {
					log.Println("user id is empty")
					_ = c.AbortWithError(http.StatusBadRequest, errors.New("UserID is empty"))
			}

			ticketID, err := primitive.ObjectIDFromHex(ticketQueryID)
			if err != nil {
					log.Println(err)
					c.AbortWithStatus(http.StatusInternalServerError)
					return
			}

			var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			err = RemoveTicketFromCart(ctx, app.UserCollection, ticketID, userQueryID)
			if err != nil {
					c.IndentedJSON(http.StatusInternalServerError, err)
					return
			}
			c.IndentedJSON(200, "Successfully removed from cart")
	}
}

func GetItemFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "invalid id"})
			c.Abort()
			return
		}

		usert_id, _ := primitive.ObjectIDFromHex(user_id)

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var filledcart models.User
		err := UserCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: usert_id}}).Decode(&filledcart)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(500, "not id found")
			return
		}

		filter_match := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: usert_id}}}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$usercart"}}}}
		grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$usercart.price"}}}}}}
		pointcursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{filter_match, unwind, grouping})
		if err != nil {
			log.Println(err)
		}
		var listing []bson.M
		if err = pointcursor.All(ctx, &listing); err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		for _, json := range listing {
			c.IndentedJSON(200, json["total"])
			c.IndentedJSON(200, filledcart.UserCart)
		}
		ctx.Done()
	}
}

func (app *Application) BuyFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		userQueryID := c.Query("id")
		productIDStr := c.Query("productID")
		productID, err := primitive.ObjectIDFromHex(productIDStr)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if userQueryID == "" {
			log.Panicln("user id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("UserID is empty"))
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		err = BuyProductFromCart(ctx, app.TicketCollection, productID, userQueryID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.IndentedJSON(200, "Successfully Placed the order")
	}
}

func (app *Application) InstantBuy() gin.HandlerFunc {
	s3BucketName := os.Getenv("S3_BUCKET_NAME")
	

	return func(c *gin.Context) {
		UserQueryID := c.Query("userId")
		if UserQueryID == "" {
			log.Println("UserID is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("UserID is empty"))
		}
		ProductQueryID := c.Query("productId")
		if ProductQueryID == "" {
			log.Println("Product_ID id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product_id is empty"))
		}
		productID, err := primitive.ObjectIDFromHex(ProductQueryID)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		// Générer un QR code avec l'ID du produit
		qrCodeData := "ProductID: " + ProductQueryID

		err = utils.GenerateAndUploadQRCode(qrCodeData, s3BucketName, "qrcode.png")		
		if err != nil {
				log.Println("Erreur lors de la génération et de l'enregistrement du QR code:", err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
		}

		err = InstantBuyer(ctx, app.TicketCollection, app.UserCollection, productID, UserQueryID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.IndentedJSON(200, "Successully placed the order")
	}
}

func RemoveProductofCart(ctx context.Context, userCollection *mongo.Collection, productID  primitive.ObjectID, userID string) error {
	filter := bson.D{{"_id", userID}}
	update := bson.D{{"$pull", bson.D{{"usercart", bson.D{{"_id", productID}}}}}}
	_, err := userCollection.UpdateOne(ctx, filter, update)
	return err
}

// BuyProductFromCart purchases a product from the user's cart in the database.
func BuyProductFromCart(ctx context.Context,  userCollection *mongo.Collection, productID  primitive.ObjectID, userID string) error {
	// This function depends on your application logic.
	// You might need to remove the product from the user's cart and decrease the product's stock.
	// Here is a simple example:
	err := RemoveProductofCart(ctx, userCollection, productID, userID)
	if err != nil {
			return err
	}
	// You might need to update the product's stock in another collection.
	// You would need to write that code here.
	return nil
}

// InstantBuyer purchases a product instantly in the database.
func InstantBuyer(ctx context.Context, ticketCollection *mongo.Collection, userCollection *mongo.Collection, productID primitive.ObjectID , userID string) error {
	// This function depends on your application logic.
	// You might need to decrease the product's stock.
	// You might need to update the product's stock in another collection.
	// You would need to write that code here.
	return nil
}

func RemoveTicketFromCart(ctx context.Context, userCollection *mongo.Collection, ticketID  primitive.ObjectID, userID string) error {
	filter := bson.D{{"_id", userID}}
	update := bson.D{{"$pull", bson.D{{"usercart", bson.D{{"_id", ticketID}}}}}}
	_, err := userCollection.UpdateOne(ctx, filter, update)
	return err
}

func AddTicketToCart(ctx context.Context, ticketCollection *mongo.Collection, userCollection *mongo.Collection, ticketID primitive.ObjectID, userID string) error {
	filter := bson.D{{Key: "email", Value: userID}}
	update := bson.D{{Key: "$push", Value: bson.D{{"usercart", bson.D{{"_id", ticketID}}}}}}
	_, err := userCollection.UpdateOne(ctx, filter, update)
	return err
}
