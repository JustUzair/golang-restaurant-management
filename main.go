package main

import (
	"log"
	"os"
	// "restaurant-management/controllers"
	"restaurant-management/database"
	"restaurant-management/middleware"
	"restaurant-management/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
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
	router := gin.Default()
	router.Use(gin.Logger())

	routes.UserRoutes(router)
	router.Use(middleware.Authentication())

	routes.FoodRoutes(router)
	routes.MenuRoutes(router)
	routes.OrderRoutes(router)
	routes.TableRoutes(router)
	routes.InvoiceRoutes(router)
	routes.OrderItemRoutes(router)

	router.Run(":" + port)

}
