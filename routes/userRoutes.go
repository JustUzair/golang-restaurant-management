package routes

import (
	"github.com/gin-gonic/gin"
	controller "restaurant-management/controllers"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	userRoute := incomingRoutes.Group("/api/v1/users")
	{
		userRoute.GET("/", controller.GetUsers())
		userRoute.GET("/:user_id", controller.GetUser())
		userRoute.POST("/signup", controller.SignUp())
		userRoute.POST("/login", controller.Login())

	}

}
