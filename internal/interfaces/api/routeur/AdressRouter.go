package routes

import (
	controllers "ginapp/internal/interfaces/api/controllers"

	"github.com/gin-gonic/gin"
)

func AdressRoutes(incomingRoutes *gin.Engine) {

	incomingRoutes.POST("/adress", controllers.AddAddress())
	incomingRoutes.PUT("/homeAdress/:id", controllers.EditHomeAddress())
	incomingRoutes.PUT("/workAdress/:id", controllers.EditWorkAddress())
	incomingRoutes.DELETE("/adress/:id", controllers.DeleteAddress())
}
