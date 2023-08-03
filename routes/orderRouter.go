package routes

import (
	controllers "ginapp/controllers"

	"github.com/gin-gonic/gin"
)

func OrderRoutes (incomingRoutes  *gin.Engine){

	incomingRoutes.GET("/orders", controllers.GetOrders())
	incomingRoutes.GET("/order/:id", controllers.GetOrderByID())
	incomingRoutes.POST("/order", controllers.PostOrder())
	incomingRoutes.PATCH("/order/:id", controllers.UpdateOrder())
	incomingRoutes.DELETE("/order/:id", controllers.DeleteOrder())
}