package sdk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"ginapp/api/models"
	generate "ginapp/api/tokens"
	"io"

	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func InitGoogleLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
			var OAuth2Config = getGoogleAuthConfig()
			url := OAuth2Config.AuthCodeURL(GetRandomOAuthStateString())
			c.JSON(200, gin.H{"url": url})
	}
}

func HandleGoogleLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		state := c.Query("state")
		code := c.Query("code")
		frontURL := os.Getenv("FRONT_URL")

		if state != GetRandomOAuthStateString() {
			c.JSON(500, gin.H{"message": "state is not valid"})
			return
		}

		OAuth2Config := getGoogleAuthConfig()

		token, err := OAuth2Config.Exchange(context.Background(), code)
		if err != nil || token == nil {
			c.JSON(500, gin.H{"message": "Error exchanging code for token"})
			c.Redirect(http.StatusTemporaryRedirect, frontURL+"/")
			return
		}

		client := OAuth2Config.Client(context.Background(), token)

		// Make a request to the Google API (e.g., userinfo)
		response, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
		if err != nil {
			c.JSON(500, gin.H{"message": "Error getting user from Google"})
			return
		}
		defer response.Body.Close()

		// Read the response body
		apiResponse, err := io.ReadAll(response.Body)
		if err != nil {
			c.JSON(500, gin.H{"message": "Error getting user from Google"})
			return
		}

		// Parse the API response into UserDetails struct
		var userDetail GoogleUserDetails
		if err := json.Unmarshal(apiResponse, &userDetail); err != nil {
			c.JSON(500, gin.H{"message": "Error getting user from Google"})
			return
		}
		// Pass the UserDetails to SignInUser function
		founduser, err := GoogleSignInUser(userDetail, UserCollection)
		fmt.Println("founduser ", founduser);
		if err != nil {
			c.Redirect(http.StatusTemporaryRedirect, frontURL+"/")
			return
		}
		if err != nil {
			log.Println("Error signing in user:", err)
			c.Redirect(http.StatusTemporaryRedirect, frontURL + "/")
			return
	}

	authtoken, refreshToken, _ := generate.TokenGenerator(*founduser.Email, *founduser.FirstName, *founduser.LastName, founduser.User_ID, *founduser.Role)      
	generate.UpdateAllTokens(authtoken, refreshToken, founduser.User_ID)
	http.SetCookie(c.Writer, &http.Cookie{
			Name:     "auth_token",
			Value:    authtoken,
			MaxAge:   60 * 60 * 240,    // 10 day
			HttpOnly: false,            // The cookie is not accessible via JavaScript
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
	c.JSON(200, gin.H{
		"user":        safeUser,
	})
	}
}

func GoogleSignInUser(userDetails GoogleUserDetails, UserCollection *mongo.Collection) (*models.User, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	if userDetails == (GoogleUserDetails{}) {
		return nil, errors.New("Google user details can't be empty")
	}

	if userDetails.Email == "" {
		return nil, errors.New("email can't be empty")
	}

	if userDetails.Name == "" {
		return nil, errors.New("name can't be empty")
	}

	var user models.User
	err := UserCollection.FindOne(ctx, bson.M{"email": userDetails.Email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			defaultRole := "user"
			emptystring := ""	

			user.Role = &defaultRole
			user.LastName = &userDetails.FamilyName
			user.FirstName = &userDetails.GivenName
			user.Username = &emptystring;
			user.Email = &userDetails.Email
			user.Avatar = userDetails.Picture
			user.GoogleID = &userDetails.ID
			user.Password = &emptystring
			user.Phone = &emptystring
			user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			user.ID = primitive.NewObjectID()
			user.User_ID = user.ID.Hex()
			token, refreshtoken, _ := generate.TokenGenerator(*user.Email, userDetails.GivenName, userDetails.FamilyName, user.User_ID, *user.Role)
			user.Token = &token
			user.Refresh_Token = &refreshtoken
			user.UserCart = make([]models.TicketUser, 0)
			user.Address_Details = make([]models.Address, 0)
			user.Order_History = make([]models.Order, 0)
			user.Order_Refunded = make([]models.Order, 0)
			user.Order_Canceled = make([]models.Order, 0)
			user.Order_History = make([]models.Order, 0)

			_, inserterr := UserCollection.InsertOne(ctx, user)
			if inserterr != nil {
				return nil, inserterr
			}

			return &user, nil
		}
		return nil, err
	}
	return &user, nil
}
