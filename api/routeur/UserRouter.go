package routes

import (
	controllers "ginapp/api/controllers"
	middleware "ginapp/api/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {

	// Public routes

	incomingRoutes.POST("/users/signup", controllers.SignUp())
	incomingRoutes.POST("/users/login", controllers.Login())
	incomingRoutes.POST("/refresh-token", controllers.RefreshToken())
	
	// Authenticated routes

	authGroup := incomingRoutes.Group("/users")
	authGroup.Use(middleware.Authentication())

	authGroup.PUT("/:id", controllers.UpdateUser())
	authGroup.DELETE("/:id", controllers.DeleteUser())

}