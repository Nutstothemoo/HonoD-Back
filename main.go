package main

import (
	"ginapp/api/controllers"
	"ginapp/database"
	"ginapp/api/middleware"
	"ginapp/api/routeur"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
    log.Println("Error loading .env file")
  }
	port:= os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	app:= controllers.NewApplication(database.ProductData(database.Client, "Products"), database.UserData(database.Client, "Users"))
	
	r := gin.New()
	r.Use(gin.Logger())
	
	// r.Use(cors.Default())
	r.Use(middleware.Authentification())
	routes.UserRoutes(r)
	routes.EventRoutes(r)
	routes.TicketRoutes(r)
	// routes.CartRoutes(r)

	r.Use(gin.Recovery())	
	log.Println("http://localhost:" + port)
	r.Run("localhost:"+ port ) 
}