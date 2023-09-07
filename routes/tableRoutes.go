package routes

import (
	"github.com/gin-gonic/gin"
	controller "restaurant-management/controllers"
)

func TableRoutes(incomingRoutes *gin.Engine) {
	tableRoute := incomingRoutes.Group("/api/v1/tables")
	{
		tableRoute.GET("/", controller.GetTables())
		tableRoute.GET("/:table_id", controller.GetTable())
		tableRoute.POST("/", controller.CreateTable())
		tableRoute.PATCH("/:table_id", controller.UpdateTable())

	}

}
