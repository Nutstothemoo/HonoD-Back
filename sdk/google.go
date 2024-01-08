package sdk

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const JWTTokenURL = "https://oauth2.googleapis.com/token"
const MTLSTokenURL = "https://oauth2.mtls.googleapis.com/token"
var Endpoint = oauth2.Endpoint{
	AuthURL:       "https://accounts.google.com/o/oauth2/auth",
	TokenURL:      "https://oauth2.googleapis.com/token",
	DeviceAuthURL: "https://oauth2.googleapis.com/device/code",
	AuthStyle:     oauth2.AuthStyleInParams,
}

func getGoogleAuthConfig() *oauth2.Config {
	return &oauth2.Config{
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("FRONT_URL") + os.Getenv("GOOGLE_REDIRECT_URL"),
			Scopes: []string{
					"https://www.googleapis.com/auth/userinfo.profile",
					"https://www.googleapis.com/auth/userinfo.email",
			},
			Endpoint: google.Endpoint,
	}
}

func GetUserInfoFromGoogle(token string) (UserDetails, error) {
	var googleUserDetails UserDetails
	googleUserDetailsRequest, _ := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	googleUserDetailsRequest.Header.Add("Authorization", "Bearer "+token)

	googleUserDetailsResponse, googleUserDetailsResponseError := http.DefaultClient.Do(googleUserDetailsRequest)

	if googleUserDetailsResponseError != nil {
			return UserDetails{}, errors.New("Erreur lors de la récupération des informations depuis Google")
	}
	defer googleUserDetailsResponse.Body.Close()

	decoder := json.NewDecoder(googleUserDetailsResponse.Body)
	decoderErr := decoder.Decode(&googleUserDetails)

	if decoderErr != nil {
			return UserDetails{}, errors.New("Erreur lors de la récupération des informations depuis Google")
	}

	return googleUserDetails, nil
}
