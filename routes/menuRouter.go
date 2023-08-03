package routes

import (
	controllers "ginapp/controllers"

	"github.com/gin-gonic/gin"
)

func MenuRoutes (incomingRoutes  *gin.Engine){

	incomingRoutes.GET("/menus", controllers.GetMenus())
	incomingRoutes.GET("/menu/:id", controllers.GetMenuByID())
	incomingRoutes.POST("/menu", controllers.PostMenu())
	incomingRoutes.PATCH("/menu/:id", controllers.UpdateMenu())
	incomingRoutes.DELETE("/menu/:id", controllers.DeleteMenu())
}