package routes

import (
	controllers "ginapp/api/controllers"
	middleware "ginapp/api/middleware"
	"github.com/gin-gonic/gin"
)

func TicketRoutes (incomingRoutes * gin.Engine){
	dealerRoutes := incomingRoutes.Group("/dealer")
	dealerRoutes.Use(middleware.DealerAuthentication())
	dealerRoutes.POST("/addproduct", controllers.AddTicketViewerAdmin())

	incomingRoutes.GET("/users/productview", controllers.SearchTicket())
	incomingRoutes.GET("/users/search", controllers.SearchTicketByQuery())
}