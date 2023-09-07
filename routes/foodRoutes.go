package routes

import (
	"github.com/gin-gonic/gin"
	controller "restaurant-management/controllers"
)

func FoodRoutes(incomingRoutes *gin.Engine) {
	foodRoute := incomingRoutes.Group("/foods")
	{
		foodRoute.GET("/", controller.GetFoods())
		foodRoute.GET("/:food_id", controller.GetFood())
		foodRoute.POST("/", controller.CreateFood())
		foodRoute.PATCH("/:food_id", controller.UpdateFood())

	}

}
