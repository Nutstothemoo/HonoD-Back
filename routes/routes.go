package routes

import (
	controllers "ginapp/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes (incomingRoutes * gin.Engine){

	incomingRoutes.POST("/users/signup", controllers.SignUp())
	incomingRoutes.POST("/users/login", controllers.Login())
	incomingRoutes.POST("/admin/addproduct", controllers.AddProductViewerAdmin())

	incomingRoutes.GET("/users/productview", controllers.SearchProduct())
	incomingRoutes.GET("/users/search", controllers.SearchProductByQuery())
	incomingRoutes.GET("/Buy/:userId/:productId", controllers.InstantBuy())

	// incomingRoutes.GET("/users", controllers.GetUsers())
	// incomingRoutes.GET("/users/:id", controllers.GetUserByID())
	// incomingRoutes.PUT("/users/:id", controllers.UpdateUser())
	// incomingRoutes.DELETE("/users/:id", controllers.DeleteUser())
	// GET ENVENTS BY FILTER
	// GET ONE EVENT
	// POST ONE EVENT
	// UPDATE ONE EVENT 
	// DELETE ONE EVENT
	// GET LAST 25 EVENT BY DATE

	
}