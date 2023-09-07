package routes

import (
	"github.com/gin-gonic/gin"
	controller "restaurant-management/controllers"
)

func InvoiceRoutes(incomingRoutes *gin.Engine) {
	invoiceRoute := incomingRoutes.Group("/invoice")
	{
		invoiceRoute.GET("/", controller.GetInvoices())
		invoiceRoute.GET("/:invoice_id", controller.GetInvoice())
		invoiceRoute.POST("/", controller.CreateInvoice())
		invoiceRoute.PATCH("/:invoice_id", controller.UpdateInvoice())

	}

}
