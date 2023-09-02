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

	// incomingRoutes.GET("/users", controllers.GetUsers())
	// incomingRoutes.GET("/users/:id", controllers.GetUserByID())
	// incomingRoutes.PUT("/users/:id", controllers.UpdateUser())
	// incomingRoutes.DELETE("/users/:id", controllers.DeleteUser())
}