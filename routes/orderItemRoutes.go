package routes

import (
	"github.com/gin-gonic/gin"
	controller "restaurant-management/controllers"
)

func OrderItemRoutes(incomingRoutes *gin.Engine) {
	orderItemRoute := incomingRoutes.Group("/order-items")
	{
		orderItemRoute.GET("/", controller.GetOrderItems())
		orderItemRoute.GET("/:order_item_id", controller.GetOrderItem())
		orderItemRoute.GET("/order/:order_id", controller.GetOrderItemsByOrder())
		orderItemRoute.POST("/", controller.CreateOrderItem())
		orderItemRoute.PATCH("/:order_item_id", controller.UpdateOrderItem())

	}

}
