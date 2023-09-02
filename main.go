package main

import (
	"ginapp/controllers"
	"ginapp/database"
	"ginapp/middleware"
	"ginapp/routes"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "Products")

func main() {
	port:= os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	app:= controllers.NewApplication(database.ProductData(database.Client, "Products"), database.UserData(database.Client, "Users"))
	
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	// r.Use(cors.Default())
	routes.UserRoutes(r)
	r.Use(middleware.Authentification())

	r.GET("/products", app.AddToCart())
	r.GET("/removeitem", app.RemoveItem())
	r.GET("/cartcheckout", app.BuyFromCart())
	r.GET("/instantbuy", app.InstantBuy())
	

	r.Run(":"+ port ) // listen and serve on 0.0.0.0:8080
}