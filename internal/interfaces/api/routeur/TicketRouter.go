package routes

import (
	controllers "ginapp/internal/interfaces/api/controllers"
	middleware "ginapp/internal/infrastructure/middleware"
	"github.com/gin-gonic/gin"
)

func TicketRoutes (incomingRoutes * gin.Engine){

	incomingRoutes.GET("/users/productview", controllers.SearchTicket())
	incomingRoutes.GET("/users/search", controllers.SearchTicketByQuery())
	incomingRoutes.GET("/tickets/:id", controllers.GetTickets())
	dealerRoutes := incomingRoutes.Group("/dealer")

	dealerRoutes.Use(middleware.DealerAuthentication())
	dealerRoutes.POST("/tickets/:eventId", controllers.AddTicket())
	dealerRoutes.PUT("/tickets/:id", controllers.UpdateTicket())
	dealerRoutes.DELETE("/tickets/:id", controllers.DeleteTicket())

}

