package main

import (
	"context"
	"firebase.google.com/go"
	"ginapp/api/controllers"
	"ginapp/api/middleware"
	"ginapp/api/routeur"
	"ginapp/database"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"	
	"ginapp/api/auth" 
	"log"
	"fmt"
	"os"
	"github.com/fatih/color"
)

var (
	app *firebase.App
)

func init() {
	opt := option.WithCredentialsFile("credential.json")
	var err error
	app, err = firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
			log.Fatalf("error initializing app: %v\n", err)
	}
	log.Println(color.GreenString("Successfully connected to Firebase"))
}

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
	)

	r := setupRouter(app)
	fmt.Println("╔════════════════════════════════════════╗")
	log.Println(color.GreenString("http://localhost:" + port))
	fmt.Println("╚════════════════════════════════════════╝")
	r.Run("localhost:"+ port ) 
}

func setupRouter(app *controllers.Application) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.POST("/googleAuth", gin.WrapF(auth.GoogleAuthHandler))
	r.POST("/facebookAuth", gin.WrapF(auth.FacebookAuthHandler))
	routes.UserRoutes(r)
	routes.EventRoutes(r)
	routes.TicketRoutes(r)
	// routes.AddressRoutes(r)
	routes.ArtistRoutes(r)
	// routes.AdminRoutes(r) 

	adminRoutes := r.Group("/admin")
	adminRoutes.Use(middleware.AdminAuthentication())
	routes.AdminEventRoutes(adminRoutes)
	adminRoutes.POST("/event", controllers.AdminAddEvent())
	adminRoutes.PUT("/event/:id", controllers.AdminUpdateEvent() )
	adminRoutes.DELETE("/event/:id", controllers.AdminDeleteEvent())

	adminRoutes.PUT("/users/:id", controllers.AdminUpdateUser())
	adminRoutes.DELETE("/users/:id", controllers.AdminDeleteUser())
	adminRoutes.GET("/users/:id", controllers.AdminGetUser())

	adminRoutes.POST("/ticket", controllers.AdminAddTicket())
	adminRoutes.PUT("/ticket/:id", controllers.AdminUpdateTicket())
	adminRoutes.DELETE("/ticket/:id", controllers.AdminDeleteTicket())
	// adminRoutes.GET("/ticket/:id", controllers.AdminGetTicket())

	// USER ROUTE

	routes.UserRoutes(r)
	routes.EventRoutes(r)
	routes.TicketRoutes(r)
	routes.AdressRoutes(r)
	routes.ArtistRoutes(r)

	r.Use(gin.Recovery())

	return r
}