package routes

import (
	controllers "ginapp/api/controllers"
	middleware "ginapp/api/middleware"
	"github.com/gin-gonic/gin"
)

func TicketRoutes (incomingRoutes * gin.Engine){

	incomingRoutes.GET("/users/productview", controllers.SearchTicket())
	incomingRoutes.GET("/users/search", controllers.SearchTicketByQuery())

	dealerRoutes := incomingRoutes.Group("/dealer")

	dealerRoutes.Use(middleware.DealerAuthentication())
	dealerRoutes.POST("/tickets", controllers.AddTicket())
	dealerRoutes.PUT("/tickets/:id", controllers.UpdateTicket())
	dealerRoutes.DELETE("/tickets/:id", controllers.DeleteTicket())

}

