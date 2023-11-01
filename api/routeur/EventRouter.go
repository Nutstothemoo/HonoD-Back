		package routes

	import (
		controllers "ginapp/api/controllers"
	
		"github.com/gin-gonic/gin"
	)
	
	func EventRoutes (incomingRoutes * gin.Engine){
	
		incomingRoutes.GET("/events", controllers.GetEvents())
		incomingRoutes.GET("/events/:id", controllers.GetEventByID())
		incomingRoutes.POST("/events", controllers.AddEvent())
		incomingRoutes.PUT("/events/:id", controllers.UpdateEvent())
		incomingRoutes.DELETE("/events/:id", controllers.DeleteEvent())
		incomingRoutes.GET("/events/:fromDate/:toDate", controllers.GetEventFromDateToDate())		

	}
	

	