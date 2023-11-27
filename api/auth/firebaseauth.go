package auth

import (
    "context"
    "firebase.google.com/go"
    "log"
    "net/http"
)

var (
    app *firebase.App
)

func GoogleAuthHandler(w http.ResponseWriter, r *http.Request) {
    ctx := context.Background()
    client, err := app.Auth(ctx)
    if err != nil {
        log.Fatalf("error getting Auth client: %v\n", err)
    }

    // Get the ID token from the request body
    idToken := r.FormValue("idToken")

    // Verify the ID token
    token, err := client.VerifyIDToken(ctx, idToken)
    if err != nil {
        log.Fatalf("error verifying ID token: %v\n", err)
    }

    // Token is valid, get the UID
    uid := token.UID
    log.Printf("Verified ID token: %v\n", uid)
}

func FacebookAuthHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	client, err := app.Auth(ctx)
	if err != nil {
			log.Fatalf("error getting Auth client: %v\n", err)
	}

	// Get the Facebook access token from the request body
	accessToken := r.FormValue("accessToken")

	// Verify the access token with Facebook
	token, err := client.VerifyIDTokenAndCheckRevoked(ctx, accessToken)
	if err != nil {
			log.Fatalf("error verifying Facebook access token: %v\n", err)
	}

	// Token is valid, get the UID
	uid := token.UID
	log.Printf("Verified Facebook access token: %v\n", uid)
}