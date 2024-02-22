package main

import (
	"context"
	"fmt"

	// "ginapp/api/auth"
	"ginapp/api/controllers"
	"ginapp/api/middleware"
	"ginapp/api/routeur"
	"ginapp/database"
	"ginapp/sdk"
	"log"
	"net/http"

	// "golang.org/x/oauth2"
	// "golang.org/x/oauth2/google"
	"os"
	"firebase.google.com/go"
	"github.com/fatih/color"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
    log.Println(color.RedString("Error loading .env file"))
  }
	port:= os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	app:= controllers.NewApplication(
		database.OpenCollection(database.Client, "Tickets"), 
		database.OpenCollection(database.Client, "Users"),		
		database.OpenCollection(database.Client, "Events"),
		database.OpenCollection(database.Client, "TicketUnit"),
	)

	r := setupRouter(app)
	fmt.Println(color.GreenString("╔════════════════════════════════════════╗"))
	log.Println(color.GreenString("http://localhost:" + port))
	fmt.Println(color.GreenString("╚════════════════════════════════════════╝"))
	r.Run("localhost:"+ port ) 
}

func setupRouter(app *controllers.Application) *gin.Engine {
		r := gin.New()
		r.Use(gin.Logger())

		// Add CORS middleware

		config := cors.DefaultConfig()
		config.AllowAllOrigins = true // Allow all origins
		config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"} // Allow all methods
		config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"} // Allow all headers
		
		r.Use(cors.New(config))

		r.GET("/ping", func(c *gin.Context) {
			c.String(http.StatusOK, "pong")
		})

	r.GET("/login/facebook", sdk.InitFacebookLogin())
	r.GET("/facebook/callback", sdk.HandleFacebookLogin())

	r.GET("/login/google", sdk.InitGoogleLogin())
	r.GET("/google/callback", sdk.HandleGoogleLogin())

	r.POST("/instantBuy/:ticketId/:eventId", middleware.Authentication(), app.InstantBuy())

	routes.UserRoutes(r)
	routes.EventRoutes(r)
	routes.TicketRoutes(r)
	routes.AdressRoutes(r)
	routes.ArtistRoutes(r)

	adminRoutes := r.Group("/admin")
	adminRoutes.Use(middleware.AdminAuthentication())
	routes.AdminEventRoutes(adminRoutes)
	adminRoutes.PUT("/users/:id", controllers.AdminUpdateUser())
	adminRoutes.DELETE("/users/:id", controllers.AdminDeleteUser())
	adminRoutes.GET("/users/:id", controllers.AdminGetUser())

	adminRoutes.POST("/ticket", controllers.AdminAddTicket())
	adminRoutes.PUT("/ticket/:id", controllers.AdminUpdateTicket())
	adminRoutes.DELETE("/ticket/:id", controllers.AdminDeleteTicket())
	// adminRoutes.GET("/ticket/:id", controllers.AdminGetTicket())

	// USER ROUTE


	r.Use(gin.Recovery())

	return r
}

var (
	app *firebase.App
)

func firebaseInit() {
	opt := option.WithCredentialsFile("credential.json")
	var err error
	app, err = firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
			log.Fatalf("error initializing app: %v\n", err)
	}
	log.Println(color.GreenString("Successfully connected to Firebase"))
}
