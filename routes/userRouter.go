package routes

import (
	controllers "ginapp/controllers"

	"github.com/gin-gonic/gin"
)
func UserRoutes (incomingRoutes  *gin.Engine){
	
	incomingRoutes.GET("/users", controllers.GetUsers())
	incomingRoutes.GET("/users/:id", controllers.GetUserByID())
	incomingRoutes.POST("/users/signup", controllers.SignUp())
	incomingRoutes.POST("/users/login", controllers.Login())
	incomingRoutes.PUT("/users/:id", controllers.UpdateUser())
	incomingRoutes.DELETE("/users/:id", controllers.DeleteUser())
}