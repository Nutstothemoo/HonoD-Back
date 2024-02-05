package sdk

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	// "go.mongodb.org/mongo-driver/bson"
	// "golang.org/x/oauth2"
	"ginapp/api/models"
	generate "ginapp/api/tokens"
	"ginapp/database"
)

var UserCollection *mongo.Collection = database.UserData(database.Client, "Users")



func InitFacebookLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
			var OAuth2Config = GetFacebookOAuthConfig()
			url := OAuth2Config.AuthCodeURL(GetRandomOAuthStateString())
			c.JSON(200, gin.H{"url": url})
	}
}

func HandleFacebookLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
			var state = c.Query("state")
			var code = c.Query("code")
			if state != GetRandomOAuthStateString() {
				c.JSON(500, gin.H{"message": "state is not valid"})
					return
			}

			var OAuth2Config = GetFacebookOAuthConfig()

			token, err := OAuth2Config.Exchange(context.Background(), code)

			if err != nil || token == nil {

					fmt.Println("Error exchanging code for token:", err)
					c.JSON(500, gin.H{"message": "Error exchanging code for token"})
					return
			}

			fbUserDetails, fbUserDetailsError := GetUserInfoFromFacebook(token.AccessToken)

			if fbUserDetailsError != nil {
					fmt.Println("Error getting user details:", fbUserDetailsError)
					c.JSON(500, gin.H{"message": "Error getting user details"})
					return
			}
			fmt.Println("fbUserDetails:", fbUserDetails)
			founduser, err := SignInUser(fbUserDetails, UserCollection)
			if err != nil {
					fmt.Println("Error signing in user:", err)
			} else {
					fmt.Println("Found user:", founduser)
			}
			if founduser.Email == nil || founduser.FirstName == nil || founduser.LastName == nil || founduser.Role == nil {
				fmt.Println("One or more fields of the user are nil")
				return
		}
			authtoken, refreshToken, _ := generate.TokenGenerator(*founduser.Email, *founduser.FirstName, *founduser.LastName, founduser.User_ID, *founduser.Role) 
    

			generate.UpdateAllTokens(authtoken, refreshToken, founduser.User_ID)
			fmt.Println("User ID:", founduser.User_ID)
			http.SetCookie(c.Writer, &http.Cookie{
					Name:     "auth_token",
					Value:    authtoken,
					MaxAge:   60 * 60 * 240,    // 10 day
					HttpOnly: true,             // The cookie is not accessible via JavaScript
					Secure:   false,  // he cookie is not sent only over HTTPS
			})
			fmt.Println("User_ID:", founduser.User_ID)
			fmt.Println("SafeUser Email:", founduser.Email)
			fmt.Println("SafeUser FirstName:", founduser.FirstName)
			fmt.Println("SafeUser LastName:", founduser.LastName)
			fmt.Println("SafeUser Phone:", founduser.Phone)
			fmt.Println("SafeUser Username:", founduser.Username)
			safeUser := SafeUser{
					UserID:    founduser.User_ID,
					Email:     *founduser.Email,
					FirstName: *founduser.FirstName,
					LastName:  *founduser.LastName,
					Phone:     *founduser.Phone,
					Username:  *founduser.Username,
			}
			fmt.Println("SafeUser:", safeUser)
			c.JSON(200, gin.H{
				"user":        safeUser,
			})
	}
}
func SignInUser(userDetail UserDetails, UserCollection *mongo.Collection) (*models.User, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	fmt.Println("User Details:", userDetail) // 
	if userDetail == (UserDetails{}) {
			return nil, errors.New("user details can't be empty")
	}

	if userDetail.Email == "" {
			return nil, errors.New("email can't be empty")
	}

	if userDetail.Name == "" {
			return nil, errors.New("name can't be empty")
	}

	var user models.User
	fmt.Println("Finding user in database:", userDetail.Email) //
	err := UserCollection.FindOne(ctx, bson.M{"email": userDetail.Email}).Decode(&user)
	fmt.Println("Finding user in database1:", user) // Print the error if any
	if err != nil {
			if err == mongo.ErrNoDocuments {
				fmt.Println("User not found in database, creating new user")
					defaultRole := "user"
					emptyString := ""

					user.Role = &defaultRole
					nameParts := strings.Split(userDetail.Name, " ")
					if len(nameParts) > 1 {
							user.FirstName = &nameParts[0]
							user.LastName = &nameParts[1]
					} else {
							user.FirstName = &userDetail.Name

							user.LastName = &emptyString
					}
					user.FirstName = &userDetail.Name
					user.Email = &userDetail.Email
					user.Password = &emptyString
					user.FacebookID = &userDetail.ID
					user.Phone = &emptyString
					user.Username = &emptyString;
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
							return nil, inserterr
					}
					
					fmt.Println("User created successfully")
					return &user, nil
			}
			fmt.Println("Error finding user in database:", err) // Print the error if any

			return nil, err
	}
	
	return &user, nil
}

