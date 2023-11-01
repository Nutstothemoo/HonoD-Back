package routes

import (
	controllers "ginapp/api/controllers"

	"github.com/gin-gonic/gin"
)

func TicketRoutes (incomingRoutes * gin.Engine){

	incomingRoutes.POST("/admin/addproduct", controllers.AddTicketViewerAdmin())
	incomingRoutes.GET("/users/productview", controllers.SearchTicket())
	incomingRoutes.GET("/users/search", controllers.SearchTicketByQuery())
}