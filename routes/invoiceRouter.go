package routes

import (
	controllers "ginapp/controllers"

	"github.com/gin-gonic/gin"
)

func InvoiceRoutes (incomingRoutes  *gin.Engine){

	incomingRoutes.GET("/invoices", controllers.GetInvoices())
	incomingRoutes.GET("/invoice/:id", controllers.GetInvoiceByID())
	incomingRoutes.POST("/invoice", controllers.PostInvoice())
	incomingRoutes.PATCH("/invoice/:id", controllers.UpdateInvoice())
	incomingRoutes.DELETE("/invoice/:id", controllers.DeleteInvoice())
}