package routes

import (
	controllers "ginapp/api/controllers"
	middleware "ginapp/api/middleware"
	"github.com/gin-gonic/gin"
)

func ArtistRoutes (incomingRoutes * gin.Engine){

	incomingRoutes.GET("/artists", controllers.GetArtists())
	incomingRoutes.GET("/artist/:id", controllers.GetArtistByID())
	dealerRoutes := incomingRoutes.Group("/dealer")
	dealerRoutes.Use(middleware.DealerAuthentication())
	dealerRoutes.POST("/artist", controllers.AddArtist())
	dealerRoutes.PUT("/artist/:id", controllers.UpdateArtist())
	dealerRoutes.DELETE("/artist/:id", controllers.DeleteArtist())	
}


func AdminArtistRoutes (incomingRoutes * gin.Engine){

	incomingRoutes.GET("/artists", controllers.GetArtists())
	incomingRoutes.GET("/artist/:id", controllers.GetArtistByID())
	incomingRoutes.POST("/artist", controllers.AddArtist())
	incomingRoutes.PUT("/artist/:id", controllers.UpdateArtist())
	incomingRoutes.DELETE("/artist/:id", controllers.DeleteArtist())	
}