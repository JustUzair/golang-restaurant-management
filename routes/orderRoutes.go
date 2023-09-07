package routes

import (
	"github.com/gin-gonic/gin"
	controller "restaurant-management/controllers"
)

func OrderRoutes(incomingRoutes *gin.Engine) {
	orderRoute := incomingRoutes.Group("/orders")
	{
		orderRoute.GET("/", controller.GetOrders())
		orderRoute.GET("/:order_id", controller.GetOrder())
		orderRoute.POST("/", controller.CreateOrder())
		orderRoute.PATCH("/:order_id", controller.UpdateOrder())

	}

}
