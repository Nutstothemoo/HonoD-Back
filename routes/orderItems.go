package routes

import (
	controllers "ginapp/controllers"

	"github.com/gin-gonic/gin"
)

func OrderItemRoutes (incomingRoutes  *gin.Engine){

	incomingRoutes.GET("/orderItems", controllers.GetOrderItems())
	incomingRoutes.GET("/orderItem/:orderItems_id", controllers.GetOrderItemByID())
	incomingRoutes.GET("/orderItem-Order/:order_id", controllers.GetOrderItemByOrderID())
	incomingRoutes.POST("/orderItem", controllers.PostOrderItem())
	incomingRoutes.PATCH("/orderItem/:orderItems_id", controllers.UpdateOrderItem())
	incomingRoutes.DELETE("/orderItem/:orderItems_id", controllers.DeleteOrderItem())
}