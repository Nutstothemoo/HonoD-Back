package routes

import (
	controllers "ginapp/controllers"

	"github.com/gin-gonic/gin"
)

func FoodRoutes (incomingRoutes  *gin.Engine){

	incomingRoutes.GET("/foods", controllers.GetFoods())
	incomingRoutes.GET("/food/:id", controllers.GetFoodByID())
	incomingRoutes.POST("/food", controllers.PostFood())
	incomingRoutes.PATCH("/food/:id", controllers.UpdateFood())
	incomingRoutes.DELETE("/food/:id", controllers.DeleteFood())
}