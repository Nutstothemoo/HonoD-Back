package main

import (
	"ginapp/controllers"
	"ginapp/database"
	// "ginapp/middleware"
	"ginapp/routes"
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
	routes.UserRoutes(r)

	// r.Use(middleware.Authentification())

	r.GET("/products", app.AddToCart())
	r.GET("/removeitem", app.RemoveItem())
	r.GET("/cartcheckout", app.BuyFromCart())
	r.GET("/instantbuy", app.InstantBuy())
	r.Use(gin.Recovery())	
	log.Println("http://localhost:" + port)
	r.Run("localhost:"+ port ) // listen and serve on 0.0.0.0:8080
}