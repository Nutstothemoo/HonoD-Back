package controllers

import (
	"context"
	"fmt"
	"ginapp/api/models"
	generate "ginapp/api/tokens"
	"ginapp/api/utils"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)
type SafeUser struct {
	Email       string `json:"email"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Phone       string `json:"phone"`
	Username    string `json:"username"`
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		if err := c.BindJSON(&user); err != nil {
			fmt.Println("Error binding JSON: ", err) // Added logging
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Set the default role to "user"
		defaultRole := "user"
		user.Role = &defaultRole

		validationErr := Validate.Struct(user)
		if validationErr != nil {
			fmt.Println("Validation error: ", validationErr) 
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr})
			return
		}

		count, err := UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			fmt.Println("Error counting documents: ", err) // Added logging

			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		if count > 0 {
			fmt.Println("User already exists") // Added logging
			c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		}
		count, err = UserCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel()
		if err != nil {			
			fmt.Println("Phone is already in use")
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Phone is already in use"})
			return
		}
		password := utils.HashPassword(*user.Password)
		user.Password = &password
		user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_ID = user.ID.Hex()
		token, refreshtoken, _ := generate.TokenGenerator(*user.Email, *user.FirstName, *user.LastName, user.User_ID, *user.Role)
		user.Token = &token
		user.Refresh_Token = &refreshtoken
		user.UserCart = make([]models.TicketUser, 0)
		user.Address_Details = make([]models.Address, 0)
		user.Order_History = make([]models.Order, 0)
		user.Order_Refunded = make([]models.Order, 0)
		user.Order_Canceled = make([]models.Order, 0)		
		user.Order_History = make([]models.Order,0)

		_, inserterr := UserCollection.InsertOne(ctx, user)
		if inserterr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "not created"})
			return
		}
		http.SetCookie(c.Writer, &http.Cookie{
            Name:     "auth_token",
            Value:    token,
            MaxAge:   60 * 60 * 240,    // 10 day
            HttpOnly: true,            // The cookie is not accessible via JavaScript
            Secure:   false,            // The cookie is not sent only over HTTPS
            SameSite: http.SameSiteStrictMode, // The cookie is sent only to the same site as the one that originated it
        })

        // Set the refresh_token cookie
        http.SetCookie(c.Writer, &http.Cookie{
            Name:     "refresh_token",
            Value:    refreshtoken,
            MaxAge:   60 * 60 * 240,    // 10 day
            HttpOnly: true,            // The cookie is not accessible via JavaScript
            Secure:   false,            // The cookie is not sent only over HTTPS
            SameSite: http.SameSiteStrictMode, // The cookie is sent only to the same site as the one that originated it
        })
		defer cancel()
		c.JSON(201, "Successfully Signed Up!!")
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		var founduser models.User
		if err := c.BindJSON(&user); err != nil {
			fmt.Println("Error binding JSON: ", err) 
			c.JSON(http.StatusBadRequest, gin.H{"Incorrect entry ": err})
			return
		}
		err := UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&founduser)
		defer cancel()
		if err != nil {
			
			fmt.Println("Error finding user: ", err) 
			c.JSON(http.StatusInternalServerError, gin.H{"error": "login or password incorrect"})
			return
		}
		PasswordIsValid, msg := utils.VerifyPassword(*user.Password, *founduser.Password)
		defer cancel()
		if !PasswordIsValid {
			fmt.Println("Password is not valid: ", msg) 

			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			fmt.Println(msg)
			return
		}
		token, refreshToken, _ := generate.TokenGenerator(*founduser.Email, *founduser.FirstName, *founduser.LastName, founduser.User_ID, *founduser.Role)		
		defer cancel()
		generate.UpdateAllTokens(token, refreshToken, founduser.User_ID)
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "auth_token",
			Value:    token,
			MaxAge:   60 * 60 * 240,    // 10 day
			HttpOnly: true,            // The cookie is not accessible via JavaScript
			Secure:   false,            // The cookie is not sent only over HTTPS
			SameSite: http.SameSiteStrictMode, // The cookie is sent only to the same site as the one that originated it
		})

		http.SetCookie(c.Writer, &http.Cookie{
            Name:     "refresh_token",
            Value:    refreshToken,
            MaxAge:   60 * 60 * 240,    // 10 day
            HttpOnly: true,            // The cookie is not accessible via JavaScript
            Secure:   false,            // The cookie is not sent only over HTTPS
            SameSite: http.SameSiteStrictMode, // The cookie is sent only to the same site as the one that originated it
        })

		safeUser := SafeUser{
			Email:     *founduser.Email,
			FirstName: *founduser.FirstName,
			LastName:  *founduser.LastName,
			Phone:     *founduser.Phone,
			Username:  *founduser.Username,
		}

		c.JSON(200, safeUser)
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		user_id := c.Query("id")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "invalid id"})
			c.Abort()
			return
		}
		// Get the user ID from the Gin context
		loggedInUserId := c.MustGet("userId").(string)

		if user_id == "" || loggedInUserId != user_id {
				c.Header("Content-Type", "application/json")
				c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to view this user"})
				c.Abort()
				return
		}
		usert_id, _ := primitive.ObjectIDFromHex(user_id)
		err := UserCollection.FindOne(ctx, bson.M{"_id": usert_id}).Decode(&user)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(500, "not id found")
			return
		}
	}
}

func DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		user_id := c.Query("id")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "invalid id"})
			c.Abort()
			return
		}

		loggedInUserId := c.MustGet("userId").(string)

		if user_id == "" || loggedInUserId != user_id {
				c.Header("Content-Type", "application/json")
				c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to update this user"})
				c.Abort()
				return
		}

		err := UserCollection.FindOneAndDelete(ctx, bson.M{"_id": user_id}).Decode(&models.User{})
		if err != nil {
			log.Println(err)
			c.IndentedJSON(500, "not id found")
			return
		}		
	}
}

func UpdateUser() gin.HandlerFunc {
	return func (c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		user_id := c.Query("id")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "invalid id"})
			c.Abort()
			return
		}

		loggedInUserId := c.MustGet("userId").(string)

		if user_id == "" || loggedInUserId != user_id {
				c.Header("Content-Type", "application/json")
				c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to update this user"})
				c.Abort()
				return
		}

		usert_id, _ := primitive.ObjectIDFromHex(user_id)
		err := UserCollection.FindOne(ctx, bson.M{"_id": usert_id}).Decode(&user)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(500, "not id found")
			return
		}
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		validationErr := Validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr})
			return
		}
		user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		_, err = UserCollection.UpdateOne(ctx, bson.M{"_id": usert_id}, bson.M{"$set": user})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "not updated"})
			return
		}
		c.JSON(http.StatusOK, "Successfully updated the user")
	}
}

func AdminUpdateUser() gin.HandlerFunc {
	return func (c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		user_id := c.Query("id")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "invalid id"})
			c.Abort()
			return
		}
		usert_id, _ := primitive.ObjectIDFromHex(user_id)
		err := UserCollection.FindOne(ctx, bson.M{"_id": usert_id}).Decode(&user)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(500, "not id found")
			return
		}
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		validationErr := Validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr})
			return
		}
		user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		_, err = UserCollection.UpdateOne(ctx, bson.M{"_id": usert_id}, bson.M{"$set": user})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "not updated"})
			return
		}
		c.JSON(http.StatusOK, "Successfully updated the user")
	}
}

func AdminDeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		user_id := c.Query("id")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "invalid id"})
			c.Abort()
			return
		}
		
		err := UserCollection.FindOneAndDelete(ctx, bson.M{"_id": user_id}).Decode(&models.User{})
		if err != nil {
			log.Println(err)
			c.IndentedJSON(500, "not id found")
			return
		}		
	}
}

func AdminGetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		user_id := c.Query("id")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "invalid id"})
			c.Abort()
			return
		}
		
		usert_id, _ := primitive.ObjectIDFromHex(user_id)
		err := UserCollection.FindOne(ctx, bson.M{"_id": usert_id}).Decode(&user)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(500, "not id found")
			return
		}
	}
}

func RefreshToken() gin.HandlerFunc {
    return func(c *gin.Context) {
        var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
        defer cancel()

        // Extract the refresh token from the request
        refreshToken := c.Request.Header.Get("Refresh-Token")
        if refreshToken == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "No refresh token provided"})
            return
        }

        // Verify the refresh token
        email, err := generate.VerifyRefreshToken(refreshToken)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
            return
        }

        // Find the user associated with the refresh token
        var user models.User
        err = UserCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
            return
        }

        // Generate a new token and refresh token
        token, newRefreshToken, _ := generate.TokenGenerator(*user.Email, *user.FirstName, *user.LastName, user.User_ID, *user.Role)

        // Update the user's tokens in the database
        generate.UpdateAllTokens(token, newRefreshToken, user.User_ID)

        // Return the new tokens
        c.JSON(http.StatusOK, gin.H{"token": token, "refresh_token": newRefreshToken})
    }
}
