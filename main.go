package main

import (
	"context"
	"fmt"
	"ginapp/api/auth"
	"ginapp/api/controllers"
	"ginapp/api/middleware"
	"ginapp/api/routeur"
	"ginapp/database"
	"log"
	"net/http"
	"os"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"firebase.google.com/go"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
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

		// Add CORS middleware

		config := cors.DefaultConfig()
		config.AllowAllOrigins = true // Allow all origins
		config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"} // Allow all methods
		config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"} // Allow all headers
		r.Use(cors.New(config))

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	
	routes.UserRoutes(r)
	routes.EventRoutes(r)
	routes.TicketRoutes(r)
	routes.AdressRoutes(r)
	routes.ArtistRoutes(r)

	r.POST("/googleAuth", gin.WrapF(auth.GoogleAuthHandler))
	r.POST("/facebookAuth", gin.WrapF(auth.FacebookAuthHandler))

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


	r.Use(gin.Recovery())

	return r
}