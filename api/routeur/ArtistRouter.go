package routes

import (
	controllers "ginapp/api/controllers"

	"github.com/gin-gonic/gin"
)

func ArtistRoutes (incomingRoutes * gin.Engine){

	incomingRoutes.GET("/artists", controllers.GetArtists())
	incomingRoutes.GET("/artist/:id", controllers.GetArtistByID())
	incomingRoutes.POST("/artist", controllers.AddArtist())
	incomingRoutes.PUT("/artist/:id", controllers.UpdateArtist())
	incomingRoutes.DELETE("/artist/:id", controllers.DeleteArtist())	
}