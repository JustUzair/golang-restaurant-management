package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"os"
	"restaurant-management/database"
	"restaurant-management/middleware"
	"restaurant-management/routes"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

func main() {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatalln(err)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	router := gin.New()
	router.Use(gin.Logger())
	router.Group("/v1/api")
	{
		routes.UserRoutes(router)
		router.Use(middleware.Authentication())

		routes.FoodRoutes(router)
		routes.MenuRoutes(router)
		routes.OrderRoutes(router)
		routes.TableRoutes(router)
		routes.InvoiceRoutes(router)
		routes.OrderItemRoutes(router)
	}

	router.Run(":" + port)

}
