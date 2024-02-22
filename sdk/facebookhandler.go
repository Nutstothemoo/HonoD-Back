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

var UserCollection *mongo.Collection = database.OpenCollection(database.Client, "Users")

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

							c.JSON(500, gin.H{"message": "Error exchanging code for token"})
					return
			}

			fbUserDetails, fbUserDetailsError := GetUserInfoFromFacebook(token.AccessToken)
			if fbUserDetailsError != nil {
					c.JSON(500, gin.H{"message": "Error getting user details"})
					return
			}
			founduser, err := SignInUser(fbUserDetails, UserCollection)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
						"error": fmt.Sprintf("Error signing in user: %v", err),
				})
				return
		} 
		if founduser.Email == nil || founduser.FirstName == nil || founduser.LastName == nil || founduser.Role == nil {
			c.JSON(http.StatusBadRequest, gin.H{
					"error": "One or more fields of the user are nil",
			})
			return
		}
		authtoken, refreshToken, _ := generate.TokenGenerator(*founduser.Email, *founduser.FirstName, *founduser.LastName, founduser.User_ID, *founduser.Role) 
    

		generate.UpdateAllTokens(authtoken, refreshToken, founduser.User_ID)
		http.SetCookie(c.Writer, &http.Cookie{
					Name:     "auth_token",
					Value:    authtoken,
					MaxAge:   60 * 60 * 240,    // 10 day
					HttpOnly: false,             // The cookie is not accessible via JavaScript
					Secure:   false, 
			})
			safeUser := SafeUser{
					UserID:    founduser.User_ID,
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

func SignInUser(userDetail UserDetails, UserCollection *mongo.Collection) (*models.User, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

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
	err := UserCollection.FindOne(ctx, bson.M{"email": userDetail.Email}).Decode(&user)
	if err != nil {
			if err == mongo.ErrNoDocuments {
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
					
					return &user, nil
			}
			return nil, err
	}
	return &user, nil
}

