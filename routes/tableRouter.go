package routes

import (
	controllers "ginapp/controllers"

	"github.com/gin-gonic/gin"
)

func TableRoutes(incomingRoutes *gin.Engine) {

	incomingRoutes.GET("/tables", controllers.GetTables())
	incomingRoutes.GET("/table/:id", controllers.GetTableByID())
	incomingRoutes.POST("/table", controllers.PostTable())
	incomingRoutes.PATCH("/table/:id", controllers.UpdateTable())
	incomingRoutes.DELETE("/table/:id", controllers.DeleteTable())
}