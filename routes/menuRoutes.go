package routes

import (
	"github.com/gin-gonic/gin"
	controller "restaurant-management/controllers"
)

func MenuRoutes(incomingRoutes *gin.Engine) {
	menuRoute := incomingRoutes.Group("/api/v1/menus")
	{
		menuRoute.GET("/", controller.GetMenus())
		menuRoute.GET("/:menu_id", controller.GetMenu())
		menuRoute.POST("/", controller.CreateMenu())
		menuRoute.PATCH("/:menu_id", controller.UpdateMenu())

	}

}
