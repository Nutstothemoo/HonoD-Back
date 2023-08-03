package controller

import (
	"context"
	"ginapp/database"
	"ginapp/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")
var validate = validator.New()

func GetFoods() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		foodId := c.Param("id")
		var food models.Food
		err:= foodCollection.FindOne(ctx, bson.M{"_id": foodId}).Decode(&food)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error while getting the food",
			})
		}
		c.JSON(http.StatusOK, food)
	}
}

func GetFoodByID() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
func PostFood() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func UpdateFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		
	}
}
func DeleteFood() gin.HandlerFunc {
	return func(c *gin.Context) {		
	}
}

func round (num float64) int {
	
}

func toFixed (num float64, precision int ) float64 {
	
}