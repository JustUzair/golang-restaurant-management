package routes

import (
	"github.com/gin-gonic/gin"
	"restaurant-management/controllers"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Group("/users")
	{
		incomingRoutes.GET("/", controllers.GetUsers())
		incomingRoutes.GET("/:user_id", controllers.GetUsers())
		incomingRoutes.GET("/users", controllers.GetUsers())

	}

}
