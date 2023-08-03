package main

import (
	"ginapp/middleware"
	"ginapp/routes"
	"os"

	"github.com/gin-gonic/gin"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

func main() {
	port:= os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	// r.Use(cors.Default())
	routes.UserRoutes(r)
	r.Use(middleware.Authentification())
	routes.FoodRoutes(r)
	routes.MenuRoutes(r)
	routes.OrderRoutes(r)
	routes.TableRoutes(r)
	routes.OrderItemRoutes(r)
	routes.OrderItemRoutes(r)
	routes.InvoiceRoutes(r)

	r.Run(":"+ port ) // listen and serve on 0.0.0.0:8080
}