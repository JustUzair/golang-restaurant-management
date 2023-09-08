package main

import (
	"log"
	"os"

	"restaurant-management/middleware"
	"restaurant-management/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

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

	serverErr := router.Run(":" + port)

	if serverErr != nil {
		log.Fatalf("Error starting the server on port %s\n", port)
	}

}
