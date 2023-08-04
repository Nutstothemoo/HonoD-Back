package controllers

import (
	"context"
	"fmt"
	"ginapp/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)


func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var UserCollection *mongo.Collection = UserData(client, "users")
		var user models.User
		if err:= c.BindJSON(&user); err!=nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
		}
		validationError := Validate.Struct(user)
		if validationError != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": validationError})
				return,
		}

		count, err := UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error":err})
			return 
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})	
		}
		count, err:= UserCollection.CountDocuments(ctx,bson.M{"phone": userPhone})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error":err})
			return 
		}


}
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}
func UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}
func DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("DeleteUser")
	}
}
func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}
func GetUserByID() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}
func SearchProductByQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func SearchProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}
func ProductViewerAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func HashPassword( password string) string {
}

func VerifyPassword( Userpassword string, providedPassword string) bool {	
}